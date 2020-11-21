package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strings"
)

var (
	Version = "1.0.1"
)

func MakeDev() *cobra.Command {
	var command = &cobra.Command{
		Use:                "dev",
		Short:              "python script manager",
		Long:               `python script manager`,
		TraverseChildren:   true,
		DisableFlagParsing: true,
		RunE: func(cmdc *cobra.Command, args []string) error {
			if len(args) == 0 {
				Help(cmdc)
				os.Exit(1)
			}

			for _, folder := range Repository() {
				scriptName := args[0]
				if !strings.HasSuffix(scriptName, ".py") {
					scriptName = "cmd/" + scriptName + ".py"
				}
				script := AppHomeDir() + "/" + folder + "/" + scriptName
				if FileExists(script) {
					commandArgs := append([]string{script}, args[1:]...)
					RunPython(commandArgs...)
				}
			}

			fmt.Println("No script found")
			os.Exit(1)
			return nil
		},
	}

	return command
}

func Repository() []string {
	return FindLastFolder(UserHomeDir() + "/.dev/")
}

func HttpDownload(url string, pattern string) (string, error) {
	fmt.Printf("Downloading %s\n", url)
	file, err := ioutil.TempFile("", pattern)
	if err != nil {
		return "", err
	}
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download %s, status=%s", url, resp.Status)
	}
	_, err = io.Copy(file, resp.Body)
	file.Close()
	if err != nil {
		return "", err
	}
	return file.Name(), nil
}

func RunCommand(command string, args ...string) (int, error) {
	comd := exec.Command(command, args...)
	comd.Stderr = os.Stderr
	comd.Stdout = os.Stdout
	comd.Stdin = os.Stdin
	if err := comd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return exitError.ExitCode(), err
		} else {
			return 1, err
		}
	} else {
		return 0, nil
	}
}

func CommandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func FindLastFolder(dir string) []string {
	var folders []string
	files, _ := ioutil.ReadDir(dir)

	for _, file := range files {
		if file.Mode().IsDir() {
			folders = append(folders, file.Name())
		}
	}
	sort.Sort(sort.Reverse(sort.StringSlice(folders)))
	return folders
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func DirExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func UserHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	} else if runtime.GOOS == "linux" {
		home := os.Getenv("XDG_CONFIG_HOME")
		if home != "" {
			return home
		}
	}
	return os.Getenv("HOME")
}

func AppHomeDir() string {
	return UserHomeDir() + "/.dev"
}

func ParsePrivateKey() (ssh.Signer, error) {
	s := fmt.Sprintf("%s/.ssh/id_rsa", UserHomeDir())
	if FileExists(s) {
		sshKey, err := ioutil.ReadFile(s)
		if err != nil {
			return nil, err
		}
		return ssh.ParsePrivateKey(sshKey)
	}
	return nil, errors.New("No private key found")
}

// ? Find string in slice
func Find(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
