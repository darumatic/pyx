package main

import (
	"devcli/cmd"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func Init() {
	repository := cmd.Repository()
	if len(repository) == 0 {
		fmt.Println("dev is not initialized, install default public script repository")
		cmd.Install("https://github.com/darumatic/dev-cli-scripts.git")
		fmt.Println("\n")
	}
}

func main() {
	Init()

	cmdDev := cmd.MakeDev()
	cmdDev.AddCommand(cmd.MakeInstall())
	cmdDev.AddCommand(cmd.MakeUninstall())
	cmdDev.AddCommand(cmd.MakeVersion())
	cmdDev.AddCommand(cmd.MakeUpdate())
	cmdDev.AddCommand(cmd.MakeHelp())
	cmdDev.AddCommand(cmd.MakeList())

	cmdDev.SetHelpCommand(&cobra.Command{
		Use:    "no-help",
		Hidden: true,
	})
	if err := cmdDev.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
