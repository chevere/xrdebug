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

use Chevere\Http\Exceptions\MiddlewareException;
use Chevere\Tests\src\Traits\CipherTrait;
use Chevere\Tests\src\Traits\KeyTrait;
use Chevere\Tests\src\Traits\Psr17Trait;
use Chevere\xrDebug\Middlewares\VerifySignatureMiddleware;
use PHPUnit\Framework\TestCase;
use Psr\Http\Message\ServerRequestInterface;
use Psr\Http\Server\RequestHandlerInterface;

final class VerifySignatureMiddlewareTest extends TestCase
{
    use CipherTrait;
    use KeyTrait;
    use Psr17Trait;

    public function testNullKey(): void
    {
        $key = null;
        $middleware = new VerifySignatureMiddleware($key);
        $request = $this->createMock(ServerRequestInterface::class);
        $handler = $this->createMock(RequestHandlerInterface::class);
        $handler
            ->expects($this->once())
            ->method('handle')
            ->with($request);
        $middleware->process($request, $handler);
    }

    public function testMissingSignature400(): void
    {
        $privateKey = $this->getPrivateKey();
        $middleware = new VerifySignatureMiddleware($privateKey);
        $request = $this->getRequest();
        $handler = $this->getRequestHandler();
        $this->expectException(MiddlewareException::class);
        $this->expectExceptionMessage('Missing signature');
        $this->expectExceptionCode(400);
        $middleware->process($request, $handler);
    }

    public function testInvalidSignatureKey(): void
    {
        $privateKey = $this->getPrivateKey();
        $middleware = new VerifySignatureMiddleware($privateKey);
        $data = ['test'];
        $serialize = serialize($data);
        $signature = $this->getPrivateKey()->sign($serialize);
        $signatureDisplay = base64_encode($signature);
        $request = $this->getRequest()
            ->withParsedBody($data)
            ->withHeader(
                'X-Signature',
                $signatureDisplay
            );
        $handler = $this->getRequestHandler();
        $this->expectException(MiddlewareException::class);
        $this->expectExceptionMessage('Invalid signature');
        $this->expectExceptionCode(400);
        $middleware->process($request, $handler);
    }

    public function dataProvider(): array
    {
        return [
            [
                null,
                [],
                ['test'],
            ],
        ];
    }

    /**
     * @dataProvider dataProvider
     */
    public function testSuccess(?array $data): void
    {
        $privateKey = $this->getPrivateKey();
        $middleware = new VerifySignatureMiddleware($privateKey);
        $serialize = serialize($data);
        $signature = $privateKey->sign($serialize);
        $signatureDisplay = base64_encode($signature);
        $request = $this->getRequest()
            ->withParsedBody($data)
            ->withHeader(
                'X-Signature',
                $signatureDisplay
            );
        $handler = $this->getRequestHandler();
        $this->expectNotToPerformAssertions();
        $middleware->process($request, $handler);
    }
}
