# Roadmap

The work is broken into small milestones so each step is shippable on its own.
Tick boxes are updated as code lands.

## M0 — Project foundation ✅

- [x] LICENSE, README, .gitignore, .editorconfig
- [x] Go module + Makefile (native, Windows, Linux, Darwin builds)
- [x] CLI dispatcher with `version` and `devices` subcommands
- [x] `internal/usb.Discoverer` interface and stub backend

## M1 — Real USB device discovery ✅

- [x] Add `quicktime_video_hack` (and its `gousb`/libusb dependency) behind a
      `cgo` build tag, with a no-op stub on `!cgo` for portability
- [x] `libusbBackend.Discover()` walks attached Apple devices and reports
      whether the QuickTime USB configuration is currently active
- [x] Read product strings (via QVH's `IosDevice.ProductName`)
- [x] Document the Zadig driver swap in [`WINDOWS.md`](WINDOWS.md)

**Done when:** `iostream devices` prints a real iPhone on Windows.

## M2 — QuickTime configuration swap ✅

- [x] Send the control transfer that flips the device into the hidden
      "QuickTime" USB configuration (see [`PROTOCOL.md`](PROTOCOL.md))
- [x] Detect the reset/re-enumeration and re-open the device on its new config
- [x] `iostream activate [--udid X]` and `deactivate [--udid X]`

**Done when:** the recording dot appears in the iPhone's status bar.

## M3 — Stream consumption ✅

- [x] Implement the QuickTime sync packet handshake (`PING`, `SYNC`, `CWPA`, `AFMT`, `CVRP`)
- [x] Bulk-read the H.264 + AAC stream from the QuickTime endpoint
- [x] Strip Apple framing and emit clean Annex-B NAL units
- [x] `iostream stream` writes the H.264 elementary stream to stdout, with
      `--audio` for an optional PCM track and SIGINT-driven clean shutdown

## M4 — Built-in player

### M4a — GUI shell ✅

- [x] WebView-based GUI (`iostream-gui`) using `webview_go`; embedded
      vanilla HTML/CSS/JS — no Node toolchain
- [x] Three-card UI: driver install, device list, start mirroring
- [x] Auto-downloads Zadig and launches it elevated for the driver step
      (`iostream setup-driver` exposes the same flow on the CLI)
- [x] Stream button spawns `iostream stream | ffplay -` so the user gets a
      working window immediately

### M4b — Embedded decoder

- [ ] Stream the H.264 NAL units to the webview over a localhost endpoint
      and decode in-page via WebCodecs (Chromium-only — fine, WebView2 is
      Chromium)
- [ ] Audio playback via the Web Audio API or WASAPI passthrough
- [ ] Replace the ffplay subprocess so the mirror lives inside the GUI

**Done when:** double-clicking `iostream-gui.exe` shows the iPhone in-window
with no external player required.

## M5 — Wi-Fi / AirPlay

- [ ] Implement an AirPlay 2 mirroring receiver (RTSP + FairPlay handshake)
- [ ] Bonjour/mDNS advertisement so the iPhone discovers the host
- [ ] Reuse the M4 player to render the AirPlay stream

**Done when:** the iPhone's AirPlay menu lists the Windows host and casts to it.

## M6 — Polish

- [ ] Code signing for the Windows binaries
- [ ] Auto-update channel
- [ ] Recording to MP4 (mux H.264 + AAC, no re-encode)
- [ ] Per-device profiles (orientation lock, audio routing)
- [ ] True zero-click driver install via `libwdi` cgo binding (skips the
      one Replace-Driver click in Zadig that `setup-driver` still requires)
- [ ] Pin Zadig SHA-256 in `internal/driver/driver_windows.go` so the
      download is verified before launch
- [ ] Vendor and patch `webview_go` so the GUI builds on CI: upstream pins
      to `webkit2gtk-4.0` (gone from Ubuntu 24.04) and emits `-mthreads`
      which mingw-w64 14+ rejects
