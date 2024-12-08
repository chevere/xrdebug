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

namespace Chevere\xrDebug\Middlewares;

use Chevere\Http\Attributes\Response;
use Chevere\Http\Exceptions\MiddlewareException;
use Chevere\Http\Status;
use phpseclib3\Crypt\EC\PrivateKey;
use phpseclib3\Crypt\EC\PublicKey;
use Psr\Http\Message\ResponseInterface;
use Psr\Http\Message\ServerRequestInterface;
use Psr\Http\Server\MiddlewareInterface;
use Psr\Http\Server\RequestHandlerInterface;
use function base64_decode;

#[Response(
    new Status(400)
)]
final class VerifySignatureMiddleware implements MiddlewareInterface
{
    public function __construct(
        private ?PrivateKey $privateKey = null
    ) {
    }

    public function process(
        ServerRequestInterface $request,
        RequestHandlerInterface $handler
    ): ResponseInterface {
        if ($this->privateKey === null) {
            return $handler->handle($request);
        }
        $signature = $request->getHeader('X-Signature');
        if ($signature === []) {
            throw new MiddlewareException(
                message: 'Missing signature',
                code: 400
            );
        }
        $body = $request->getParsedBody();
        $serialize = serialize($body);
        $signature = base64_decode($signature[0], true);
        // @codeCoverageIgnoreStart
        if ($signature === false) {
            throw new MiddlewareException(
                message: 'Invalid signature',
                code: 400
            );
        }
        // @codeCoverageIgnoreEnd
        /** @var PublicKey $publicKey */
        $publicKey = $this->privateKey->getPublicKey();
        if (! $publicKey->verify($serialize, $signature)) {
            throw new MiddlewareException(
                message: 'Invalid signature',
                code: 400
            );
        }

        return $handler->handle($request);
    }
}
