/*
 * This file is part of xrDebug.
 *
 * (c) Rodolfo Berrios <rodolfo@chevere.org>
 *
 * For the full copyright and license information, please view the LICENSE
 * file that was distributed with this source code.
 */

package message

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

type mockLogger struct {
	lastMsg string
}

func (m *mockLogger) Printf(format string, v ...interface{}) {
	m.lastMsg = fmt.Sprintf(format, v...)
}

func TestMessage(t *testing.T) {
	tests := []struct {
		name           string
		formData       url.Values
		expectedStatus int
		expectMessage  bool
		expectLog      bool
	}{
		{
			name: "valid message",
			formData: url.Values{
				"body":      {"test message"},
				"file_path": {"test.go"},
				"file_line": {"10"},
				"emote":     {"🚀"},
				"topic":     {"test"},
				"id":        {"123"},
			},
			expectedStatus: http.StatusOK,
			expectMessage:  true,
			expectLog:      true,
		},
		{
			name:           "empty form",
			formData:       url.Values{},
			expectedStatus: http.StatusBadRequest,
			expectMessage:  false,
			expectLog:      false,
		},
		{
			name:           "malformed form data",
			formData:       nil,
			expectedStatus: http.StatusBadRequest,
			expectMessage:  false,
			expectLog:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			messages := make(chan string, 1)
			logger := &mockLogger{}
			handler := Handle(messages, logger)
			var req *http.Request
			s := "%gh&%ij"
			if tt.formData != nil {
				s = tt.formData.Encode()
			}
			req = httptest.NewRequest(http.MethodPost, "/message", strings.NewReader(s))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}
			if tt.expectMessage {
				select {
				case msg := <-messages:
					if msg == "" {
						t.Error("received empty message from channel")
					}
				default:
					t.Error("no message received from channel")
				}
				if tt.expectLog && logger.lastMsg == "" {
					t.Error("expected log message but got none")
				}
				if !tt.expectLog && logger.lastMsg != "" {
					t.Error("got unexpected log message")
				}
			}
		})
	}
}
