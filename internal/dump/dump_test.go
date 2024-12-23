/*
 * This file is part of xrDebug.
 *
 * (c) Rodolfo Berrios <rodolfo@chevere.org>
 *
 * For the full copyright and license information, please view the LICENSE
 * file that was distributed with this source code.
 */

package dump

import (
	"path/filepath"
	"testing"
)

func TestStripScriptTags(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "no script tags",
			input:    "Hello, world!",
			expected: "Hello, world!",
		},
		{
			name:     "single script tag",
			input:    "Before<script>alert('test')</script>After",
			expected: "BeforeAfter",
		},
		{
			name:     "multiple script tags",
			input:    "Start<script>alert(1)</script>Middle<script>console.log('test')</script>End",
			expected: "StartMiddleEnd",
		},
		{
			name:     "case insensitive",
			input:    "<SCRIPT>test</SCRIPT><script>test</script><ScRiPt>test</ScRiPt>",
			expected: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StripScriptTags(tt.input)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name     string
		action   string
		body     string
		filePath string
		fileLine string
		emote    string
		topic    string
		id       string
		want     *Dump
	}{
		{
			name:     "basic dump",
			action:   "message",
			body:     "test message",
			filePath: "/path/to/file.go",
			fileLine: "42",
			emote:    "🐛",
			topic:    "testing",
			id:       "123",
			want: &Dump{
				Action:           "message",
				Message:          "test message",
				FilePath:         "/path/to/file.go",
				FileLine:         "42",
				FileDisplay:      filepath.Clean("/path/to/file.go") + ":42",
				FileDisplayShort: "file.go:42",
				Emote:            "🐛",
				Topic:            "testing",
				ID:               "123",
			},
		},
		{
			name:     "with script tags in body",
			action:   "message",
			body:     "test<script>alert('hack')</script>message",
			filePath: "/path/to/file.go",
			fileLine: "",
			emote:    "🔍",
			topic:    "security",
			id:       "456",
			want: &Dump{
				Action:           "message",
				Message:          "testmessage",
				FilePath:         "/path/to/file.go",
				FileLine:         "",
				FileDisplay:      filepath.Clean("/path/to/file.go"),
				FileDisplayShort: "file.go",
				Emote:            "🔍",
				Topic:            "security",
				ID:               "456",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.action, tt.body, tt.filePath, tt.fileLine, tt.emote, tt.topic, tt.id)
			if *got != *tt.want {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
