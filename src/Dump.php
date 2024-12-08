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

namespace Chevere\xrDebug;

use LogicException;
use function json_encode;

final class Dump
{
    public function __construct(
        public readonly string $message,
        public readonly string $file_path,
        public readonly string $file_line,
        public readonly string $file_display,
        public readonly string $file_display_short,
        public readonly string $emote,
        public readonly string $topic,
        public readonly string $id,
        public readonly string $action,
    ) {
    }

    /**
     * @return array<string, string>
     */
    public function toArray(): array
    {
        return (array) $this;
    }

    public function toJson(): string
    {
        $json = json_encode($this->toArray());
        if ($json === false) {
            throw new LogicException('Unable to encode to JSON');
        }

        return $json;
    }
}
