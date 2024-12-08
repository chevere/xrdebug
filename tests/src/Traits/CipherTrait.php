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

use phpseclib3\Crypt\AES;
use phpseclib3\Crypt\Random;

trait CipherTrait
{
    private function getCipher(): AES
    {
        $symmetricKey = Random::string(32);
        $cipher = new AES('gcm');
        $cipher->setKey($symmetricKey);

        return $cipher;
    }
}
