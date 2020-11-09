package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	gitssh "github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

func MakeUpdate() *cobra.Command {
	var command = &cobra.Command{
		Use:          "update",
		Short:        "update script repositories",
		Long:         `update script repositories`,
		Example:      `  dev update`,
		SilenceUsage: false,
	}

	command.RunE = func(command *cobra.Command, args []string) error {
		path := AppHomeDir()

		for _, f := range Repository() {
			r, err := git.PlainOpen(path + "/" + f)
			if err != nil {
				fmt.Println(filepath.Base(f), "-", err.Error())
				continue
			}

			w, err := r.Worktree()
			if err != nil {
				fmt.Println(filepath.Base(f), "-", err.Error())
				continue
			}
			remote, _ := r.Remote("origin")
			cfg := remote.Config()

			if strings.HasPrefix(cfg.URLs[0], "git@") {
				s := fmt.Sprintf("%s/.ssh/id_rsa", UserHomeDir())
				sshKey, err := ioutil.ReadFile(s)
				if err != nil {
					fmt.Print(err)
				}
				var signer ssh.Signer
				var errP error

				signer, errP = ssh.ParsePrivateKey(sshKey)
				_, cerr := errP.(*ssh.PassphraseMissingError)
				if cerr {
					fmt.Printf("\nPlease type ssh password: ")
					STDIN := int(os.Stdin.Fd())
					password, _ := terminal.ReadPassword(STDIN)
					signer, _ = ssh.ParsePrivateKeyWithPassphrase(sshKey, password)
				}

				auth := &gitssh.PublicKeys{
					User:   "git",
					Signer: signer,
					HostKeyCallbackHelper: gitssh.HostKeyCallbackHelper{
						HostKeyCallback: ssh.InsecureIgnoreHostKey(),
					},
				}

				err = w.Pull(&git.PullOptions{
					RemoteName: "origin",
					Auth:       auth,
				})

				if err != nil {
					fmt.Println(filepath.Base(f), "-", err.Error())
				} else {
					fmt.Println(filepath.Base(f), "- Updated")

					err := PythonInstall(path + "/" + f)
					if err != nil {
						fmt.Printf("failed to run setup.py, %s", err.Error())
						os.Exit(1)
					}
				}
			} else {
				err = w.Pull(&git.PullOptions{RemoteName: "origin"})
				if err != nil {
					fmt.Println(filepath.Base(f), "-", err.Error())
				} else {
					fmt.Println(filepath.Base(f), "- Updated")
					err := PythonInstall(path + "/" + f)
					if err != nil {
						fmt.Printf("failed to run setup.py, %s", err.Error())
						os.Exit(1)
					}
				}
			}
		}
		return nil
	}
	return command
}
