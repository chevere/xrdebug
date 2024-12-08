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

use ArgumentCountError;
use Chevere\xrDebug\Controllers\MessagePostController;
use Chevere\xrDebug\Debugger;
use PHPUnit\Framework\TestCase;
use function Chevere\Parameter\assertArray;

final class MessagePostControllerTest extends TestCase
{
    public function dataProvider(): array
    {
        return [
            [
                [
                    'body' => 'body',
                    'emote' => 'emote',
                    'file_line' => 'file_line',
                    'file_path' => 'file_path',
                    'id' => 'id',
                    'topic' => 'topic',
                ],
            ],
        ];
    }

    /**
     * @dataProvider dataProvider
     */
    public function testAcceptBody(array $binds): void
    {
        $body = MessagePostController::acceptBody();
        foreach (array_keys($binds) as $key) {
            assertArray($body, $binds);
            unset($binds[$key]);
        }
        $this->expectException(ArgumentCountError::class);
        $this->expectExceptionMessage(
            <<<PLAIN
            Requires minimum **1** optional argument(s), **0** provided
            PLAIN
        );
        assertArray($body, []);
    }

    /**
     * @dataProvider dataProvider
     */
    public function test201(array $body): void
    {
        $remoteAddress = 'remote_address';
        $debugger = $this->createMock(Debugger::class);
        $debugger
            ->expects($this->once())
            ->method('sendMessage')
            ->with(
                $this->equalTo($body),
                $this->equalTo('remote_address')
            );

        $controller = new MessagePostController(
            $debugger,
            $remoteAddress
        );
        $controller = $controller->withBody($body);
        $response = $controller->__invoke();
        $this->assertSame(null, $response);
    }
}
