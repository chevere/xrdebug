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

use Chevere\Filesystem\File;
use Chevere\Filesystem\Interfaces\DirectoryInterface;
use Chevere\Filesystem\Interfaces\FileInterface;
use function Chevere\Filesystem\directoryForPath;
use function realpath;

trait DirectoryTrait
{
    private function getWritableFile(string $id): FileInterface
    {
        return new File(
            $this->getWritableDirectory()->path()->getChild($id)
        );
    }

    private function getWritableDirectory(): DirectoryInterface
    {
        return directoryForPath(
            realpath(__DIR__ . '/../../writable')
        );
    }
}
