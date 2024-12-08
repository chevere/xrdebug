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

use Chevere\ThrowableHandler\ThrowableHandler;
use Chevere\Writer\StreamWriter;
use Chevere\Writer\Writers;
use Chevere\Writer\WritersInstance;
use function Chevere\Filesystem\directoryForPath;
use function Chevere\Message\message;
use function Chevere\Router\router;
use function Chevere\Writer\streamFor;

foreach (['/', '/../../../'] as $path) {
    $autoload = __DIR__ . $path . 'vendor/autoload.php';
    if (file_exists($autoload)) {
        require $autoload;

        break;
    }
}

if (! in_array(PHP_SAPI, ['cli', 'phpdbg', 'embed', 'micro'], true)) {
    echo message(
        'xrDebug may only be invoked from a command line, called from %sapi%',
        sapi: PHP_SAPI
    );
    exit(1);
}
$writersInstance = new WritersInstance(
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
$appDirectory = directoryForPath(__DIR__ . '/app');
require_once $appDirectory->path()->getChild('meta.php');
$routes = include $appDirectory->path()->getChild('routes.php');
$router = router($routes);
$routeCollector = $router->routeCollector();
