/*
 * This file is part of xrDebug.
 *
 * (c) Rodolfo Berrios <rodolfo@chevere.org>
 *
 * For the full copyright and license information, please view the LICENSE
 * file that was distributed with this source code.
 */

// Package main implements the xrDebug server, a web-based debugging interface
// that provides real-time debug information and execution control for network applications.
package main

import (
	"crypto/ed25519"
	"embed"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/xrdebug/xrdebug/internal/build"
	"github.com/xrdebug/xrdebug/internal/cipher"
	"github.com/xrdebug/xrdebug/internal/cli"
	"github.com/xrdebug/xrdebug/internal/controller/message"
	"github.com/xrdebug/xrdebug/internal/controller/pause"
	"github.com/xrdebug/xrdebug/internal/controller/spa"
	"github.com/xrdebug/xrdebug/internal/controller/sse"
	"github.com/xrdebug/xrdebug/internal/pausectl"
	"github.com/xrdebug/xrdebug/internal/server"
)

//go:embed web/* assets/*
var filesystem embed.FS

type ServerDeps struct {
	Logger    cli.Logger
	Messages  chan string
	Clients   map[*sse.Client]bool
	ClientsMu *sync.Mutex
}

func middleware(handler http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}
	return handler
}

func main() {
	var clientsMu sync.Mutex
	deps := &ServerDeps{
		Logger:    cli.NewLogger(),
		Messages:  make(chan string, 100),
		Clients:   make(map[*sse.Client]bool),
		ClientsMu: &clientsMu,
	}
	if err := run(deps); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

// run initializes and starts the xrDebug server. It sets up HTTP routes,
// initializes message channels, and manages client connections.
func run(deps *ServerDeps) error {
	options, err := cli.NewOptions(flags)
	if err != nil {
		return err
	}
	if options.Version {
		fmt.Printf("%s %s\n", name, version)
		os.Exit(0)
	}
	if err := validateEditor(options.Editor); err != nil {
		return err
	}
	if err := server.ValidateTLSFiles(options.TLSCert, options.TLSPrivateKey); err != nil {
		return err
	}
	protocol := "http"
	if options.TLSCert != "" && options.TLSPrivateKey != "" {
		protocol += "s"
	}
	var generatedKeys []string
	var signPrivateKey ed25519.PrivateKey
	var symmetricKey []byte
	if options.EnableSignVerification {
		signPrivateKey, err = cipher.LoadPrivateKey(options.SignPrivateKey)
		if err != nil {
			return err
		}
		if options.SignPrivateKey == "" {
			pemKey, err := cipher.PemKey(signPrivateKey)
			if err != nil {
				return err
			}
			generatedKeys = append(generatedKeys, pemKey)
		}
	}
	if options.EnableEncryption {
		symmetricKey, err = cipher.LoadSymmetricKey(options.SymmetricKey)
		if err != nil {
			return err
		}
		if options.SymmetricKey == "" {
			symmetricKeyDisplay := cipher.Base64(symmetricKey)
			generatedKeys = append(generatedKeys,
				fmt.Sprintf("ENCRYPTION KEY\n%s", symmetricKeyDisplay))
		}
	}
	html, err := filesystem.ReadFile("web/index.html")
	if err != nil {
		return err
	}
	build, err := build.New(html, filesystem, "web/", version, options.SessionName, options.Editor, options.EnableEncryption, options.EnableSignVerification)
	if err != nil {
		return err
	}
	content := build.Bytes()
	gzipped, err := server.GzipContent(content)
	if err != nil {
		return err
	}
	displayAddress := server.DisplayAddress(options.Address, anyIPv4, anyIPv6)
	listener, err := server.NewListener(options.Address, options.Port)
	if err != nil {
		return err
	}
	displayPort := listener.Addr().(*net.TCPAddr).Port
	displayAddress = server.FormatDisplayAddress(protocol, displayAddress, displayPort)
	lockManager := pausectl.NewManager(5*time.Minute, 1*time.Minute)
	pauseController := pause.New(lockManager, deps.Messages, deps.Logger)
	sse.StartDispatcher(deps.Messages, deps.Clients, deps.ClientsMu, symmetricKey)
	middlewares := []func(http.Handler) http.Handler{server.WithHeaders}
	clientSignMiddleware := append([]func(http.Handler) http.Handler{}, middlewares...)
	if options.EnableSignVerification {
		clientSignMiddleware = append(
			clientSignMiddleware,
			server.VerifySignature(signPrivateKey.Public().(ed25519.PublicKey)),
		)
	}
	http.Handle("GET /", middleware(spa.Handle(gzipped), middlewares...))
	http.Handle("POST /messages", middleware(message.Handle(deps.Messages, deps.Logger), clientSignMiddleware...))
	http.Handle("POST /pauses", middleware(pauseController.Post(), clientSignMiddleware...))
	http.Handle("GET /pauses/{id}", middleware(pauseController.Get(), clientSignMiddleware...))
	http.Handle("GET /stream", middleware(sse.Handle(deps.Messages, deps.Logger, deps.Clients, deps.ClientsMu), middlewares...))
	// These are meant to be issued from the user interface (no need to sign)
	http.Handle("PATCH /pauses/{id}", middleware(pauseController.Patch(), middlewares...))
	http.Handle("DELETE /pauses/{id}", middleware(pauseController.Delete(), middlewares...))
	logo, err := filesystem.ReadFile("assets/logo")
	if err != nil {
		return err
	}
	cli.WriteTemplate(os.Stdout, templateHeader, &Replacements{
		Logo:           string(logo),
		Version:        version,
		Url:            url,
		Copyright:      copyright,
		Name:           name,
		DisplayAddress: displayAddress,
		GeneratedKeys:  joinGeneratedKeys(generatedKeys),
	})
	if protocol == "https" {
		return http.ServeTLS(listener, nil, options.TLSCert, options.TLSPrivateKey)
	}
	return http.Serve(listener, nil)
}

func joinGeneratedKeys(keys []string) string {
	if len(keys) > 0 {
		return "\n" + strings.Join(keys, "\n\n") + "\n"
	}
	return ""
}
