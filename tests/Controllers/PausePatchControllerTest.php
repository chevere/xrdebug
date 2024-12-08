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
use Chevere\xrDebug\Controllers\PausePatchController;
use PHPUnit\Framework\TestCase;

final class PausePatchControllerTest extends TestCase
{
    use DirectoryTrait;

    public function test404(): void
    {
        $id = 'b1cabc9a-145f-11ee-be56-0242ac120002';
        $directory = $this->getWritableDirectory();
        $controller = new PausePatchController($directory);
        $this->expectException(ControllerException::class);
        $this->expectExceptionCode(404);
        $controller->__invoke(id: $id);
    }

    public function test200(): void
    {
        $id = '93683d90-145f-11ee-be56-0242ac120002';
        $directory = $this->getWritableDirectory();
        $path = $directory->path()->getChild($id);
        $file = new File($path);
        $file->createIfNotExists();
        $controller = new PausePatchController($directory);
        $response = $controller->__invoke(id: $id);
        $decoded = json_decode($file->getContents(), true);
        $expected = [
            'stop' => true,
        ];
        $this->assertSame($expected, $response);
        $this->assertSame($expected, $decoded);
        $this->assertTrue($file->exists());
        $file->remove();
    }
}
