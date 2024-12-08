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

use Chevere\Tests\src\Traits\DirectoryTrait;
use Chevere\xrDebug\Controllers\PausePostController;
use Chevere\xrDebug\Debugger;
use PHPUnit\Framework\TestCase;
use React\Http\Message\ServerRequest;

final class PausePostControllerTest extends TestCase
{
    use DirectoryTrait;

    public function test201(): void
    {
        $id = 'b1cabc9a-145f-11ee-be56-0242ac120002';
        $array = [
            'stop' => false,
        ];
        $encode = json_encode($array);
        $file = $this->getWritableFile($id);
        $file->createIfNotExists();
        $file->put($encode);
        $request = new ServerRequest(
            method: 'POST',
            url: '/locks',
            headers: [],
            body: '',
            version: '1.1',
            serverParams: [
                'REMOTE_ADDR' => '0.0.0.0',
            ]
        );
        $body = [
            'id' => $id,
        ];
        $directory = $this->getWritableDirectory();
        $remoteAddress = $request->getServerParams()['REMOTE_ADDR'];
        $debugger = $this->createMock(Debugger::class);
        $debugger
            ->expects($this->once())
            ->method('sendPause')
            ->with(
                $this->equalTo($body),
                $this->equalTo($remoteAddress)
            );
        $controller = new PausePostController(
            $directory,
            $debugger,
            $remoteAddress
        );
        $controller = $controller->withBody($body);
        $response = $controller->__invoke();
        $decoded = json_decode($file->getContents(), true);
        $this->assertSame($decoded, $response);
        $this->assertTrue($file->exists());
        $this->assertSame($encode, $file->getContents());
        $file->remove();
    }
}
