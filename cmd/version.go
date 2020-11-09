package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func MakeVersion() *cobra.Command {
	var command = &cobra.Command{
		Use:          "version",
		Short:        "dev version",
		Example:      `  dev version`,
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
