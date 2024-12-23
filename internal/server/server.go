package server

import (
	"bytes"
	"compress/gzip"
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"os"
	"sort"
)

// createListener creates a TCP listener on the specified address and port.
// Returns the listener and any error encountered during creation.
func NewListener(address string, port int) (net.Listener, error) {
	return net.Listen("tcp", fmt.Sprintf("%s:%d", address, port))
}

// GzipContent compresses the provided content using gzip compression.
// Returns the compressed content as bytes and any error encountered during compression.
func GzipContent(content []byte) ([]byte, error) {
	var gzippedContent bytes.Buffer
	gzipWriter := gzip.NewWriter(&gzippedContent)
	if _, err := gzipWriter.Write(content); err != nil {
		return nil, err
	}
	if err := gzipWriter.Close(); err != nil {
		return nil, err
	}
	return gzippedContent.Bytes(), nil
}

// FormatDisplayAddress creates a properly formatted HTTP URL string
// from the given address and port.
func FormatDisplayAddress(protocol string, address string, port int) string {
	return fmt.Sprintf("%s://%s:%d", protocol, address, port)
}

// DisplayAddress returns the display-friendly address for the server.
// It converts empty, IPv4, or IPv6 any-address to localhost.
func DisplayAddress(address, anyIPv4, anyIPv6 string) string {
	if address == "" || address == anyIPv4 || address == anyIPv6 {
		return "localhost"
	}
	return address
}

// WithHeaders is a middleware that wraps an http.Handler with common response headers.
func WithHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}

// VerifySignature is a middleware that checks for the presence of a signature header in the request.
func VerifySignature(publicKey ed25519.PublicKey) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			signature := r.Header.Get("X-Signature")
			if signature == "" {
				http.Error(w, "Missing signature", http.StatusUnauthorized)
				return
			}
			if err := r.ParseForm(); err != nil {
				http.Error(w, "Invalid form data", http.StatusBadRequest)
				return
			}
			keys := make([]string, 0, len(r.Form))
			for k := range r.Form {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			var contents []byte
			for _, k := range keys {
				contents = append(contents, []byte(k+r.Form.Get(k))...)
			}
			sig, err := base64.StdEncoding.DecodeString(signature)
			if err != nil {
				http.Error(w, "Invalid signature format", http.StatusBadRequest)
				return
			}
			if !ed25519.Verify(publicKey, contents, sig) {
				http.Error(w, "Invalid signature", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// validateTLSFiles checks if the provided TLS certificate and private key files
// are valid and exist. Returns an error if the validation fails.
func ValidateTLSFiles(certFile, keyFile string) error {
	if certFile != "" || keyFile != "" {
		if certFile == "" || keyFile == "" {
			return fmt.Errorf("both -c and -z options must be provided")
		}
		if _, err := os.Stat(certFile); os.IsNotExist(err) {
			return fmt.Errorf("certificate file '%s' not found", certFile)
		}
		if _, err := os.Stat(keyFile); os.IsNotExist(err) {
			return fmt.Errorf("private key file '%s' not found", keyFile)
		}
	}
	return nil
}
