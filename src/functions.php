<?php

/*
 * This file is part of xrDebug.
 *
 * (c) Rodolfo Berrios <rodolfo@chevere.org>
 *
 * For the full copyright and license information, please view the LICENSE
 * file that was distributed with this source code.
 */

declare(strict_types=1);

namespace Chevere\xrDebug;

use Chevere\Http\Controller;
use Chevere\Http\Exceptions\ControllerException;
use Chevere\Router\Interfaces\DependenciesInterface;
use Chevere\Router\Interfaces\DispatcherInterface;
use Chevere\Writer\Interfaces\WriterInterface;
use Colors\Color;
use LogicException;
use phpseclib3\Crypt\AES;
use phpseclib3\Crypt\Common\SymmetricKey;
use phpseclib3\Crypt\EC;
use phpseclib3\Crypt\EC\PrivateKey;
use phpseclib3\Crypt\Random;
use Psr\Http\Message\ResponseInterface;
use Psr\Http\Message\ServerRequestInterface;
use Psr\Http\Server\MiddlewareInterface;
use Psr\Http\Server\RequestHandlerInterface;
use React\Http\Message\Response;
use React\Stream\ThroughStream;
use Relay\Relay;
use Throwable;
use function base64_decode;
use function Chevere\Http\responseAttribute;
use function json_encode;

function encrypt(SymmetricKey $symmetricKey, string $message, ?string $nonce = null): string
{
    if ($nonce === null) {
        $nonce = Random::string(cipherNonceLength());
    }
    $tag = Random::string(cipherTagLength());
    $symmetricKey = clone $symmetricKey;
    $symmetricKey->setNonce($nonce);
    $symmetricKey->setTag($tag);
    $cipherText = $symmetricKey->encrypt($message);

    return base64_encode($nonce . $cipherText . $symmetricKey->getTag());
}

function decrypt(SymmetricKey $symmetricKey, string $encodedCipherText): string
{
    $decode = base64_decode($encodedCipherText, true);
    if ($decode === false) {
        throw new LogicException('Invalid base64 encoded cipher text');
    }
    $nonce = mb_substr(
        $decode,
        0,
        cipherNonceLength(),
        '8bit'
    );
    $tag = mb_substr(
        $decode,
        -cipherTagLength(),
        null,
        '8bit'
    );
    $cipherText = mb_substr(
        $decode,
        cipherNonceLength(),
        -cipherTagLength(),
        '8bit'
    );
    $symmetricKey = clone $symmetricKey;
    $symmetricKey->setNonce($nonce);
    $symmetricKey->setTag($tag);

    return $symmetricKey->decrypt($cipherText);
}

function cipherNonceLength(): int
{
    return 16;
}

function cipherTagLength(): int
{
    return 16;
}

/**
 * @param array<int|string, string> $body
 */
function getDump(array $body, string $action): Dump
{
    $message = $body['body'] ?? '';
    $message = preg_replace('#<script(.*?)>(.*?)</script>#is', '', $message) ?? '';
    $file = $body['file_path'] ?? '<file>';
    $line = $body['file_line'] ?? '<line>';
    $fileDisplay = $file;
    $fileDisplayShort = basename($file);
    if ($line !== '') {
        $fileDisplay .= ':' . $line;
        $fileDisplayShort .= ':' . $line;
    }

    return new Dump(
        message: $message,
        file_path: $file,
        file_line: $line,
        file_display: $fileDisplay,
        file_display_short: $fileDisplayShort,
        emote: $body['emote'] ?? '',
        topic: $body['topic'] ?? '',
        id: $body['id'] ?? '',
        action: $action,
    );
}

/**
 * @param array<string, mixed> $container
 */
