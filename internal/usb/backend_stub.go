//go:build !cgo

package usb

type stubBackend struct{}

func newBackend() Backend { return stubBackend{} }

func (stubBackend) Discover() ([]Device, error)        { return nil, ErrBackendUnavailable }
func (stubBackend) Activate(string) error              { return ErrBackendUnavailable }
func (stubBackend) Deactivate(string) error            { return ErrBackendUnavailable }
func (stubBackend) Stream(string, StreamOptions) error { return ErrBackendUnavailable }
