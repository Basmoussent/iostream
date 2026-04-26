# iostream

Mirror an iPhone screen to Windows over USB or Wi-Fi, in the highest quality the device can produce.

> **Status:** USB pipeline (M1–M3) is wired up — `devices`, `activate`, `deactivate`, `stream`, and `setup-driver` work. A small WebView GUI (`iostream-gui`) ships alongside. AirPlay (M5) and an embedded decoder (M4 v2) are next. See [`docs/ROADMAP.md`](docs/ROADMAP.md).

## Goals

- **Best-in-class image quality.** Capture the native H.264 stream Apple already produces — no re-encoding, no scaling.
- **Low latency.** Sub-100 ms over USB; AirPlay-grade over Wi-Fi.
- **Single self-contained binary.** No browser, no Java runtime, no Electron.
- **Windows first**, but the core is portable to Linux and macOS.

## Transports

| Transport     | Quality      | Latency   | Setup                               | Status     |
| ------------- | ------------ | --------- | ----------------------------------- | ---------- |
| USB / Lightning / USB-C | 1080p60, native H.264 | ~50 ms | Driver swap (Zadig, one-time)       | in progress |
| Wi-Fi / AirPlay         | 1080p, AAC audio      | ~150 ms | None (same network as the phone)    | planned    |

USB uses Apple's hidden "QuickTime" USB configuration — the same mechanism QuickTime on macOS uses when you plug in an iPhone and pick it as a camera source. Wi-Fi uses a reverse-engineered AirPlay 2 mirroring receiver.

## Quick start

### GUI (recommended)

Double-click `iostream-gui.exe`. Three cards:

1. **Install / repair driver** — auto-downloads Zadig and elevates it. One click in Zadig (`WinUSB → Replace Driver`) and you're done forever.
2. **Devices** — refreshes the connected iPhone list.
3. **Start mirroring** — launches `ffplay` in a separate window with the live stream. `ffplay` must be on `PATH` (`winget install Gyan.FFmpeg`).

### CLI

```powershell
# One-time WinUSB driver install — opens Zadig with admin elevation.
iostream setup-driver

# List connected iPhones
iostream devices

# Pipe the raw H.264 to ffplay / mpv for absolute minimum latency
iostream stream | ffplay -fflags nobuffer -flags low_delay -framedrop -

# Or save a recording to disk (no re-encoding — just the device's native H.264)
iostream stream -o recording.h264

# Manage the QuickTime USB configuration explicitly (stream auto-activates)
iostream activate
iostream deactivate
```

## Building from source

Requires Go 1.23+. The libusb backend is gated behind `CGO_ENABLED=1` so the
default build is pure Go and works everywhere — useful for CI and for reading
the source. To actually talk to a phone you need a cgo build with libusb.

```sh
# Pure-Go build (stub backend — useful for development/CI lint):
make build

# Real build with libusb (talks to the phone):
make cgo

# Cross-compile from WSL/Linux to Windows:
make windows-cgo   # needs mingw-w64 + libusb headers under MINGW
make windows       # stub-only build, no cgo

# Per-platform variants:
make linux-cgo darwin-cgo
```

On Windows, install libusb via [vcpkg](https://vcpkg.io/) (`vcpkg install libusb:x64-windows`)
or grab the prebuilt `iostream-windows` artifact from CI. The CI workflow
in `.github/workflows/ci.yml` is the canonical reference for how to set up
`PKG_CONFIG_PATH` and `CGO_LDFLAGS` for a Windows cgo build.

## Repository layout

```
cmd/iostream/         CLI entry point
cmd/iostream-gui/     WebView GUI (cgo)
internal/cli/         Subcommand wiring
internal/usb/         USB device discovery + QuickTime protocol (cgo + libusb)
internal/driver/      Zadig auto-installer for the WinUSB driver swap
internal/airplay/     AirPlay receiver (planned)
docs/                 Protocol notes, Windows setup, roadmap
```

## Prior art and credits

This project stands on the shoulders of years of community reverse engineering:

- [`quicktime_video_hack`](https://github.com/danielpaulus/quicktime_video_hack) — the canonical Go implementation of the iOS USB QuickTime stream protocol.
- [`libimobiledevice`](https://libimobiledevice.org/) — for everything pairing- and lockdown-related.
- [`UxPlay`](https://github.com/FDH2/UxPlay) and [`RPiPlay`](https://github.com/FD-/RPiPlay) — open AirPlay 2 mirroring receivers.

Where we depend on or borrow from these, attribution lives in [`docs/CREDITS.md`](docs/CREDITS.md).

## License

[MIT](LICENSE).
