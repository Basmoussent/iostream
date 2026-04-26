//go:build cgo

package usb

import (
	"fmt"
	"io"
	"strings"

	qvh "github.com/danielpaulus/quicktime_video_hack/screencapture"
	"github.com/danielpaulus/quicktime_video_hack/screencapture/coremedia"
)

// libusbBackend is the production backend. It speaks to the iPhone over raw
// USB through gousb (libusb) and implements Apple's QuickTime stream protocol
// via the quicktime_video_hack library.
type libusbBackend struct{}

func newBackend() Backend { return libusbBackend{} }

func (libusbBackend) Discover() ([]Device, error) {
	list, err := qvh.FindIosDevices()
	if err != nil {
		return nil, fmt.Errorf("usb discover: %w", err)
	}
	out := make([]Device, len(list))
	for i := range list {
		d := list[i]
		out[i] = Device{
			UDID:             d.SerialNumber,
			ProductName:      d.ProductName,
			USBInfo:          d.UsbInfo,
			QuickTimeEnabled: d.IsActivated(),
		}
	}
	return out, nil
}

func (libusbBackend) Activate(udid string) error {
	d, err := qvh.FindIosDevice(udid)
	if err != nil {
		return fmt.Errorf("usb activate: %w", wrapNotFound(err))
	}
	if _, err := qvh.EnableQTConfig(d); err != nil {
		return fmt.Errorf("usb activate: %w", err)
	}
	return nil
}

func (libusbBackend) Deactivate(udid string) error {
	d, err := qvh.FindIosDevice(udid)
	if err != nil {
		return fmt.Errorf("usb deactivate: %w", wrapNotFound(err))
	}
	if _, err := qvh.DisableQTConfig(d); err != nil {
		return fmt.Errorf("usb deactivate: %w", err)
	}
	return nil
}

func (libusbBackend) Stream(udid string, opts StreamOptions) error {
	d, err := qvh.FindIosDevice(udid)
	if err != nil {
		return fmt.Errorf("usb stream: %w", wrapNotFound(err))
	}

	if opts.AutoActivate && !d.IsActivated() {
		d, err = qvh.EnableQTConfig(d)
		if err != nil {
			return fmt.Errorf("usb stream: activate: %w", err)
		}
	}

	video := opts.Video
	if video == nil {
		video = io.Discard
	}
	audio := opts.Audio
	if audio == nil {
		audio = io.Discard
	}
	audioOnly := opts.Video == nil && opts.Audio != nil
	consumer := coremedia.NewAVFileWriter(video, audio)

	adapter := qvh.UsbAdapter{}
	stop := make(chan interface{})
	if opts.Stop != nil {
		go bridgeStop(opts.Stop, stop)
	}

	mp := qvh.NewMessageProcessor(&adapter, stop, consumer, audioOnly)
	if err := adapter.StartReading(d, &mp, stop); err != nil {
		return fmt.Errorf("usb stream: %w", err)
	}
	return nil
}

// bridgeStop forwards a close on the typed Stop channel to the untyped one
// quicktime_video_hack expects. It uses select+default so a late stop after
// StartReading already returned does not block the goroutine forever.
func bridgeStop(in <-chan struct{}, out chan<- interface{}) {
	<-in
	select {
	case out <- struct{}{}:
	default:
	}
}

// wrapNotFound converts QVH's "no device found" errors into ErrDeviceNotFound
// so callers can match on it without depending on QVH's error text.
func wrapNotFound(err error) error {
	if err == nil {
		return nil
	}
	// QVH reports missing devices via fmt.Errorf with messages like
	// "no device found to activate" or "Unable to find device:…". Cheap text
	// match is fine — the upstream errors are not typed.
	msg := err.Error()
	if strings.Contains(msg, "no device found") || strings.Contains(msg, "Unable to find device") {
		return fmt.Errorf("%w: %s", ErrDeviceNotFound, msg)
	}
	return err
}
