package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"os/exec"

	"github.com/Basmoussent/iostream/internal/driver"
	"github.com/Basmoussent/iostream/internal/usb"
)

// registerAPI wires the JSON endpoints the frontend calls. Endpoints stay
// chatty-style (one verb each) instead of REST-ful so the JS client can be a
// dumb fetch wrapper.
func registerAPI(mux *http.ServeMux) {
	mux.HandleFunc("/api/devices", handleDevices)
	mux.HandleFunc("/api/activate", handleActivate)
	mux.HandleFunc("/api/deactivate", handleDeactivate)
	mux.HandleFunc("/api/setup-driver", handleSetupDriver)
	mux.HandleFunc("/api/stream", handleStream)
}

func handleDevices(w http.ResponseWriter, _ *http.Request) {
	devs, err := usb.NewBackend().Discover()
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, devs)
}

func handleActivate(w http.ResponseWriter, r *http.Request) {
	udid := r.URL.Query().Get("udid")
	if err := usb.NewBackend().Activate(udid); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, map[string]bool{"ok": true})
}

func handleDeactivate(w http.ResponseWriter, r *http.Request) {
	udid := r.URL.Query().Get("udid")
	if err := usb.NewBackend().Deactivate(udid); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, map[string]bool{"ok": true})
}

func handleSetupDriver(w http.ResponseWriter, _ *http.Request) {
	if err := driver.Setup(); err != nil {
		if errors.Is(err, driver.ErrUnsupportedPlatform) {
			writeJSON(w, map[string]string{"status": "skipped", "reason": err.Error()})
			return
		}
		writeError(w, err)
		return
	}
	writeJSON(w, map[string]string{"status": "ok"})
}

// handleStream spawns `iostream stream | ffplay -` as a detached pipeline.
// We do not stream bytes through the HTTP layer because (a) we'd lose the
// hardware-accelerated decode that ffplay gives us and (b) the webview cannot
// decode raw H.264 NAL units anyway.
func handleStream(w http.ResponseWriter, r *http.Request) {
	udid := r.URL.Query().Get("udid")
	if err := launchStream(udid); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, map[string]string{"status": "started"})
}

// launchStream wires `iostream stream` into ffplay via an OS pipe and returns
// once both processes have started. Output is left on the parent's stderr so
// failures show up in the GUI's terminal log.
func launchStream(udid string) error {
	args := []string{"stream"}
	if udid != "" {
		args = append(args, "--udid", udid)
	}
	streamer := exec.Command("iostream", args...)
	player := exec.Command("ffplay",
		"-fflags", "nobuffer",
		"-flags", "low_delay",
		"-framedrop",
		"-i", "-",
	)

	pipe, err := streamer.StdoutPipe()
	if err != nil {
		return err
	}
	player.Stdin = pipe

	if err := streamer.Start(); err != nil {
		return err
	}
	if err := player.Start(); err != nil {
		_ = streamer.Process.Kill()
		return err
	}
	// Reap children in the background so they don't turn into zombies. We do
	// not wait on either: the user controls them via their own UI / Ctrl-C.
	go func() { _ = streamer.Wait() }()
	go func() { _ = player.Wait() }()
	return nil
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}
