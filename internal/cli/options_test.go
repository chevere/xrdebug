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
	"reflect"
	"strings"
	"testing"
)

func TestParseOptions(t *testing.T) {
	type args struct {
		flags map[string]Flag
	}
	tests := []struct {
		name         string
		args         args
		errorMessage string
		expected     Options
	}{
		{
			name: "basic options test",
			args: args{
				flags: map[string]Flag{
					"a": {Variable: "Address", Type: "string", Default: "localhost"},
					"p": {Variable: "Port", Type: "int", Default: 8080},
					"n": {Variable: "SessionName", Type: "string", Default: "debug"},
					"e": {Variable: "Editor", Type: "string", Default: "vim"},
					"t": {Variable: "TLSCert", Type: "string", Default: "tls"},
					"k": {Variable: "TLSPrivateKey", Type: "string", Default: "tls_key"},
					"c": {Variable: "EnableEncryption", Type: "bool", Default: false},
					"s": {Variable: "SymmetricKey", Type: "string", Default: "key"},
					"v": {Variable: "EnableSignVerification", Type: "bool", Default: false},
					"g": {Variable: "SignPrivateKey", Type: "string", Default: "sign_key"},
					"version": {Variable: "Version", Type: "bool", Default: false},
				},
			},
			expected: Options{
				Address:                "localhost",
				Port:                   8080,
				SessionName:            "debug",
				Editor:                 "vim",
				TLSCert:                "tls",
				TLSPrivateKey:          "tls_key",
				EnableEncryption:       false,
				SymmetricKey:           "key",
				EnableSignVerification: false,
				SignPrivateKey:         "sign_key",
				Version:				false,
			},
		},
		{
			name: "invalid flag name",
			args: args{
				flags: map[string]Flag{
					"": {Variable: "Address", Type: "string", Default: "localhost", Description: ""},
				},
			},
			errorMessage: "Invalid name for `` flag",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewOptions(tt.args.flags)
			if tt.errorMessage != "" {
				if err == nil {
					t.Errorf("ParseOptions() expected error containing %v, got nil", tt.errorMessage)
					return
				}
				if !strings.Contains(err.Error(), tt.errorMessage) {
					t.Errorf("ParseOptions() error = %v, want error containing %v", err, tt.errorMessage)
				}
				return
			}
			if err != nil {
				t.Errorf("ParseOptions() unexpected error = %v", err)
				return
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("ParseOptions() = %v, want %v", got, tt.expected)
			}
		})
	}
}
