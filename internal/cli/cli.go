// Package cli wires subcommands together for the iphone-mirror binary.
//
// Each subcommand is a small struct with a Name, Summary, and Run method so
// that adding one is a matter of registering it in commands().
package cli

import (
	"flag"
	"fmt"
	"io"
	"sort"
)

// Env carries process-level dependencies into commands so they stay testable.
type Env struct {
	Version string
	Stdout  io.Writer
	Stderr  io.Writer
}

// Command is the contract every subcommand satisfies.
type Command interface {
	Name() string
	Summary() string
	Run(args []string, env Env) int
}

func commands() []Command {
	return []Command{
		&devicesCmd{},
		&versionCmd{},
	}
}

// Run dispatches argv to the right subcommand. It returns a process exit code.
func Run(argv []string, env Env) int {
	if len(argv) == 0 || argv[0] == "help" || argv[0] == "-h" || argv[0] == "--help" {
		printRootHelp(env.Stdout)
		if len(argv) == 0 {
			return 1
		}
		return 0
	}

	name := argv[0]
	for _, c := range commands() {
		if c.Name() == name {
			return c.Run(argv[1:], env)
		}
	}

	fmt.Fprintf(env.Stderr, "iphone-mirror: unknown command %q\n\n", name)
	printRootHelp(env.Stderr)
	return 2
}

func printRootHelp(w io.Writer) {
	fmt.Fprintln(w, "iphone-mirror — mirror an iPhone screen to a Windows host.")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Usage:")
	fmt.Fprintln(w, "  iphone-mirror <command> [flags]")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Commands:")

	cmds := commands()
	sort.Slice(cmds, func(i, j int) bool { return cmds[i].Name() < cmds[j].Name() })
	for _, c := range cmds {
		fmt.Fprintf(w, "  %-10s %s\n", c.Name(), c.Summary())
	}

	fmt.Fprintln(w)
	fmt.Fprintln(w, "Run 'iphone-mirror <command> -h' for command-specific flags.")
}

// newFlagSet builds a FlagSet that prints its usage to env.Stderr so that
// subcommand help text follows the same convention as Run().
func newFlagSet(name string, env Env) *flag.FlagSet {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.SetOutput(env.Stderr)
	return fs
}
