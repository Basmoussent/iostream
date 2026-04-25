# Roadmap

The work is broken into small milestones so each step is shippable on its own.
Tick boxes are updated as code lands.

## M0 — Project foundation ✅

- [x] LICENSE, README, .gitignore, .editorconfig
- [x] Go module + Makefile (native, Windows, Linux, Darwin builds)
- [x] CLI dispatcher with `version` and `devices` subcommands
- [x] `internal/usb.Discoverer` interface and stub backend

## M1 — Real USB device discovery

- [ ] Add `gousb` (libusb) dependency, gated behind a `usb_libusb` build tag
- [ ] Implement `libusbDiscoverer` that walks attached devices and filters by
      Apple's vendor ID (`0x05ac`)
- [ ] Detect whether the QuickTime USB configuration is currently active
- [ ] Read product strings (USB string descriptor `iProduct`)
- [ ] Document the Zadig driver swap in [`WINDOWS.md`](WINDOWS.md)

**Done when:** `iphone-mirror devices` prints a real iPhone on Windows.

## M2 — QuickTime configuration swap

- [ ] Send the magic control transfer that flips the device into the hidden
      "QuickTime" USB configuration (see [`PROTOCOL.md`](PROTOCOL.md))
- [ ] Detect the reset/re-enumeration and re-open the device on its new config
- [ ] Add `iphone-mirror activate <udid>` and `deactivate <udid>` commands

**Done when:** the recording dot appears in the iPhone's status bar.

## M3 — Stream consumption

- [ ] Implement the QuickTime sync packet handshake (`PING`, `SYNC`, `CWPA`, `AFMT`, `CVRP`)
- [ ] Bulk-read the H.264 + AAC stream from the QuickTime endpoint
- [ ] Strip Apple framing and emit clean Annex-B NAL units
- [ ] `iphone-mirror stream` writes the H.264 elementary stream to stdout

**Done when:** `iphone-mirror stream | ffplay -` shows a live picture.

## M4 — Built-in player

- [ ] Embed a video sink so `iphone-mirror play` works without an external
      player. Candidates: libmpv, GStreamer, Wails + `<video>` MSE.
- [ ] Decide on the UI shell (likely Wails — webview2 ships on Win10/11)
- [ ] Audio playback via WASAPI

**Done when:** double-clicking `iphone-mirror.exe` shows the iPhone.

## M5 — Wi-Fi / AirPlay

- [ ] Implement an AirPlay 2 mirroring receiver (RTSP + FairPlay handshake)
- [ ] Bonjour/mDNS advertisement so the iPhone discovers the host
- [ ] Reuse the M4 player to render the AirPlay stream

**Done when:** the iPhone's AirPlay menu lists the Windows host and casts to it.

## M6 — Polish

- [ ] Code signing for the Windows binary
- [ ] Auto-update channel
- [ ] Recording to MP4 (mux H.264 + AAC, no re-encode)
- [ ] Per-device profiles (orientation lock, audio routing)
