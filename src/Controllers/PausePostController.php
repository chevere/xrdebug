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

namespace Chevere\xrDebug\Controllers;

use Chevere\Filesystem\File;
use Chevere\Filesystem\Interfaces\DirectoryInterface;
use Chevere\Http\Attributes\Description;
use Chevere\Http\Attributes\Request;
use Chevere\Http\Attributes\Response;
use Chevere\Http\Controller;
use Chevere\Http\Header;
use Chevere\Http\Status;
use Chevere\Parameter\Interfaces\ArrayParameterInterface;
use Chevere\Parameter\Interfaces\ParameterInterface;
use Chevere\xrDebug\Constants\UrlPathRegex;
use Chevere\xrDebug\Debugger;
use function Chevere\Parameter\arrayp;
use function Chevere\Parameter\bool;
use function Chevere\Parameter\string;
use function json_encode;

#[Description('Create a pause')]
#[Request(
    new Header('Content-Type', 'application/json'),
)]
#[Response(
    new Status(201),
    new Header('Content-Type', 'application/json')
)]
final class PausePostController extends Controller
{
    public function __construct(
        private DirectoryInterface $directory,
        private Debugger $debugger,
        private string $remoteAddress,
    ) {
    }

    public static function acceptResponse(): ParameterInterface
    {
        return arrayp(
            stop: bool(),
        );
    }

    public static function acceptBody(): ArrayParameterInterface
    {
        return arrayp(
            id: string(UrlPathRegex::UUID),
        )->withOptional(
            body: string('/.*?/'),
            emote: string(),
            file_line: string(),
            file_path: string(),
            topic: string(),
        );
    }

    protected function main(): array
    {
        $id = $this->body()->required('id')->string();
        $path = $this->directory->path()->getChild($id);
        $file = new File($path);
        $file->removeIfExists();
        $file->create();
        $data = [
            'stop' => false,
        ];
        $encoded = json_encode($data);
        $encoded = strval($encoded);
        $file->put($encoded);
        $this->debugger->sendPause(
            $this->body()->toArray(),
            $this->remoteAddress
        );

        return $data;
    }
}
