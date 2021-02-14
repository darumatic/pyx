package cmd

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"

	"github.com/cheggaaa/pb/v3"
	"github.com/spf13/cobra"
	giturls "github.com/whilp/git-urls"
)

var (
	Version = "1.0.2"
)

func MakePyx() *cobra.Command {
	var command = &cobra.Command{
		Use:                "pyx",
		Short:              "python script runner",
		Long:               `python script runner`,
		TraverseChildren:   true,
		DisableFlagParsing: false,
		RunE: func(cmdc *cobra.Command, args []string) error {
			if len(args) == 0 {
				Help(cmdc)
				os.Exit(0)
			}
			if isGithubScript(args) {
				repoURL := fmt.Sprintf("https://github.com/%s", args[0])
				branch := cmdc.Flag("branch").Value.String()

				dir, err := cloneGitRepo(repoURL, branch)
				if err != nil {
					fmt.Printf("Error: failed to clone git repository %s, error=%s\n", repoURL, err.Error())
					os.Exit(1)
				}
				script := path.Join(dir, args[1])
				commandArgs := []string{script}
				if len(args) > 2 {
					commandArgs = append(commandArgs, args[2:]...)
				}
				RunPython(commandArgs...)
			} else if isGitScript(args) {
				repoURL := args[0]
				branch := cmdc.Flag("branch").Value.String()
				dir, err := cloneGitRepo(repoURL, branch)
				if err != nil {
					fmt.Printf("Error: failed to checkout git repository %s, error=%s\n", repoURL, err.Error())
					os.Exit(1)
				}
				script := path.Join(dir, args[1])
				commandArgs := []string{script}
				if len(args) > 2 {
					commandArgs = append(commandArgs, args[2:]...)
				}
				RunPython(commandArgs...)
			} else if isHTTPScript(args) {
				script, err := downloadURL(args[0])
				if err != nil {
					fmt.Printf("failed to download script %s, error=%s\n", args[0], err.Error())
					os.Exit(1)
				}
				commandArgs := []string{script}
				fmt.Println(len(args))
				if len(args) > 1 {
					commandArgs = append(commandArgs, args[1:]...)
				}
				RunPython(commandArgs...)
			} else if isLocalScript(args) {
				script := args[0]
				commandArgs := []string{script}
				if len(args) > 1 {
					commandArgs = append(commandArgs, args[1:]...)
				}
				RunPython(commandArgs...)
			} else {
				fmt.Printf("Error: unknown script\n\n")

				ExampleUsage()
				os.Exit(1)
			}
			return nil
		},
	}

	command.Flags().StringP("branch", "b", "master", "Git branch")
	return command
}

func isGithubScript(args []string) bool {
	if len(args) > 1 {
		r, _ := regexp.Compile("^[^/:]+/[^/:]+$")
		if r.MatchString(args[0]) && isPythonFile(args[1]) {
			return true
		}
	}
	return false
}

func isGitScript(args []string) bool {
	if len(args) > 1 {
		_, err := giturls.Parse(args[0])
		return err == nil && isPythonFile(args[1])
	}
	return false
}

func isHTTPScript(args []string) bool {
	if len(args) > 0 {
		return isURL(args[0]) && strings.HasSuffix(strings.ToLower(args[0]), ".py")
	}
	return false
}

func isLocalScript(args []string) bool {
	if len(args) > 0 {
		return isPythonFile(args[0])
	}
	return false
}

func isPythonFile(name string) bool {
	r, _ := regexp.Compile(".*\\.py")
	return r.MatchString(name)
}

func isURL(name string) bool {
	return strings.HasPrefix(strings.ToLower(name), "http:") || strings.HasPrefix(strings.ToLower(name), "https:")
}

func cloneGitRepo(repo string, branch string) (string, error) {
	targetDir := path.Join(RepositoryHome(), normalizeRepoName(repo))
	if DirExists(targetDir) {
		err := GitUpdate(targetDir, branch)
		return targetDir, err
	} else {
		err := GitClone(repo, branch, targetDir)
		return targetDir, err
	}
}

func downloadURL(scriptURL string) (string, error) {
	targetFile := path.Join(RepositoryHome(), "http", normalizeURLName(scriptURL))
	err := HttpDownload(scriptURL, targetFile)
	return targetFile, err
}

func normalizeRepoName(repoURL string) string {
	u, _ := giturls.Parse(repoURL)
	return path.Join(strings.ReplaceAll(u.Hostname(), ".", "_"), strings.ReplaceAll(u.Path, "/", "_"))
}

func normalizeURLName(scriptURL string) string {
	u, _ := url.Parse(scriptURL)
	return path.Join(strings.ReplaceAll(u.Hostname(), ".", "_"), strings.ReplaceAll(u.Path, "/", "_"))
}

func RepositoryHome() string {
	return filepath.Join(UserHomeDir(), ".pyx", "cache")
}

func PythonHome() string {
	return filepath.Join(UserHomeDir(), ".pyx", "python")
}

func HttpDownload(url string, output string) error {
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	os.MkdirAll(path.Dir(output), os.ModePerm)
	file, _ := os.OpenFile(output, os.O_CREATE|os.O_WRONLY, 0644)
	bar := pb.Full.Start64(resp.ContentLength)
	writer := bufio.NewWriter(file)
	barWriter := bar.NewProxyWriter(writer)
	_, _ = io.Copy(barWriter, resp.Body)
	barWriter.Close()
	bar.Finish()
	_ = file.Close()
	return nil
}

func RunCommand(command string, args ...string) (int, error) {
	cmd := exec.Command(command, args...)

	cmd.Env = os.Environ()
	executable, _ := os.Executable()
	cmd.Env = append(cmd.Env, fmt.Sprintf("PYX_CLI=%s", executable))
	cmd.Env = append(cmd.Env, fmt.Sprintf("PYX_HOME=%s", PYXHome()))

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
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

func FindFolders(dir string) []string {
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

func PYXHome() string {
	return filepath.Join(UserHomeDir(), ".pyx")
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
