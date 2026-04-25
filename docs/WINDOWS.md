# Windows setup

iphone-mirror talks to the iPhone over raw USB through libusb. Windows ships
without a generic USB driver, so the iPhone needs to be re-bound to **WinUSB**
once. After that it stays put — you only redo this if you reset the device
manager.

## 1. Install Zadig

Zadig is the standard tool for swapping Windows USB drivers.

- Download from <https://zadig.akeo.ie/>
- Run it as **Administrator** (right-click → Run as administrator)

## 2. Plug the iPhone in and trust the computer

Unlock the phone, plug it into the host, tap **Trust** when prompted. Leave
iTunes / Apple Devices / Apple Mobile Device Service running — they handle
pairing and won't conflict with what we're about to do.

## 3. Show the QuickTime interface in Zadig

By default Zadig hides composite-device interfaces. Enable:

- **Options → List All Devices**
- **Options → Ignore Hubs or Composite Parents** (leave unchecked)

The iPhone shows up multiple times because it's a USB composite device. Pick
the entry whose USB ID is `05AC` (Apple) and whose interface description
contains **"Apple Mobile Device USB Composite"** or **"USB Composite Device"** —
the QuickTime configuration is exposed as one of its interfaces.

## 4. Replace the driver with WinUSB

In the right-hand dropdown, choose **WinUSB** and click **Replace Driver** (or
**Install Driver** if there is no current driver). This takes a few seconds.

> **What this does:** it tells Windows "for this specific interface of this
> specific device, route I/O through libusb instead of the Apple driver."
> The rest of the iPhone (sync, file transfer, charging) keeps working.

## 5. Verify

```powershell
iphone-mirror devices
```

Once the libusb backend lands, this should print a line for the connected
iPhone. Until then it prints `device discovery not yet wired up` — see
[`ROADMAP.md`](ROADMAP.md) M1.

## Reverting

If you want to put the phone back to "stock" (e.g. for iTunes-only use):

- Open **Device Manager**
- Find the iPhone interface that now shows up under **libusb (WinUSB) devices**
- Right-click → **Uninstall device** → check **Delete the driver software**
- Unplug, replug. Windows installs the Apple driver again.
