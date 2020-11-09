package cmd

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
	"path/filepath"
)

func MakeList() *cobra.Command {
	var command = &cobra.Command{
		Use:   "list",
		Short: "list",
		RunE: func(cmdc *cobra.Command, args []string) error {
			for _, f := range Repository() {
				r, err := git.PlainOpen(AppHomeDir() + "/" + f)
				if err != nil {
					fmt.Println(filepath.Base(f), "-", err.Error())
				} else {
					remote, err := r.Remote("origin")
					if err != nil {
						fmt.Println(filepath.Base(f), "-", err.Error())
					} else {
						fmt.Println(filepath.Base(f), "-", remote.Config().URLs[0])
					}
				}
			}
			return nil
		},
	}

	return command
}
