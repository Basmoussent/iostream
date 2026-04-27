package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

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
	logPath, err := launchStream(udid)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, map[string]string{
		"status": "started",
		"log":    logPath,
	})
}

// launchStream wires `iostream stream` into ffplay via an OS pipe. It returns
// the path of the log file we attach to both children's stderr so the
// frontend can show the user where to look when something dies.
func launchStream(udid string) (string, error) {
	args := []string{"stream"}
	if udid != "" {
		args = append(args, "--udid", udid)
	}
	streamerPath, err := resolveSibling("iostream")
	if err != nil {
		return "", err
	}

	logFile, err := openLogFile()
	if err != nil {
		return "", err
	}

	streamer := exec.Command(streamerPath, args...)
	player := exec.Command("ffplay",
		"-fflags", "nobuffer",
		"-flags", "low_delay",
		"-framedrop",
		"-i", "-",
	)

	streamer.Stderr = logFile
	player.Stderr = logFile
	hideConsole(streamer)
	hideConsole(player)

	pipe, err := streamer.StdoutPipe()
	if err != nil {
		_ = logFile.Close()
		return "", err
	}
	player.Stdin = pipe

	if err := streamer.Start(); err != nil {
		_ = logFile.Close()
		return "", fmt.Errorf("start iostream: %w", err)
	}
	if err := player.Start(); err != nil {
		_ = streamer.Process.Kill()
		_ = logFile.Close()
		return "", fmt.Errorf("start ffplay: %w", err)
	}
	// Reap children in the background so they don't turn into zombies, and
	// close the log file once the slower of the two has exited.
	go func() {
		_ = streamer.Wait()
		_ = player.Wait()
		_ = logFile.Close()
	}()
	return logFile.Name(), nil
}

// openLogFile creates a fresh log file under the user's temp dir for each
// stream attempt. Caller owns the file; caller closes it.
func openLogFile() (*os.File, error) {
	dir := filepath.Join(os.TempDir(), "iostream")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	return os.CreateTemp(dir, "stream-*.log")
}

// resolveSibling returns the absolute path of `name` if it lives next to the
// running GUI binary, falling back to a regular PATH lookup. Go 1.19 stopped
// resolving relative paths through exec.LookPath for security, so we have to
// do this ourselves to make `iostream-gui.exe` find its sibling `iostream.exe`.
func resolveSibling(name string) (string, error) {
	exeName := name
	if runtime.GOOS == "windows" && filepath.Ext(exeName) == "" {
		exeName += ".exe"
	}
	if self, err := os.Executable(); err == nil {
		candidate := filepath.Join(filepath.Dir(self), exeName)
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
	}
	return exec.LookPath(name)
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
