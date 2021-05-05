package cmd

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	gitssh "github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"os"
	"strings"
)

func GitClone(repoURL string, branch string, dir string) error {
	if strings.HasPrefix(repoURL, "git@") {
		s := fmt.Sprintf("%s/.ssh/id_rsa", UserHomeDir())
		sshKey, err := ioutil.ReadFile(s)
		if err != nil {
			fmt.Println("missing ssh private key")
			return err
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
			Auth:          auth,
			URL:           repoURL,
			ReferenceName: referenceName(branch),
			Progress:      os.Stdout,
			SingleBranch:  true,
		})
		fmt.Println("")
		return cloneError
	} else {
		_, cloneError := git.PlainClone(dir, false, &git.CloneOptions{
			URL:           repoURL,
			ReferenceName: referenceName(branch),
			Progress:      os.Stdout,
			SingleBranch:  true,
		})
		fmt.Println("")
		return cloneError
	}
}

func GitPull(repo string, branch string) error {
	r, err := git.PlainOpen(repo)
	if err != nil {
		return err
	}
	w, err := r.Worktree()
	if err != nil {
		return err
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
		if err != nil && !strings.Contains(err.Error(), "up-to-date") {
			return err
		}
		err = w.Checkout(&git.CheckoutOptions{
			Create: false,
			Force:  false,
			Branch: referenceName(branch),
		})
		return err
	} else {
		err = w.Pull(&git.PullOptions{RemoteName: "origin"})
		if err != nil && !strings.Contains(err.Error(), "up-to-date") {
			return err
		}
		err = w.Checkout(&git.CheckoutOptions{
			Create: false,
			Force:  false,
			Branch: referenceName(branch),
		})
		return err
	}
}

func referenceName(branch string) (name plumbing.ReferenceName) {
	if branch == "" {
		return ""
	}
	return plumbing.NewBranchReferenceName(branch)
}