function getResponse(
    ServerRequestInterface $request,
    DispatcherInterface $dispatcher,
    DependenciesInterface $dependencies,
    array $container,
): ResponseInterface {
    $path = $request->getUri()->getPath();
    $body = $request->getParsedBody() ?? [];

    try {
        $routed = $dispatcher->dispatch($request->getMethod(), $path);
        $queue = [];
        $middlewares = $routed->bind()->middlewares();
        foreach ($middlewares as $middlewareName) {
            $className = (string) $middlewareName;
            if (array_key_exists($className, $queue)) {
                continue;
            }
            $middlewareArguments = getArguments($dependencies, $className, $container);
            $queue[$className] = new $className(...$middlewareArguments);
        }
        $queue[] = new class() implements MiddlewareInterface {
            public function process(
                ServerRequestInterface $request,
                RequestHandlerInterface $handler
            ): ResponseInterface {
                return new Response();
            }
        };
        $relay = new Relay($queue);
        $response = $relay->handle($request);
    } catch (Throwable $e) {
        return new Response(404);
    }
    $container = array_merge($container, [
        'request' => $request,
        'lastEventId' => $request->getHeaderLine('Last-Event-ID'),
        'remoteAddress' => $request->getServerParams()['REMOTE_ADDR'] ?? '',
    ]);
    $view = $routed->bind()->view();
    $controllerName = $routed->bind()->controllerName()->__toString();
    $controllerArguments = getArguments($dependencies, $controllerName, $container);
    /** @var Controller $controller */
    $controller = new $controllerName(...$controllerArguments);
    if ($request->getMethod() === 'POST') {
        try {
            $controller = $controller->withBody((array) $body);
        } catch (Throwable $e) {
            $headers = [
                'Content-Type' => 'text/json',
            ];
            $error = json_encode([
                'error' => $e->getMessage(),
            ]);

            return new Response(
                status: 400,
                headers: $headers,
                body: strval($error)
            );
        }
    }

    try {
        $response = $controller->__invoke(...$routed->arguments());
    } catch (ControllerException $e) {
        return new Response(
            $e->getCode(),
        );
    }
    $isStream = $response instanceof ThroughStream;
    $status = responseAttribute($controllerName)->status->primary;
    $headers = [
        'Content-Type' => match (true) {
            $isStream => 'text/event-stream',
            $view === 'spa/GET' => 'text/html',
            default => 'text/json',
        },
    ];
    /** @var ThroughStream|string $body */
    $body = match (true) {
        $isStream => $response,
        $view === 'spa/GET' => $response,
        default => json_encode($response),
    };

    return new Response($status, $headers, $body);
}

/**
 * @param array<string, mixed> $container
 * @return array<string, mixed>
 */
function getArguments(
    DependenciesInterface $dependencies,
    string $class,
    array $container
): array {
    $arguments = [];
    foreach ($dependencies->get($class)->keys() as $key) {
        $arguments[$key] = array_key_exists($key, $container)
            ? $container[$key]
            : throw new LogicException("Missing container key {$key}");
    }

    return $arguments;
}

function getCipher(
    ?string $symmetricKey,
    WriterInterface $logger,
    Color $color
): AES {
    if ($symmetricKey === null) {
        $symmetricKey = Random::string(32);
        $logger->write(
            $color('INFO: Generated encryption key (empty -k)')->magenta() . "\n"
        );
    } else {
        $symmetricKey = base64_decode($symmetricKey, true);
        if ($symmetricKey === false) {
            throw new LogicException('Invalid symmetric key');
        }
    }
    $cipher = new AES('gcm');
    $cipher->setKey($symmetricKey);
    $encryptionKeyDisplay = base64_encode($symmetricKey);
    $logger->write(<<<LOG

    > ENCRYPTION KEY
    {$encryptionKeyDisplay}

    LOG);

    return $cipher;
}

function getPrivateKey(
    ?string $privateKey,
    WriterInterface $logger,
    Color $color
): PrivateKey {
    if ($privateKey === null) {
        $privateKey = EC::createKey('ed25519');
        $logger->write(
            $color('INFO: Generated private key (empty -x)')->magenta() . "\n"
        );
    } else {
        $privateKey = EC::load($privateKey);
    }
    $privateKeyDisplay = $privateKey->toString('PKCS8');
    $logger->write(
        <<<LOG

        > PRIVATE KEY
        {$privateKeyDisplay}

        LOG
    );
    /** @var PrivateKey */
    return $privateKey;
}
