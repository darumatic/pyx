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

func InstallPython() (bool, error) {
	version := PYTHON_VERSION
	if runtime.GOOS == "windows" {
		return windowsInstallPython(version)
	} else if runtime.GOOS == "linux" {
		return linuxInstallPython(version)
	} else if runtime.GOOS == "darwin" {
		return darwinInstallPython(version)
	} else {
		return false, errors.New("couldn't install python3, please manually install python3")
	}
}

func MakePython() *cobra.Command {
	var command = &cobra.Command{
		Use:                "python",
		Short:              "run python script",
		Long:               `run python script`,
		Example:            `  dev python script.py`,
		TraverseChildren:   true,
		DisableFlagParsing: true,
	}
	command.RunE = func(command *cobra.Command, args []string) error {
		RunPython(args...)
		return nil
	}
	return command
}

func RunPython(args ...string) {
	python := EnsurePythonInstalled()
	_, err := RunCommand(python, args...)
	if err != nil {
		os.Exit(1)
	}
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

func EnsurePythonInstalled() string {
	python, err := GetPython()
	if err != nil {
		fmt.Print("python3 not found, would you like to install python3? (y/n)")
		var input string
		fmt.Scanln(&input)
		if input == "y" || input == "Y" {
			_, err := InstallPython()
			if err != nil {
				fmt.Println("couldn't install python3, please manually install python3")
			}
			python, err = GetPython()
		} else {
			os.Exit(1)
		}
	}

	version, err := GetPythonVersion(python)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("couldn't install python3, please manually install python3")
		os.Exit(1)
	}

	majorVersion, err := strconv.Atoi(strings.Split(version, ".")[0])
	if majorVersion < 3 {
		fmt.Println("couldn't install python3, please manually install python3")
		os.Exit(1)
	}
	return python
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
	return "", err
}

func windowsInstallPython(version string) (bool, error) {
	fileURL := fmt.Sprintf("https://www.python.org/ftp/python/%s/python-%s.exe", version, version)
	file, err := HttpDownload(fileURL, "python."+version+".*.exe")
	if err != nil {
		return false, errors.New(fmt.Sprintf("Failed to download %s, please check your netowrk", fileURL))
	}
	fmt.Printf("Installing python-%s\n", version)
	status, err := RunCommand(file, "/quiet")
	if status != 0 && status != 1602 && status != 1638 {
		return false, errors.New(fmt.Sprintf("Failed to install %s", file))
	}
	fmt.Println("python installed")
	return true, nil
}

func linuxInstallPython(version string) (bool, error) {
	if CommandExists("apt") {
		_, err := RunCommand("/bin/sh", "-c", "sudo apt install python3")
		if err != nil {
			return false, errors.New("failed to install python3")
		}
	} else if CommandExists("apt") {
		_, err := RunCommand("/bin/sh", "-c", "sudo yum install python3")
		if err != nil {
			return false, errors.New("failed to install python3")
		}
	} else {
		return false, errors.New("please install python3 manually")
	}
	return true, nil
}

func darwinInstallPython(version string) (bool, error) {
	fileURL := fmt.Sprintf("https://www.python.org/ftp/python/%s/python-%s-macosx10.9.pkg", version, version)
	file, err := HttpDownload(fileURL, "python."+version+".*.pkg")
	if err != nil {
		return false, errors.New(fmt.Sprintf("Failed to download %s, please check your netowrk", fileURL))
	}
	fmt.Printf("Installing python-%s\n", version)
	status, err := RunCommand("/bin/sh", "-c", fmt.Sprintf("sudo installer -pkg %s -target /Applications", file))
	if status != 0 {
		return false, errors.New(fmt.Sprintf("Failed to install %s", file))
	}
	fmt.Println("python installed")
	return true, nil
}
