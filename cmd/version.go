package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func MakeVersion() *cobra.Command {
	var command = &cobra.Command{
		Use:          "version",
		Short:        "pyx version",
		Example:      `  pyx version`,
		SilenceUsage: false,
	}
	command.Run = func(cmd *cobra.Command, args []string) {
		if len(Version) == 0 {
			fmt.Println("0.1")
		} else {
			fmt.Println(Version)
		}
	}
	return command
}
