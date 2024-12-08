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

namespace Chevere\Tests\Controllers;

use Chevere\Filesystem\File;
use Chevere\Http\Exceptions\ControllerException;
use Chevere\Tests\src\Traits\DirectoryTrait;
use Chevere\xrDebug\Controllers\PauseDeleteController;
use PHPUnit\Framework\TestCase;

final class PauseDeleteControllerTest extends TestCase
{
    use DirectoryTrait;

    public function test404(): void
    {
        $id = 'b1cabc9a-145f-11ee-be56-0242ac120002';
        $directory = $this->getWritableDirectory();
        $path = $directory->path()->getChild($id);
        $file = new File($path);
        $file->removeIfExists();
        $controller = new PauseDeleteController($directory);
        $this->expectException(ControllerException::class);
        $this->expectExceptionCode(404);
        $controller->__invoke(id: $id);
    }

    public function test204(): void
    {
        $id = '93683d90-145f-11ee-be56-0242ac120002';
        $directory = $this->getWritableDirectory();
        $path = $directory->path()->getChild($id);
        $file = new File($path);
        $file->createIfNotExists();
        $controller = new PauseDeleteController($directory);
        $response = $controller->__invoke(id: $id);
        $this->assertSame(null, $response);
        $this->assertFalse($file->exists());
    }
}
