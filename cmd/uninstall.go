package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func Uninstall(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println(cmd.Usage())
		os.Exit(1)
	}

	dir := AppHomeDir() + "/" + args[0]
	if DirExists(dir) {
		os.RemoveAll(dir)
	}
}

func MakeUninstall() *cobra.Command {
	var command = &cobra.Command{
		Use:     "uninstall",
		Short:   "uninstall script git repository",
		Long:    `uninstall script git repository`,
		Example: `  dev uninstall dir`,
	}

	command.RunE = func(command *cobra.Command, args []string) error {
		Uninstall(command, args)
		return nil
	}

	return command
}
