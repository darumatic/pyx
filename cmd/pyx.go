package cmd

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/cheggaaa/pb/v3"
	giturls "github.com/whilp/git-urls"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

var version = "1.0.5"

type Pyx struct {
}

func (pyx Pyx) Run() (code int) {
	if len(os.Args) == 1 {
		pyx.help()
		return 0
	}
	p := &Args{}
	fs := flag.NewFlagSet("pyx", flag.ContinueOnError)

	args, err := p.Parse(fs)
	if err != nil {
		Error("Incorrect parameters")
		Usage()
		return 1
	}

	if args.version {
		pyx.version()
		return 0
	}

	if args.help {
		pyx.help()
		return 0
	}
	fmt.Printf("Args %s", args.scriptArgs)

	if args.repo != "" && args.script != "" {
		EnsurePythonInstalled()

		if isGithubScript(args.repo, args.script) {
			repo := fmt.Sprintf("https://github.com/%s.git", args.repo)
			return pyx.runGitScript(repo, args.branch, args.script, args.scriptArgs)
		} else if isGitScript(args.repo, args.script) {
			return pyx.runGitScript(args.repo, args.branch, args.script, args.scriptArgs)
		} else if isLocalScript(args.repo, args.script) {
			return pyx.runLocalScript(args.repo, args.script, args.scriptArgs)
		}
	}

	Error("Incorrect parameters")
	Usage()
	return 1
}

func (pyx Pyx) version() (status int) {
	fmt.Printf("pyx %s\n", version)
	python, err := GetPython()
	if err == nil {
		pythonVersion, _ := GetPythonVersion(python)
		fmt.Printf("python %s\n", pythonVersion)
	}
	return 0
}

func (pyx Pyx) help() (status int) {
	fmt.Printf("Single command to run python3 script anywhere.\n\n")
	Usage()
	return 0
}

func (pyx Pyx) runGitScript(repoURL string, branch string, script string, scriptArgs []string) (code int) {
	dir, err := cloneGitRepo(repoURL, branch)
	if err != nil {
		Error("failed to checkout git repository %s, error=%s\n", repoURL, err.Error())
		os.Exit(1)
	}

	InitPythonProject(dir)
	scriptFile := path.Join(dir, script)
	commandArgs := []string{scriptFile}
	if len(scriptArgs) > 0 {
		commandArgs = append(commandArgs, scriptArgs...)
	}
	RunPython(commandArgs...)
	return 0
}

func (pyx Pyx) runLocalScript(dir string, script string, scriptArgs []string) (code int) {
	scriptFile := path.Join(dir, script)
	if !FileExists(scriptFile) {
		Error("script doesn't exist, path=%s", scriptFile)
		return 1
	}

	InitPythonProject(dir)
	commandArgs := []string{scriptFile}
	if len(scriptArgs) > 0 {
		commandArgs = append(commandArgs, scriptArgs...)
	}
	RunPython(commandArgs...)
	return 0
}

func Usage() {
	fmt.Print("Usage :\npyx [args] REPO SCRIPT\n\n")
	fmt.Println("REPO : Git repository URL or local folder")
	fmt.Println("SCRIPT : python script relative path")
	fmt.Println("Arguments :")
	fmt.Println("--branch : Git repository branch")
	fmt.Println("--version : display pyx version")
	fmt.Println("--help : display help and exit")

	fmt.Println("\nExamples :")
	fmt.Println("1) Run remote git repository python script")
	fmt.Println("$ pyx --branch=master https://github.com/darumatic/pyx.git scripts/hello.py")
	fmt.Println("or")
	fmt.Print("$ pyx git@github.com:darumatic/pyx.git scripts/hello.py\n\n")
	fmt.Println("2) Support short name for github repositories")
	fmt.Print("$ pyx darumatic/pyx scripts/hello.py\n\n")
	fmt.Println("3) Run local python script")
	fmt.Println("$ pyx /opt/darumatic/pyx scripts/hello.py")
}

func Error(message string, args ...interface{}) {
	var errorMessage string
	if len(args) > 0 {
		errorMessage = fmt.Sprintf(message, args...)
	} else {
		errorMessage = message
	}
	fmt.Printf("Error: %s\n", errorMessage)
}

func isGithubScript(project string, script string) bool {
	r, _ := regexp.Compile("^[^/:]+/[^/:]+$")
	if strings.Index(project, "..") != 0 && strings.Index(project, ".") != 0 && r.MatchString(project) && isPythonFile(script) {
		return true
	}
	return false
}

func isGitScript(project string, script string) bool {
	url, err := giturls.Parse(project)
	return err == nil && url.Scheme != "file" && isPythonFile(script)
}

func isLocalScript(project string, script string) bool {
	scriptFile := path.Join(project, script)
	return isPythonFile(scriptFile)
}

func isPythonFile(name string) bool {
	r, _ := regexp.Compile(".*\\.py")
	return r.MatchString(name)
}

func cloneGitRepo(repo string, branch string) (string, error) {
	targetDir := path.Join(RepositoryHome(), normalizeRepoName(repo))
	if DirExists(targetDir) {
		err := GitPull(targetDir, branch)
		return targetDir, err
	} else {
		err := GitClone(repo, branch, targetDir)
		return targetDir, err
	}
}

func normalizeRepoName(repoURL string) string {
	u, _ := giturls.Parse(repoURL)
	return path.Join(strings.ReplaceAll(u.Hostname(), ".", "_"), strings.ReplaceAll(u.Path[1:], "/", "_"))
}

func RepositoryHome() string {
	return filepath.Join(UserHomeDir(), ".pyx", "repo")
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
