/*
 * This file is part of xrDebug.
 *
 * (c) Rodolfo Berrios <rodolfo@chevere.org>
 *
 * For the full copyright and license information, please view the LICENSE
 * file that was distributed with this source code.
 */

package main

import (
	"fmt"
	"slices"
)

var editors = []string{
	"atom",
	"bracket",
	"emacs",
	"espresso",
	"fleet",
	"idea",
	"macvim",
	"netbeans",
	"nova",
	"phpstorm",
	"sublime",
	"textmate",
	"vscode",
	"zed",
}

// validateEditor checks if the provided editor is supported by the application.
// Returns an error if the editor is not supported.
func validateEditor(editor string) error {
	if slices.Contains(editors, editor) {
        return nil
    }
    return fmt.Errorf("editor '%s' not supported", editor)
}
