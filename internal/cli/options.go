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
	"flag"
	"fmt"
)

// Options holds the configuration for the CLI application.
type Options struct {
	// Address specifies the network address to listen on
	Address string
	// Port specifies the network port number to listen on
	Port int
	// TLSCert is the path to the TLS certificate file
	TLSCert string
	// TLSPrivateKey is the path to the TLS private key file
	TLSPrivateKey string
	// EnableEncryption determines if encryption should be used
	EnableEncryption bool
	// SymmetricKey is the path for the key used for encryption (AES-GCM AE)
	SymmetricKey string
	// EnableSignVerification determines if signature verification should be performed
	EnableSignVerification bool
	// SignPrivateKey is the path to the private key used for signing (ed25519)
	SignPrivateKey string
	// SessionName specifies the name of the debug session
	SessionName string
	// Editor specifies the preferred text editor
	Editor string
	// Version specifies the `-version` flag to return the version
	Version bool
}

// NewOptions creates a new Options instance from the provided flags configuration.
// It parses command-line flags and returns Options and any error encountered.
// The flags parameter should contain the flag definitions for all supported options.
func NewOptions(flags map[string]Flag) (Options, error) {
	var validationErrors []error
	flagValues := make(map[string]interface{})
	for name, item := range flags {
		item.Name = name
		flagItem, err := NewFlag(item)
		if err != nil {
			validationErrors = append(validationErrors, err)
			continue
		}
		switch flagItem.Type {
		case "string":
			flagValues[item.Variable] = flag.String(name, item.Default.(string), item.Description)
		case "int":
			flagValues[item.Variable] = flag.Int(name, item.Default.(int), item.Description)
		case "bool":
			flagValues[item.Variable] = flag.Bool(name, item.Default.(bool), item.Description)
		}
	}
	if len(validationErrors) > 0 {
		return Options{}, fmt.Errorf("%v", validationErrors)
	}
	flag.Parse()
	return Options{
		Address:                *flagValues["Address"].(*string),
		Port:                   *flagValues["Port"].(*int),
		TLSCert:                *flagValues["TLSCert"].(*string),
		TLSPrivateKey:          *flagValues["TLSPrivateKey"].(*string),
		EnableEncryption:       *flagValues["EnableEncryption"].(*bool),
		SymmetricKey:           *flagValues["SymmetricKey"].(*string),
		EnableSignVerification: *flagValues["EnableSignVerification"].(*bool),
		SignPrivateKey:         *flagValues["SignPrivateKey"].(*string),
		SessionName:            *flagValues["SessionName"].(*string),
		Editor:                 *flagValues["Editor"].(*string),
		Version:                *flagValues["Version"].(*bool),
	}, nil
}
