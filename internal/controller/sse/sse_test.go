/*
 * This file is part of xrDebug.
 *
 * (c) Rodolfo Berrios <rodolfo@chevere.org>
 *
 * For the full copyright and license information, please view the LICENSE
 * file that was distributed with this source code.
 */

package sse

import (
	"context"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

type mockLogger struct {
	messages []string
}

func (m *mockLogger) Printf(format string, v ...interface{}) {
	m.messages = append(m.messages, format)
}

func TestStartDispatcher(t *testing.T) {
	messages := make(chan string)
	clients := make(map[*Client]bool)
	clientsMu := &sync.Mutex{}
	w := httptest.NewRecorder()
	client := &Client{w: w, flusher: w}
	clientsMu.Lock()
	clients[client] = true
	clientsMu.Unlock()
	StartDispatcher(messages, clients, clientsMu, nil)
	testMessage := "test message"
	messages <- testMessage
	time.Sleep(100 * time.Millisecond)
	response := w.Body.String()
	expected := "data: " + testMessage + "\n\n"
	if response != expected {
		t.Errorf("Expected response %q, got %q", expected, response)
	}
}

func TestHandleDisconnection(t *testing.T) {
	messages := make(chan string)
	clients := make(map[*Client]bool)
	clientsMu := &sync.Mutex{}
	logger := &mockLogger{}
	handler := Handle(messages, logger, clients, clientsMu)
	ctx, cancel := context.WithCancel(context.Background())
	req := httptest.NewRequest("GET", "/sse", nil).WithContext(ctx)
	w := httptest.NewRecorder()
	go func() {
		handler.ServeHTTP(w, req)
	}()
	time.Sleep(100 * time.Millisecond)
	clientsMu.Lock()
	initialClients := len(clients)
	clientsMu.Unlock()
	cancel()
	time.Sleep(100 * time.Millisecond)
	clientsMu.Lock()
	finalClients := len(clients)
	clientsMu.Unlock()
	if initialClients != 1 {
		t.Errorf("Expected 1 initial client, got %d", initialClients)
	}
	if finalClients != 0 {
		t.Errorf("Expected 0 clients after disconnection, got %d", finalClients)
	}
}
