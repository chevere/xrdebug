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

use phpseclib3\Crypt\EC;
use phpseclib3\Crypt\EC\PrivateKey;

trait KeyTrait
{
    public function getPrivateKey(?string $key = null): PrivateKey
    {
        return $key === null
            ? EC::createKey('ed25519')
            : EC::load($key);
    }
}
