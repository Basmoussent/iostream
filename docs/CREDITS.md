# Credits

iostream builds on years of community work to reverse-engineer Apple's
iOS USB and AirPlay protocols. Every section below is something we depend on,
were inspired by, or used to learn the protocol from.

## Direct dependencies

- **github.com/danielpaulus/quicktime_video_hack** — MIT
  Powers the entire USB pipeline (M1–M3): device discovery, QuickTime config
  activation, the message processor, and `coremedia.AVFileWriter` which
  emits H.264 in Annex-B framing. This dependency is what makes the project
  feasible without months of solo reverse engineering.
  Patched via `go.mod replace` to a fork
  (https://github.com/Basmoussent/quicktime_video_hack) — newer iPhones
  (iPhone 14 / 15 / 16) renumbered the QuickTime USB subclass from 0x2A
  to 0xFD; the fork carries the one-line constant change.
- **github.com/google/gousb** — Apache-2.0 (transitive, via QVH)
  Cgo bindings around libusb-1.0.
- **libusb-1.0** — LGPL-2.1 (system library, dynamically linked on cgo builds)
- **github.com/webview/webview_go** — MIT
  Cgo bindings around the platform webview (WebView2 on Windows, WebKit2GTK
  on Linux, WebKit on macOS). Powers `iostream-gui` without a Node toolchain.
- **Zadig** — GPL-3.0
  Downloaded on demand by `iostream setup-driver` to perform the WinUSB
  driver swap. Not redistributed; fetched from the official release.

## Inspiration / reference implementations

- **quicktime_video_hack** — Daniel Paulus
  <https://github.com/danielpaulus/quicktime_video_hack> — MIT
  Reference Go implementation of the iOS USB QuickTime stream protocol. Our
  M2/M3 milestones follow its packet layout closely; see
  [`PROTOCOL.md`](PROTOCOL.md).

- **libimobiledevice** — Nikias Bassen and contributors
  <https://libimobiledevice.org/> — LGPL-2.1
  The canonical implementation of usbmuxd, lockdownd, and the rest of the iOS
  pairing stack. Anything we do that involves UDIDs or device info ultimately
  traces back here.

- **UxPlay** — Florian Draschbacher / FDH2
  <https://github.com/FDH2/UxPlay> — GPL-3.0
  Open-source AirPlay 2 mirroring receiver. The M5 Wi-Fi transport will draw
  on its handling of the FairPlay handshake and RTSP setup.

- **RPiPlay** — FD-
  <https://github.com/FD-/RPiPlay> — GPL-3.0
  The original Raspberry Pi AirPlay receiver UxPlay forked from.

If you spot a project we're using or borrowing from that isn't listed,
please open a PR.
