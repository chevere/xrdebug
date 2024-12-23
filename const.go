/*
 * This file is part of xrDebug.
 *
 * (c) Rodolfo Berrios <rodolfo@chevere.org>
 *
 * For the full copyright and license information, please view the LICENSE
 * file that was distributed with this source code.
 */

package main

type Replacements struct {
	Logo           string
	Version        string
	Url            string
	Copyright      string
	Name           string
	DisplayAddress string
	GeneratedKeys  string
}

const (
	anyIPv4            = "0.0.0.0"
	anyIPv6            = "::"
	defaultAddress     = ""
	defaultPort        = 27420
	defaultSessionName = name
	defaultEditor      = "vscode"
	templateHeader     = `{{ .Logo }}
{{ .Name }} {{ .Version }}
{{ .Url }}
{{ .Copyright }}
{{ .GeneratedKeys }}
Running at {{ .DisplayAddress }}
--
`
)
