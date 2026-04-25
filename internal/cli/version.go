package cli

import "fmt"

type versionCmd struct{}

func (versionCmd) Name() string    { return "version" }
func (versionCmd) Summary() string { return "Print the iphone-mirror version." }

func (versionCmd) Run(_ []string, env Env) int {
	fmt.Fprintln(env.Stdout, env.Version)
	return 0
}
