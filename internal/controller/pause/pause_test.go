package pause

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/xrdebug/xrdebug/internal/pausectl"
)

type mockLogger struct{}

func (m *mockLogger) Printf(format string, v ...interface{}) {}

func setupTest() (*Controller, chan string) {
	messages := make(chan string, 10)
	manager := pausectl.NewManager(5*time.Minute, 10*time.Minute)
	logger := &mockLogger{}
	controller := New(manager, messages, logger)
	return controller, messages
}

func TestPauseController(t *testing.T) {
	controller, messages := setupTest()
	lockID := "123-456-789"
	t.Run("GET non-existent lock", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/pauses/"+lockID, nil)
		w := httptest.NewRecorder()
		controller.Get()(w, req)
		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
		}
	})
	t.Run("PATCH non-existent lock", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/pauses/"+lockID, nil)
		w := httptest.NewRecorder()
		controller.Patch()(w, req)
		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
		}
	})
	t.Run("DELETE non-existent lock", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/pauses/"+lockID, nil)
		w := httptest.NewRecorder()
		controller.Delete()(w, req)
		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
		}
	})
	t.Run("POST create lock", func(t *testing.T) {
		body := strings.NewReader("id=" + lockID + "&body=test&file_path=/test&file_line=1&emote=🔍&topic=debug")
		req := httptest.NewRequest(http.MethodPost, "/pauses", body)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		controller.Post()(w, req)
		if w.Code != http.StatusCreated {
			t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
		}
		var lock pausectl.Lock
		if err := json.NewDecoder(w.Body).Decode(&lock); err != nil {
			t.Fatal(err)
		}
		select {
		case msg := <-messages:
			if !strings.Contains(msg, lockID) {
				t.Error("Expected message to contain lock ID")
			}
		default:
			t.Error("Expected message to be sent")
		}
		w = httptest.NewRecorder()
		controller.Post()(w, req)
		if w.Code != http.StatusConflict {
			t.Errorf("Expected status %d, got %d", http.StatusConflict, w.Code)
		}
		controller.lockManager.Delete(lockID)
	})
	t.Run("GET existing lock", func(t *testing.T) {
		_, err := controller.lockManager.New(lockID)
		if err != nil {
			t.Fatal(err)
		}
		req := httptest.NewRequest(http.MethodGet, "/pauses/"+lockID, nil)
		req.SetPathValue("id", lockID)
		w := httptest.NewRecorder()
		controller.Get()(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}
		var lock pausectl.Lock
		if err := json.NewDecoder(w.Body).Decode(&lock); err != nil {
			t.Fatal(err)
		}
		controller.lockManager.Delete(lockID)
	})
	t.Run("PATCH existing lock", func(t *testing.T) {
		_, err := controller.lockManager.New(lockID)
		if err != nil {
			t.Fatal(err)
		}
		req := httptest.NewRequest(http.MethodPatch, "/pauses/"+lockID, nil)
		req.SetPathValue("id", lockID)
		w := httptest.NewRecorder()
		controller.Patch()(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}
		controller.lockManager.Delete(lockID)
	})
	t.Run("DELETE existing lock", func(t *testing.T) {
		_, err := controller.lockManager.New(lockID)
		if err != nil {
			t.Fatal(err)
		}
		req := httptest.NewRequest(http.MethodDelete, "/pauses/"+lockID, nil)
		req.SetPathValue("id", lockID)
		w := httptest.NewRecorder()
		controller.Delete()(w, req)
		if w.Code != http.StatusNoContent {
			t.Errorf("Expected status %d, got %d", http.StatusNoContent, w.Code)
		}
	})
}
