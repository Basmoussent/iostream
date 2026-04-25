# iphone-mirror

Mirror an iPhone screen to Windows over USB or Wi-Fi, in the highest quality the device can produce.

> **Status:** early development. The CLI scaffold and device discovery work; the QuickTime stream pipeline is being built. See [`docs/ROADMAP.md`](docs/ROADMAP.md).

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

## Quick start (USB, when ready)

```powershell
# One-time: swap the iPhone USB driver to WinUSB using Zadig.
# See docs/WINDOWS.md for screenshots.

# List connected iPhones
iphone-mirror devices

# Stream to a built-in window
iphone-mirror play

# Or pipe the raw H.264 to ffplay / mpv for absolute minimum latency
iphone-mirror stream | ffplay -fflags nobuffer -flags low_delay -framedrop -
```

## Building from source

Requires Go 1.23+.

```sh
# From WSL or Linux, cross-compile a Windows binary:
make windows

# Native build for whatever you're on:
make build
```

The Windows build does not yet require cgo — the libusb dependency is added in a later milestone (see roadmap).

## Repository layout

```
cmd/iphone-mirror/    Entry point
internal/cli/         Subcommand wiring
internal/usb/         USB device discovery and QuickTime protocol (planned)
internal/airplay/     AirPlay receiver (planned)
internal/decoder/     H.264 NAL handling (planned)
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
