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
use Chevere\xrDebug\Controllers\SPAController;
use PHPUnit\Framework\TestCase;

final class SPAControllerTest extends TestCase
{
    use DirectoryTrait;

    public function test200(): void
    {
        $contents = 'contents';
        $file = $this->getWritableFile('app.html');
        $file->createIfNotExists();
        $file->put($contents);
        $controller = new SPAController($file);
        $response = $controller->__invoke();
        $this->assertSame($contents, $response);
        $file->remove();
    }
}
