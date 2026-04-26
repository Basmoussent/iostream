//go:build !windows

package driver

func setup() error { return ErrUnsupportedPlatform }
