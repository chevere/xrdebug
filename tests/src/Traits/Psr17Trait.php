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

namespace Chevere\Tests\src\Traits;

use Middlewares\Utils\CallableHandler;
use Nyholm\Psr7\Factory\Psr17Factory;
use Psr\Http\Message\ServerRequestInterface;
use Psr\Http\Server\RequestHandlerInterface;

trait Psr17Trait
{
    private function getRequestHandler(?callable $callable = null): RequestHandlerInterface
    {
        return new CallableHandler($callable ?? function () {
        });
    }

    private function getRequest(
        string $method = 'GET',
        string $uri = '/',
    ): ServerRequestInterface {
        return (new Psr17Factory())->createServerRequest($method, $uri);
    }
}
