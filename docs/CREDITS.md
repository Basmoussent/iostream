# Credits

iphone-mirror builds on years of community work to reverse-engineer Apple's
iOS USB and AirPlay protocols. Every section below is something we depend on,
were inspired by, or used to learn the protocol from.

## Direct dependencies

_None yet — the project is pure stdlib at M0._

When dependencies are added their licenses will be enumerated here.

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
