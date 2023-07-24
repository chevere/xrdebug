#!/usr/bin/env php
<?php

/*
 * This file is part of Chevere.
 *
 * (c) Rodolfo Berrios <rodolfo@chevere.org>
 *
 * For the full copyright and license information, please view the LICENSE
 * file that was distributed with this source code.
 */

declare(strict_types=1);

foreach (['/', '/../../../'] as $path) {
    $autoload = __DIR__ . $path . 'vendor/autoload.php';
    if (stream_resolve_include_path($autoload)) {
        require $autoload;

        break;
    }
}

use Chevere\Router\Dependencies;
use Chevere\Router\Dispatcher;
use Chevere\Schwager\DocumentSchema;
use Chevere\Schwager\ServerSchema;
use Chevere\Schwager\Spec;
use Chevere\ThrowableHandler\ThrowableHandler;
use Chevere\Writer\StreamWriter;
use Chevere\Writer\Writers;
use Chevere\Writer\WritersInstance;
use Chevere\XrServer\Build;
use Chevere\XrServer\Debugger;
use Clue\React\Sse\BufferedChannel;
use Colors\Color;
use Psr\Http\Message\ServerRequestInterface;
use React\EventLoop\Loop;
use React\Http\HttpServer;
use React\Http\Middleware\LimitConcurrentRequestsMiddleware;
use React\Http\Middleware\RequestBodyBufferMiddleware;
use React\Http\Middleware\RequestBodyParserMiddleware;
use React\Http\Middleware\StreamingRequestMiddleware;
use React\Socket\SocketServer;
use React\Stream\ThroughStream;
use samejack\PHP\ArgvParser;
use function Chevere\Filesystem\directoryForPath;
use function Chevere\Filesystem\fileForPath;
use function Chevere\Router\router;
use function Chevere\Standard\arrayFilterBoth;
use function Chevere\ThrowableHandler\handleAsConsole;
use function Chevere\Writer\streamFor;
use function Chevere\Writer\writers;
use function Chevere\XrServer\getCipher;
use function Chevere\XrServer\getPrivateKey;
use function Chevere\XrServer\getResponse;

include __DIR__ . '/meta.php';

new WritersInstance(
    (new Writers())
        ->with(
            new StreamWriter(
                streamFor('php://output', 'w')
            )
        )
        ->withError(
            new StreamWriter(
                streamFor('php://stderr', 'w')
            )
        )
);
set_error_handler(ThrowableHandler::ERROR_AS_EXCEPTION);
register_shutdown_function(ThrowableHandler::SHUTDOWN_ERROR_AS_EXCEPTION);
set_exception_handler(ThrowableHandler::CONSOLE);
$color = new Color();
$logger = writers()->log();
$logger->write(
    $color(file_get_contents(__DIR__ . '/logo'))->cyan()
    . "\n"
    . $color(strtr('XR Debug %v (%c) by Rodolfo Berrios', [
        '%v' => XR_SERVER_VERSION,
        '%c' => XR_SERVER_CODENAME,
    ]))->green()
    . "\n\n"
);
$options = (new ArgvParser())->parseConfigs();
if (array_key_exists('h', $options) || array_key_exists('help', $options)) {
    $logger->write(
        <<<LOG
        -p Port (default 27420)
        -e Enable end-to-end encryption
        -k Symmetric key (for -e option)
        -v Enable sign verification
        -s Private key (for -v option)
        -c Cert file for TLS

        LOG
    );
    exit(0);
}
$host = '0.0.0.0';
$port = $options['p'] ?? '27420';
$cert = $options['c'] ?? null;
$isEncryptionEnabled = $options['e'] ?? false;
$isSignVerificationEnabled = $options['v'] ?? false;
$scheme = isset($cert) ? 'tls' : 'tcp';
$uri = "{$scheme}://{$host}:{$port}";
$context = $scheme === 'tcp'
    ? []
    : [
        'tls' => [
            'local_cert' => $cert,
        ],
    ];
$cipher = null;
if ($isEncryptionEnabled) {
    $symmetricKey = $options['k'] ?? null;
    if ($symmetricKey === true) {
        $symmetricKey = null;
    }
    $cipher = getCipher($symmetricKey, $logger, $color);
}
$privateKey = null;
if ($isSignVerificationEnabled) {
    $privateKey = $options['s'] ?? null;
    if ($privateKey === true) {
        $privateKey = null;
    }
    $privateKey = getPrivateKey($privateKey, $logger, $color);
}
$rootDirectory = directoryForPath(__DIR__);
$locksDirectory = $rootDirectory->getChild('locks/');

try {
    $locksDirectory->removeContents();
} catch (Throwable) {
}
$build = new Build(
    $rootDirectory->getChild('app/src/'),
    XR_SERVER_VERSION,
    XR_SERVER_CODENAME,
    $isEncryptionEnabled,
);
$app = fileForPath($rootDirectory->getChild('app/build/')->path()->__toString() . 'en.html');
$app->removeIfExists();
$app->create();
$app->put($build->__toString());
$routes = include 'routes.php';
$dependencies = new Dependencies($routes);
$router = router($routes);
$routeCollector = $router->routeCollector();
$dispatcher = new Dispatcher($routeCollector);
$loop = Loop::get();
$channel = new BufferedChannel();
$debugger = new Debugger($channel, $logger, $cipher);
$containerMap = [
    'app' => $app,
    'channel' => $channel,
    'cipher' => $cipher,
    'debugger' => $debugger,
    'directory' => $locksDirectory,
    'logger' => $logger,
    'loop' => $loop,
    'privateKey' => $privateKey,
    'stream' => new ThroughStream(),
];
$handler = function (ServerRequestInterface $request) use ($dispatcher, $dependencies, $containerMap) {
    try {
        return getResponse($request, $dispatcher, $dependencies, $containerMap);
    } catch (Throwable $e) {
        handleAsConsole($e);
    }
};
$http = new HttpServer(
    $loop,
    new StreamingRequestMiddleware(),
    new LimitConcurrentRequestsMiddleware(100),
    new RequestBodyBufferMiddleware(8 * 1024 * 1024),
    new RequestBodyParserMiddleware(100 * 1024, 1),
    $handler
);
$socket = new SocketServer($uri, $context, $loop);
$http->listen($socket);
$socket->on('error', 'printf');
$scheme = parse_url($socket->getAddress(), PHP_URL_SCHEME);
$httpAddress = strtr(
    $socket->getAddress(),
    [
        'tls' => 'https',
        'tcp' => 'http',
    ]
);
$logger->write(
    <<<LOG
    Server listening on {$scheme} {$httpAddress}
    Press Ctrl+C to quit
    --

    LOG
);

$document = new DocumentSchema(
    api: 'xr',
    name: 'XR Debug API',
    version: '1.0.0'
);
$server = new ServerSchema(
    url: $httpAddress,
    description: 'XR Debug Server',
);
$spec = new Spec($router, $document, $server);
$array = arrayFilterBoth($spec->toArray(), function ($v, $k) {
    return match (true) {
        $v === null => false,
        $v === [] => false,
        $v === '' => false,
        $k === 'required' && $v === true => false,
        $k === 'regex' && $v === '^.*$' => false,
        $k === 'body' && $v === [
            'type' => 'array#map',
        ] => false,
        default => true,
    };
});
$json = json_encode($array, JSON_PRETTY_PRINT);
$file = fileForPath(__DIR__ . '/index.json');
$file->createIfNotExists();
$file->put($json);
$loop->run();
