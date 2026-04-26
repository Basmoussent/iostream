// iostream-gui is a small WebView-based front end for iostream.
//
// It serves an embedded HTML/JS bundle from a localhost HTTP server, then
// opens a native window pointing at it (Microsoft.Web.WebView2 on Windows,
// WebKit2GTK on Linux, WebKit on macOS). The frontend talks to a tiny JSON
// API (see api.go) which wraps internal/usb and internal/driver.
//
// Streaming hands off to ffplay/mpv as a child process — embedding a video
// decoder in the webview is M4 v2 (see docs/ROADMAP.md).
package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	webview "github.com/webview/webview_go"
)

func main() {
	srv, addr, err := startAPIServer()
	if err != nil {
		fmt.Fprintf(os.Stderr, "iostream-gui: %v\n", err)
		os.Exit(1)
	}

	w := webview.New(false)
	defer w.Destroy()
	w.SetTitle("iostream")
	w.SetSize(960, 640, webview.HintNone)
	w.Navigate("http://" + addr)
	w.Run()

	// Webview returned → tear the API server down so the process exits cleanly.
	_ = srv.Close()
}

// startAPIServer binds an HTTP server to a random localhost port. We let the
// kernel pick the port (port 0) so two GUI instances can coexist and we never
// collide with anything the user happens to be running.
func startAPIServer() (*http.Server, string, error) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, "", fmt.Errorf("listen: %w", err)
	}
	mux := http.NewServeMux()
	registerAPI(mux)
	registerAssets(mux)

	srv := &http.Server{Handler: mux}
	go func() {
		if err := srv.Serve(ln); err != nil && err != http.ErrServerClosed {
			log.Printf("iostream-gui: api server: %v", err)
		}
	}()
	return srv, ln.Addr().String(), nil
}
