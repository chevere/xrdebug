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
	"regexp"
)

// scriptTagPattern is a compiled regex for matching and removing script tags
const scriptTagPattern = `(?i)<script.*?>.*?</script>`

var scriptTagRegex = regexp.MustCompile(scriptTagPattern)

// Dump represents a debug message with associated metadata and file information
type Dump struct {
	// Action represents the debug action type
	Action string `json:"action"`
	// Message contains the debug message content
	Message string `json:"message"`
	// FilePath contains the full path to the debugged file
	FilePath string `json:"file_path"`
	// FileLine represents the line number in the debugged file
	FileLine string `json:"file_line"`
	// FileDisplay contains the formatted file path with line number
	FileDisplay string `json:"file_display"`
	// FileDisplayShort contains the basename of the file with line number
	FileDisplayShort string `json:"file_display_short"`
	// Emote represents an emotion indicator
	Emote string `json:"emote"`
	// Topic categorizes the debug message
	Topic string `json:"topic"`
	// ID uniquely identifies the debug entry
	ID string `json:"id"`
}

// StripScriptTags removes any script tags from the input string for security
func StripScriptTags(input string) string {
	return scriptTagRegex.ReplaceAllString(input, "")
}

// New creates a new Dump instance with the provided parameters
func New(action, body, filePath, fileLine, emote, topic, id string) *Dump {
	body = StripScriptTags(body)
	fileDisplay := filepath.Clean(filePath)
	fileDisplayShort := filepath.Base(fileDisplay)
	if fileLine != "" {
		fileDisplay += ":" + fileLine
		fileDisplayShort += ":" + fileLine
	}
	return &Dump{
		Action:           action,
		Message:          body,
		FilePath:         filePath,
		FileLine:         fileLine,
		FileDisplay:      fileDisplay,
		FileDisplayShort: fileDisplayShort,
		Emote:            emote,
		Topic:            topic,
		ID:               id,
	}
}
