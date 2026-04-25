# iOS USB QuickTime stream — implementation notes

These are working notes, not Apple documentation. The protocol is undocumented;
everything below is the consensus of community reverse engineering, summarised
here so the implementation has a single reference.

## How macOS does it

When you plug an iPhone into a Mac and open QuickTime → New Movie Recording →
pick the iPhone as the camera, macOS sends a vendor-specific USB control
transfer to the device. The device responds by **resetting itself and
re-enumerating with an additional USB configuration** that exposes:

- One bulk IN endpoint carrying H.264 video (Annex-B-ish, with Apple framing)
- One bulk IN endpoint carrying linear PCM audio
- One bulk OUT endpoint for ping / sync packets back to the host

The recording dot in the iPhone's status bar lights up while this configuration
is active. Setting the configuration back to the default unloads it.

## Activating the QuickTime configuration

The activation request, observed on the wire:

| Field           | Value                                |
| --------------- | ------------------------------------ |
| `bmRequestType` | `0x40` (Vendor, Host→Device, Device) |
| `bRequest`      | `0x52`                               |
| `wValue`        | `0x0000`                             |
| `wIndex`        | `0x0002`                             |
| `wLength`       | `0`                                  |

After this control transfer the device disappears from the bus for a few
hundred milliseconds and reappears with the extra configuration. The host has
to re-enumerate, claim the new interface, and start reading the bulk endpoints.

Deactivation uses the same request with `wIndex = 0x0000`.

## Stream framing

Each bulk read returns a stream of length-prefixed packets:

```
+---------+---------+---------+---------+---------+
| length  | magic   | payload …                   |
+---------+---------+---------+---------+---------+
  4 bytes   4 bytes   length-8 bytes
```

`length` is little-endian and includes the 8-byte header. `magic` is a 4-CC
identifying the packet type. The ones we care about:

| Magic   | Direction | Meaning                                  |
| ------- | --------- | ---------------------------------------- |
| `PING ` | host → dev | keep-alive, host pings every ~1 s       |
| `SYNC ` | dev → host | initial handshake; device asks for clock |
| `CWPA ` | dev → host | clock + audio format announcement        |
| `AFMT ` | dev → host | audio format details                     |
| `CVRP ` | dev → host | video format announcement                |
| `MEDA ` | dev → host | media data (audio or video sample)       |

Inside `MEDA`, video samples are H.264 NAL units in AVCC format (4-byte length
prefix, no start codes). Converting to Annex-B (replacing the length prefix
with `00 00 00 01`) makes the stream playable by ffplay/mpv/anything else.

## References

The two cleanest write-ups of this protocol:

- danielpaulus, *quicktime_video_hack* — Go implementation, very readable
  source: https://github.com/danielpaulus/quicktime_video_hack
- jonpryor, *Reverse engineering the iOS USB camera protocol* —
  blog post that walks through the original capture analysis
