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
	"bytes"
	"testing"
)

func TestWriteTemplate(t *testing.T) {
	tests := []struct {
		name     string
		template string
		replace  map[string]string
		expected string
		error    bool
	}{
		{
			name:     "replacement",
			template: "Hello {{.Name}}!",
			replace: map[string]string{
				"Name": "World",
			},
			expected: "Hello World!",
		},
		{
			name:     "invalid template",
			template: "Hello {{.Name}!",
			error:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := WriteTemplate(&buf, tt.template, tt.replace)
			if tt.error && err == nil {
				t.Fatal("WriteTemplate() expected error, got nil")
			}
			if !tt.error && err != nil {
				t.Fatalf("WriteTemplate() unexpected error: %v", err)
			}
			if !tt.error && buf.String() != tt.expected {
				t.Fatalf("WriteTemplate() = %v, want %v", buf.String(), tt.expected)
			}
		})
	}
}
