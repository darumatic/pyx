package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/go-git/go-git/v5"
	gitssh "github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
)

func Install(repoURL string) {
	fmt.Printf("clone %s\n", repoURL)
	project := path.Base(repoURL)
	dir := AppHomeDir() + "/" + project

	if strings.HasPrefix(repoURL, "git@") {
		s := fmt.Sprintf("%s/.ssh/id_rsa", UserHomeDir())
		sshKey, err := ioutil.ReadFile(s)
		if err != nil {
			fmt.Println("missing ssh private key")
			os.Exit(1)
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

		_, cloneError := git.PlainClone(dir, false, &git.CloneOptions{
			Auth:         auth,
			URL:          repoURL,
			Progress:     os.Stdout,
			SingleBranch: true,
		})

		if cloneError != nil {
			fmt.Println("failed to clone repository", cloneError)
			os.Exit(1)
		}
	} else {
		_, cloneError := git.PlainClone(dir, false, &git.CloneOptions{
			URL:          repoURL,
			Progress:     os.Stdout,
			SingleBranch: true,
		})

		if cloneError != nil {
			fmt.Println("failed to clone repository", cloneError)
			os.Exit(1)
		}
	}

	err := PythonInstall(dir)
	if err != nil {
		fmt.Printf("failed to run setup.py, %s", err.Error())
		os.Exit(1)
	}
}

func PythonInstall(dir string) error {
	if FileExists(dir + "/setup.py") {
		EnsurePythonInstalled()

		python, _ := GetPython()
		comd := exec.Command(python, "setup.py", "install")
		comd.Dir = dir
		comd.Stderr = os.Stderr
		comd.Stdout = os.Stdout
		comd.Stdin = os.Stdin
		if err := comd.Run(); err != nil {
			return err
		}
		return nil
	}
	return nil
}

func MakeInstall() *cobra.Command {
	var command = &cobra.Command{
		Use:          "install",
		Short:        "install script git repository",
		Long:         `install script git repository`,
		Example:      `  dev install https://github.com/x1/x2`,
		SilenceUsage: false,
	}

	command.Run = func(command *cobra.Command, args []string) {
		if len(args) == 0 {
			usage := command.UsageString()
			fmt.Printf(usage)
			os.Exit(0)
		}
		for _, repoURL := range args {
			Install(repoURL)
		}
	}
	return command
}
