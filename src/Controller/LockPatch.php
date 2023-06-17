<?php

/*
 * This file is part of Chevere.
 *
 * (c) Rodolfo Berrios <rodolfo@chevere.org>
 *
 * For the full copyright and license information, please view the LICENSE
 * file that was distributed with this source code.
 */

declare(strict_types=1);

namespace Chevere\XrServer\Controller;

use Chevere\Filesystem\File;
use Chevere\Filesystem\Interfaces\DirectoryInterface;
use Chevere\Http\Controller;
use function Chevere\Parameter\arrayp;
use function Chevere\Parameter\boolean;
use Chevere\Parameter\Interfaces\ArrayTypeParameterInterface;
use function Safe\json_encode;

final class LockPatch extends Controller
{
    public function __construct(
        private DirectoryInterface $directory
    ) {
    }

    public static function acceptResponse(): ArrayTypeParameterInterface
    {
        return arrayp(
            lock: boolean(),
            stop: boolean(),
        );
    }

    /**
     * @return array<string, boolean>
     */
    public function run(string $id): array
    {
        $lockFile = new File(
            $this->directory->path()->getChild('locks/' . $id)
        );
        $lockFile->removeIfExists();
        $lockFile->create();
        $data = [
            'lock' => true,
            'stop' => true,
        ];
        $json = json_encode($data);
        $lockFile->put($json);

        return $data;
    }
}
