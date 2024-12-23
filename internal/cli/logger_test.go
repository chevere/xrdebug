package cli

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"
)

func TestLoggerPrintf(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr)
	logger := NewLogger()
	testMsg := "test message %s"
	testArg := "arg"
	logger.Printf(testMsg, testArg)
	output := buf.String()
	expected := "test message arg"
	if !strings.Contains(output, expected) {
		t.Errorf("Printf() output = %v, want %v", output, expected)
	}
}
