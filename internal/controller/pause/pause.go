/*
 * This file is part of xrDebug.
 *
 * (c) Rodolfo Berrios <rodolfo@chevere.org>
 *
 * For the full copyright and license information, please view the LICENSE
 * file that was distributed with this source code.
 */

package pause

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/xrdebug/xrdebug/internal/cli"
	"github.com/xrdebug/xrdebug/internal/dump"
	"github.com/xrdebug/xrdebug/internal/pausectl"
)

// Controller handles HTTP requests for pause operations.
// It manages pause locks and messaging for debugging sessions.
type Controller struct {
	lockManager *pausectl.Manager
	messages    chan string
	logger      cli.Logger
}

// New creates a Controller with the given dependencies.
func New(lockManager *pausectl.Manager, messages chan string, logger cli.Logger) *Controller {
	return &Controller{
		lockManager: lockManager,
		messages:    messages,
		logger:      logger,
	}
}

// Post handles POST /pauses requests.
// It creates a new pause lock and broadcasts the pause message.
func (c *Controller) Post() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.FormValue("id")
		lock, err := c.lockManager.New(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		msg := dump.New(
			"pause",
			r.FormValue("body"),
			r.FormValue("file_path"),
			r.FormValue("file_line"),
			r.FormValue("emote"),
			r.FormValue("topic"),
			id,
		)
		c.logger.Printf("Pause %s %s", r.RemoteAddr, msg.FileDisplay)
		jsonMsg, _ := json.Marshal(msg)
		c.messages <- string(jsonMsg)
		w.Header().Set("Location", fmt.Sprintf("/pauses/%s", id))
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(lock)
	}
}

// Get handles GET /pauses/{id} requests.
// It retrieves the status of an existing pause lock.
func (c *Controller) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		lock, err := c.lockManager.Get(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(lock)
	}
}

// Patch handles PATCH /pauses/{id} requests.
// It updates the state of an existing pause lock.
func (c *Controller) Patch() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		lock, err := c.lockManager.Update(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		c.logger.Printf("Stop %s", r.RemoteAddr)
		json.NewEncoder(w).Encode(lock)
	}
}

// Delete handles DELETE /pauses/{id} requests.
// It removes an existing pause lock and continues execution.
func (c *Controller) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if _, err := c.lockManager.Get(id); err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		c.lockManager.Delete(id)
		c.logger.Printf("Continue %s", r.RemoteAddr)
		w.WriteHeader(http.StatusNoContent)
	}
}
