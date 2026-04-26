// Package usb describes iOS devices reachable over USB and the contract that
// any platform-specific backend must satisfy.
//
// See backend.go for the Backend interface and the cgo / non-cgo split.
package usb

// Device is one connected iOS device as seen from the host's USB stack.
//
// UDID is the device's USB serial number — for older iPhones it is the full
// 40-character UDID, for newer ones it is a 24-character prefix that uniquely
// identifies the device on this host.
//
// QuickTimeEnabled reports whether the hidden QuickTime USB configuration is
// currently active. When true the H.264 video stream is exposed and Stream()
// can be called without activating first.
type Device struct {
	UDID             string `json:"udid"`
	ProductName      string `json:"product_name"`
	USBInfo          string `json:"usb_info"`
	QuickTimeEnabled bool   `json:"quicktime_enabled"`
}
