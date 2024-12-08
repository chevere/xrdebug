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

namespace Chevere\Tests\Middlewares;

use Chevere\Tests\src\Traits\CipherTrait;
use Chevere\Tests\src\Traits\Psr17Trait;
use Chevere\xrDebug\Middlewares\DecryptMiddleware;
use PHPUnit\Framework\TestCase;
use Psr\Http\Message\ServerRequestInterface;
use Psr\Http\Server\RequestHandlerInterface;
use function Chevere\Writer\streamTemp;
use function Chevere\xrDebug\encrypt;

final class DecryptMiddlewareTest extends TestCase
{
    use CipherTrait;
    use Psr17Trait;

    public function testNoCipher(): void
    {
        $middleware = new DecryptMiddleware();
        $request = $this->createMock(ServerRequestInterface::class);
        $handler = $this->createMock(RequestHandlerInterface::class);
        $handler
            ->expects($this->once())
            ->method('handle')
            ->with($request);
        $middleware->process($request, $handler);
    }

    public function testCipher(): void
    {
        $message = '{"action":"message"}';
        $cipher = $this->getCipher();
        $middleware = new DecryptMiddleware($cipher);
        $function = function ($request) use ($message) {
            $this->assertSame($message, $request->getBody()->__toString());
        };
        $handler = $this->getRequestHandler($function);
        $cipherBody = encrypt($cipher, $message);
        $request = $this->getRequest()
            ->withBody(streamTemp($cipherBody));
        $middleware->process($request, $handler);
    }
}
