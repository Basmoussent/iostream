//go:build windows

package main

import (
	"os/exec"
	"syscall"
)

// hideConsole asks Windows not to allocate a console window for this child
// process. Without it, every spawned console-subsystem child (iostream.exe,
// ffplay.exe) flashes a black terminal because the GUI is a windowsgui app
// and has no console of its own to inherit.
//
// 0x08000000 is CREATE_NO_WINDOW. HideWindow toggles SW_HIDE for any window
// the process tries to create itself.
func hideConsole(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: 0x08000000,
	}
}
