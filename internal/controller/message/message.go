/*
 * This file is part of xrDebug.
 *
 * (c) Rodolfo Berrios <rodolfo@chevere.org>
 *
 * For the full copyright and license information, please view the LICENSE
 * file that was distributed with this source code.
 */

// Package controller provides HTTP handlers for the xrDebug service.
// It handles incoming debug messages and routes them to appropriate channels.
package message

import (
	"encoding/json"
	"net/http"

	"github.com/xrdebug/xrdebug/internal/cli"
	"github.com/xrdebug/xrdebug/internal/dump"
)

var (
	errParseForm = "Error parsing form data"
	errEmptyForm = "Form data is empty"
)

// Handle returns an http.HandlerFunc that handles incoming debug messages.
// It takes a messages channel where the processed debug messages will be sent,
// and a logger for logging the received messages.
func Handle(messages chan string, logger cli.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, errParseForm, http.StatusBadRequest)
			return
		}
		if len(r.Form) == 0 {
			http.Error(w, errEmptyForm, http.StatusBadRequest)
			return
		}
		msg := dump.New(
			"message",
			r.FormValue("body"),
			r.FormValue("file_path"),
			r.FormValue("file_line"),
			r.FormValue("emote"),
			r.FormValue("topic"),
			r.FormValue("id"),
		)
		jsonMsg, _ := json.Marshal(msg)
		messages <- string(jsonMsg)
		w.WriteHeader(http.StatusOK)
		logger.Printf("Message %s %s", r.RemoteAddr, msg.FileDisplay)
	}
}
