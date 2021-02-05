package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"pyx/cmd"
)

func main() {
	pyx := cmd.MakePyx()
	pyx.AddCommand(cmd.MakeVersion())

	pyx.SetHelpCommand(&cobra.Command{
		Use:    "no-help",
		Hidden: true,
	})
	if err := pyx.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
