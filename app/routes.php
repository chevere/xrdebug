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

use Chevere\xrDebug\Controllers\MessagePostController;
use Chevere\xrDebug\Controllers\PauseDeleteController;
use Chevere\xrDebug\Controllers\PauseGetController;
use Chevere\xrDebug\Controllers\PausePatchController;
use Chevere\xrDebug\Controllers\PausePostController;
use Chevere\xrDebug\Controllers\SPAController;
use Chevere\xrDebug\Controllers\StreamController;
use Chevere\xrDebug\Middlewares\DecryptMiddleware;
use Chevere\xrDebug\Middlewares\VerifySignatureMiddleware;
use function Chevere\Router\bind;
use function Chevere\Router\route;
use function Chevere\Router\routes;

return routes(
    route(
        path: '/',
        GET: bind(SPAController::class)->withView('spa')
    ),
    route(
        path: '/pauses',
        POST: bind(
            PausePostController::class,
            // DecryptMiddleware::class,
            VerifySignatureMiddleware::class,
        ),
    ),
    route(
        path: '/pauses/{id}',
        GET: bind(
            PauseGetController::class,
            VerifySignatureMiddleware::class,
        ),
        PATCH: bind(
            PausePatchController::class,
            DecryptMiddleware::class,
        ),
        DELETE: bind(
            PauseDeleteController::class,
            DecryptMiddleware::class,
        ),
    ),
    route(
        path: '/messages',
        POST: bind(
            MessagePostController::class,
            // DecryptMiddleware::class,
            VerifySignatureMiddleware::class
        ),
    ),
    route(
        path: '/stream',
        GET: StreamController::class,
    ),
);
