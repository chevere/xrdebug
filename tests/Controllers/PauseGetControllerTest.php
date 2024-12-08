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

use Chevere\Http\Exceptions\ControllerException;
use Chevere\Tests\src\Traits\DirectoryTrait;
use Chevere\xrDebug\Controllers\PauseGetController;
use PHPUnit\Framework\TestCase;

final class PauseGetControllerTest extends TestCase
{
    use DirectoryTrait;

    public function test404(): void
    {
        $id = 'b1cabc9a-145f-11ee-be56-0242ac120002';
        $directory = $this->getWritableDirectory();
        $controller = new PauseGetController($directory);
        $this->expectException(ControllerException::class);
        $this->expectExceptionCode(404);
        $controller->__invoke(id: $id);
    }

    public function test200(): void
    {
        $id = '93683d90-145f-11ee-be56-0242ac120002';
        $array = [
            'stop' => false,
        ];
        $encode = json_encode($array);
        $file = $this->getWritableFile($id);
        $file->createIfNotExists();
        $file->put($encode);
        $directory = $this->getWritableDirectory();
        $controller = new PauseGetController($directory);
        $response = $controller->__invoke(id: $id);
        $this->assertSame($array, $response);
        $file->remove();
    }
}
