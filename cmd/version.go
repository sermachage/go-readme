package cmd

import (
	"fmt"
	"runtime/debug"

	"github.com/spf13/cobra"
)

// Version is set at build time via -ldflags.
var Version = "dev"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the go-readme version",
	Run: func(cmd *cobra.Command, _ []string) {
		version := Version
		if version == "dev" {
			if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
				version = info.Main.Version
			}
		}
		fmt.Fprintf(cmd.OutOrStdout(), "go-readme %s\n", version)
	},
}
