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

func TestNewFlag(t *testing.T) {
	type args struct {
		flag Flag
	}
	tests := []struct {
		name         string
		args         args
		expected     *Flag
		errorMessage string
	}{
		{
			name: "valid flag",
			args: args{
				flag: Flag{
					Name:        "test",
					Variable:    "test_var",
					Type:        "string",
					Default:     "default",
					Description: "test description",
				},
			},
			expected: &Flag{
				Name:        "test",
				Variable:    "test_var",
				Type:        "string",
				Default:     "default",
				Description: "test description",
			},
		},
		{
			name: "empty name",
			args: args{
				flag: Flag{
					Name:        "",
					Variable:    "test_var",
					Type:        "string",
					Default:     "default",
					Description: "test description",
				},
			},
			errorMessage: "Invalid name for `` flag",
		},
		{
			name: "empty variable",
			args: args{
				flag: Flag{
					Name:        "test",
					Variable:    "",
					Type:        "string",
					Default:     "default",
					Description: "test description",
				},
			},
			errorMessage: "Invalid variable for `test` flag",
		},
		{
			name: "unsupported type",
			args: args{
				flag: Flag{
					Name:        "test",
					Variable:    "test_var",
					Type:        "float",
					Default:     1.23,
					Description: "test description",
				},
			},
			errorMessage: "Unsupported type for `test` flag",
		},
		{
			name: "invalid default type",
			args: args{
				flag: Flag{
					Name:        "test",
					Variable:    "test_var",
					Type:        "string",
					Default:     123,
					Description: "test description",
				},
			},
			errorMessage: "Invalid default value for `test` flag",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewFlag(tt.args.flag)
			if tt.errorMessage != "" {
				if err == nil {
					t.Error("expected error but got none")
				}
				if !strings.Contains(err.Error(), tt.errorMessage) {
					t.Errorf("error message does not contain %q, got: %v", tt.errorMessage, err)
				}
				if got != nil {
					t.Error("expected nil result when error occurs")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(tt.expected, got) {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
			assertFlagFields(t, got, tt.args.flag)
		})
	}
}

// assertFlagFields verifies that a Flag instance matches the expected Flag values
func assertFlagFields(t *testing.T, got *Flag, flag Flag) {
	t.Helper()
	fields := map[string]struct {
		got, want interface{}
	}{
		"Name":         {got.Name, flag.Name},
		"Variable":     {got.Variable, flag.Variable},
		"Type":         {got.Type, flag.Type},
		"DefaultValue": {got.Default, flag.Default},
	}

	for field, values := range fields {
		if values.got != values.want {
			t.Errorf("expected %s %v, got %v", field, values.want, values.got)
		}
	}
}
