package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

const PYTHON_VERSION = "3.8.0"

func Python(cmd *cobra.Command) {
	usage := cmd.UsageString()
	fmt.Printf(usage)
}

func InstallPython() (bool, error) {
	version := PYTHON_VERSION
	if runtime.GOOS == "windows" {
		return windowsInstallPython(version)
	} else if runtime.GOOS == "linux" {
		return linuxInstallPython(version)
	} else if runtime.GOOS == "darwin" {
		return darwinInstallPython(version)
	} else {
		return false, errors.New("platform not support, please manually install python")
	}
}

func MakePython() *cobra.Command {
	var command = &cobra.Command{
		Use:     "python",
		Short:   "manage python installation",
		Long:    `manage python installation`,
		Example: `  dev python install`,
	}
	var installCommand = &cobra.Command{
		Use:     "install",
		Short:   "install python",
		Long:    `install python`,
		Example: `  dev python install`,
	}
	installCommand.RunE = func(command *cobra.Command, args []string) error {
		_, err := InstallPython()
		if err != nil {
			fmt.Printf("failed to install python, %s", err.Error())
			os.Exit(1)
		}
		return nil
	}
	command.AddCommand(installCommand)
	command.RunE = func(command *cobra.Command, args []string) error {
		Python(command)
		return nil
	}
	return command
}

func GetPythonVersion(python string) (string, error) {
	cmd := exec.Command(python, "--version")
	output := new(bytes.Buffer)
	cmd.Stdout = output
	cmd.Stderr = output
	if err := cmd.Run(); err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	version := strings.ReplaceAll(output.String(), "Python", "")
	version = strings.TrimSpace(version)
	return version, nil
}

func EnsurePythonInstalled() {
	python, err := GetPython()
	if err != nil {
		if runtime.GOOS == "windows" {
			fmt.Println("python not found, would you like to install python? (y/n)")
			var input string
			fmt.Scanln(&input)
			if input == "y" || input == "Y" {
				_, err := InstallPython()
				if err != nil {
					fmt.Printf("failed to install python, %s\n", err.Error())
				}
				python, err = GetPython()
			} else {
				os.Exit(1)
			}
		} else {
			fmt.Printf("Error: python not found, please install python3\n")
			os.Exit(1)
		}
	}

	version, err := GetPythonVersion(python)
	if err != nil {
		fmt.Printf(err.Error())
		fmt.Printf("Warning: failed to get python version, dev requires python3\n")
		return
	}

	majorVersion, err := strconv.Atoi(strings.Split(version, ".")[0])
	if majorVersion < 3 {
		fmt.Printf("Warning: python version is %s, dev requires python3\n", version)
		return
	}
}

func GetPython() (string, error) {
	if runtime.GOOS == "windows" {
		path := fmt.Sprintf("%s\\Programs\\Python\\Python38-32\\python.exe", os.Getenv("LocalAppData"))
		if FileExists(path) {
			return path, nil
		}
	}

	path, err := exec.LookPath("python3")
	if err == nil {
		return path, nil
	}

	path, err = exec.LookPath("python")
	if err == nil {
		return path, nil
	}

	return "", err
}

func windowsInstallPython(version string) (bool, error) {
	fileURL := fmt.Sprintf("https://www.python.org/ftp/python/%s/python-%s.exe", version, version)
	file, err := HttpDownload(fileURL, "python."+version+".*.exe")
	if err != nil {
		return false, errors.New(fmt.Sprintf("Failed to download %s, please check your netowrk", fileURL))
	}

	fmt.Printf("Installing python-%s\n", version)

	fmt.Printf("file=%s", file)
	status, err := Exec(file, "/quiet")
	if status != 0 && status != 1602 && status != 1638 {
		return false, errors.New(fmt.Sprintf("Failed to install %s", file))
	}
	fmt.Println("python installed")
	return true, nil
}

func linuxInstallPython(version string) (bool, error) {
	return false, nil
}

func darwinInstallPython(version string) (bool, error) {
	return false, nil
}
