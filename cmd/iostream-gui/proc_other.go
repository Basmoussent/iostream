//go:build !windows

package main

import "os/exec"

// hideConsole is a no-op on platforms that do not allocate console windows
// for child processes (i.e. everywhere that isn't Windows).
func hideConsole(_ *exec.Cmd) {}
