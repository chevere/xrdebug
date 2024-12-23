/*
 * This file is part of xrDebug.
 *
 * (c) Rodolfo Berrios <rodolfo@chevere.org>
 *
 * For the full copyright and license information, please view the LICENSE
 * file that was distributed with this source code.
 */

package spa

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSpa(t *testing.T) {
	content := []byte("Test Content")
	handler := Handle(content)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if got := w.Header().Get("Content-Type"); got != "text/html" {
		t.Errorf("Content-Type = %q; want text/html", got)
	}
	if got := w.Code; got != http.StatusOK {
		t.Errorf("Status Code = %d; want %d", got, http.StatusOK)
	}
	if got := w.Body.String(); got != string(content) {
		t.Errorf("Body = %q; want %q", got, string(content))
	}
}

func TestSpaNotFound(t *testing.T) {
	content := []byte("<html><body>Test Content</body></html>")
	handler := Handle(content)
	req := httptest.NewRequest(http.MethodGet, "/fail", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if got := w.Code; got != http.StatusNotFound {
		t.Errorf("Status Code = %d; want %d", got, http.StatusNotFound)
	}
}
