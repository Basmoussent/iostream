package cli

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/0xBasmoussent/iphone-mirror/internal/usb"
)

type streamCmd struct{}

func (streamCmd) Name() string    { return "stream" }
func (streamCmd) Summary() string { return "Stream the iPhone screen as raw H.264 to stdout." }

func (streamCmd) Run(args []string, env Env) int {
	fs := newFlagSet("stream", env)
	udid := fs.String("udid", "", "device UDID (USB serial). Empty = first device found.")
	output := fs.String("o", "-", "video output path. '-' means stdout.")
	audioPath := fs.String("audio", "", "if set, write linear PCM audio samples to this path.")
	noActivate := fs.Bool("no-activate", false, "do not auto-activate the QuickTime config first.")
	if cont, code := parseFlags(fs, args); !cont {
		return code
	}

	video, closeVideo, err := openOutput(*output, env.Stdout)
	if err != nil {
		fmt.Fprintf(env.Stderr, "iphone-mirror: %v\n", err)
		return 1
	}
	defer closeVideo()

	var audio io.Writer
	var closeAudio func()
	if *audioPath != "" {
		audio, closeAudio, err = openOutput(*audioPath, env.Stdout)
		if err != nil {
			fmt.Fprintf(env.Stderr, "iphone-mirror: %v\n", err)
			return 1
		}
		defer closeAudio()
	}

	stop := make(chan struct{})
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigs
		close(stop)
	}()
	defer signal.Stop(sigs)

	fmt.Fprintln(env.Stderr, "iphone-mirror: streaming — Ctrl-C to stop.")

	err = usb.NewBackend().Stream(*udid, usb.StreamOptions{
		Video:        video,
		Audio:        audio,
		Stop:         stop,
		AutoActivate: !*noActivate,
	})
	if err != nil {
		fmt.Fprintf(env.Stderr, "iphone-mirror: %v\n", err)
		return 1
	}
	return 0
}

// openOutput resolves a CLI output path. "-" is special-cased to mean stdout.
// The returned closer is always safe to call, even for stdout (it is a no-op).
func openOutput(path string, stdout io.Writer) (io.Writer, func(), error) {
	if path == "-" {
		return stdout, func() {}, nil
	}
	f, err := os.Create(path)
	if err != nil {
		return nil, nil, fmt.Errorf("open %s: %w", path, err)
	}
	return f, func() { _ = f.Close() }, nil
}
