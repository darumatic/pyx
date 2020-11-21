package main

import (
	"devcli/cmd"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func main() {
	cmdDev := cmd.MakeDev()
	cmdDev.AddCommand(cmd.MakeInstall())
	cmdDev.AddCommand(cmd.MakeUninstall())
	cmdDev.AddCommand(cmd.MakeVersion())
	cmdDev.AddCommand(cmd.MakeUpdate())
	cmdDev.AddCommand(cmd.MakeHelp())
	cmdDev.AddCommand(cmd.MakeList())
	cmdDev.AddCommand(cmd.MakePython())

	cmdDev.SetHelpCommand(&cobra.Command{
		Use:    "no-help",
		Hidden: true,
	})
	if err := cmdDev.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
