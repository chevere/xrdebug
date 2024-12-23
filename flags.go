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

	"github.com/xrdebug/xrdebug/internal/cli"
)

var flags = map[string]cli.Flag{
	"a": {
		Variable:    "Address",
		Type:        "string",
		Default:     defaultAddress,
		Description: "IP address to listen on",
	},
	"p": {
		Variable:    "Port",
		Type:        "int",
		Default:     defaultPort,
		Description: "Port to listen on [use 0 for random]",
	},
	"c": {
		Variable:    "TLSCert",
		Type:        "string",
		Default:     "",
		Description: "Path to TLS certificate file [PEM]",
	},
	"z": {
		Variable:    "TLSPrivateKey",
		Type:        "string",
		Default:     "",
		Description: "Path to TLS private key",
	},
	"e": {
		Variable:    "EnableEncryption",
		Type:        "bool",
		Default:     false,
		Description: "Enable end-to-end encryption",
	},
	"k": {
		Variable:    "SymmetricKey",
		Type:        "string",
		Default:     "",
		Description: "[for -e option] Path to symmetric key (AES-GCM AE)",
	},
	"s": {
		Variable:    "EnableSignVerification",
		Type:        "bool",
		Default:     false,
		Description: "Enable sign verification",
	},
	"x": {
		Variable:    "SignPrivateKey",
		Type:        "string",
		Default:     "",
		Description: "[for -s option] Path to private key (ed25519)",
	},
	"n": {
		Variable:    "SessionName",
		Type:        "string",
		Default:     defaultSessionName,
		Description: "Session name",
	},
	"i": {
		Variable:    "Editor",
		Type:        "string",
		Default:     defaultEditor,
		Description: fmt.Sprintf("Editor to use %v", editors),
	},
}
