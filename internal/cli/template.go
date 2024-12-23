/*
 * This file is part of xrDebug.
 *
 * (c) Rodolfo Berrios <rodolfo@chevere.org>
 *
 * For the full copyright and license information, please view the LICENSE
 * file that was distributed with this source code.
 */

package cli

import (
	"io"
	"text/template"
)

// WriteTemplate writes a formatted string to the provided writer using the given template
// and replacements. The tpl is parsed as a text/template and executed with
// the replacements map providing the template variables.
func WriteTemplate(w io.Writer, tpl string, replacements any) error {
	t, err := template.New("xrdebug").Parse(tpl)
	if err != nil {
		return err
	}
	return t.Execute(w, replacements)
}
