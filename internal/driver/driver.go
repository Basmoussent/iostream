// Package driver installs the WinUSB driver Apple devices need before libusb
// can talk to them on Windows.
//
// The actual install is delegated to Zadig (https://zadig.akeo.ie/), which is
// a pre-signed, well-trusted UI on top of libwdi. We download it once, launch
// it elevated, and the user clicks "Replace Driver" — which is one click
// instead of the five-step manual flow documented in docs/WINDOWS.md.
//
// Bundling libwdi directly to drop the last click would require either a cgo
// binding (libwdi headers are not in vcpkg) or compiling wdi-simple from
// source in CI. Both are larger projects than this; see ROADMAP M6.
package driver

import "errors"

// ErrUnsupportedPlatform is returned when Setup is called on a platform that
// does not need a driver swap (Linux, macOS — libusb works out of the box).
var ErrUnsupportedPlatform = errors.New("driver: install only required on Windows")

// Setup downloads Zadig if necessary and launches it elevated so the user can
// bind the iPhone interface to WinUSB. The platform-specific implementation
// lives in driver_windows.go; on every other platform this returns
// ErrUnsupportedPlatform.
func Setup() error {
	return setup()
}
