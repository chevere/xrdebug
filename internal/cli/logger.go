/*
 * This file is part of xrDebug.
 *
 * (c) Rodolfo Berrios <rodolfo@chevere.org>
 *
 * For the full copyright and license information, please view the LICENSE
 * file that was distributed with this source code.
 */

package cli

import "log"

// Logger defines the interface for logging operations
type Logger interface {
	// Printf formats and prints a message according to the format specifier
	Printf(format string, v ...interface{})
}

// stdLogger implements the Logger interface using the standard log package
type stdLogger struct{}

func (l *stdLogger) Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

// NewLogger creates and returns a new Logger instance
func NewLogger() Logger {
	return &stdLogger{}
}
