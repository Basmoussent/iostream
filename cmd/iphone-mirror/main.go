// iphone-mirror mirrors an iPhone screen to a Windows host over USB or Wi-Fi.
package main

import (
	"os"

	"github.com/0xBasmoussent/iphone-mirror/internal/cli"
)

// version is overridden at build time via -ldflags.
var version = "dev"

func main() {
	os.Exit(cli.Run(os.Args[1:], cli.Env{
		Version: version,
		Stdout:  os.Stdout,
		Stderr:  os.Stderr,
	}))
}
