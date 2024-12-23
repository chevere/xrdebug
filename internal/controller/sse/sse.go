/*
 * This file is part of xrDebug.
 *
 * (c) Rodolfo Berrios <rodolfo@chevere.org>
 *
 * For the full copyright and license information, please view the LICENSE
 * file that was distributed with this source code.
 */

// Package sse implements Server-Sent Events functionality
package sse

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/xrdebug/xrdebug/internal/cipher"
	"github.com/xrdebug/xrdebug/internal/cli"
)

// Client represents a connected SSE client with its associated writer
// and flusher for sending events.
type Client struct {
	w       http.ResponseWriter
	flusher http.Flusher
}

// StartDispatcher initializes the SSE message dispatcher that broadcasts
// messages to all connected clients.
func StartDispatcher(messages chan string, clients map[*Client]bool, clientsMu *sync.Mutex, symmetricKey []byte) {
	go func() {
		for msg := range messages {
			if symmetricKey != nil {
				msg = cipher.Encrypt(symmetricKey, msg)
			}
			clientsMu.Lock()
			for client := range clients {

				fmt.Fprintf(client.w, "data: %s\n\n", msg)
				client.flusher.Flush()
			}
			clientsMu.Unlock()
		}
	}()
}

// Handle manages SSE connections, setting up appropriate headers and
// maintaining the connection until the client disconnects.
func Handle(messages chan string, logger cli.Logger, clients map[*Client]bool, clientsMu *sync.Mutex) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		client := &Client{w, w.(http.Flusher)}
		clientsMu.Lock()
		clients[client] = true
		clientsMu.Unlock()
		logger.Printf("Connected %s", r.RemoteAddr)
		defer func() {
			clientsMu.Lock()
			delete(clients, client)
			clientsMu.Unlock()
			logger.Printf("Disconnected %s", r.RemoteAddr)
		}()
		<-r.Context().Done()
	}
}
