package cmd

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

const PYTHON_VERSION = "3.9.2"

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
		_, err := InstallPython()
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("couldn't install python3, please manually install python3")
			os.Exit(1)
		}
		python, err = GetPython()
		if err != nil {
			fmt.Println("couldn't install python3, please manually install python3")
			os.Exit(1)
		}
	}

	version, err := GetPythonVersion(python)
	if err != nil {
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
		pythonPath := filepath.Join(PythonHome(), "python.exe")
		if FileExists(pythonPath) {
			return pythonPath, nil
		}
	} else {
		pythonPath := filepath.Join(PythonHome(), "bin", "python3")
		if FileExists(pythonPath) {
			return pythonPath, nil
		}
	}
	pythonPath, err := exec.LookPath("python3")
	if err == nil {
		return pythonPath, nil
	}
	return "", errors.New("python3 not installed")
}

func InstallPython() (bool, error) {
	version := PYTHON_VERSION
	fmt.Printf("Installing python-%s\n", version)
	url, err := pythonBuildURL()
	if err != nil {
		return false, err
	}
	file := path.Join(PythonHome(), "python.zst")
	err = HttpDownload(url, file)
	if err != nil {
		return false, err
	}
	f1, err := os.Open(file)
	if err != nil {
		return false, nil
	}
	defer f1.Close()
	reader := bufio.NewReader(f1)

	tarFile, err := ioutil.TempFile("", "python."+version+".*.tar")

	writer := bufio.NewWriter(tarFile)
	err = DecompressZstd(reader, tarFile)
	if err != nil {
		return false, err
	}
	writer.Flush()
	tarFile.Close()

	newTarFile, err := os.Open(tarFile.Name())
	tarReader := bufio.NewReader(newTarFile)
	err = DecompressTar(tarReader, PYXHome())
	if err != nil {
		return false, err
	}
	fmt.Println("")
	return true, nil
}

func InitPythonProject(dir string) {
	requirementsTxt := path.Join(dir, "requirements.txt")
	if FileExists(requirementsTxt) {
		RunPython("-m", "pip", "install", "-r", "requirements.txt")
	}
}

func pythonBuildURL() (string, error) {
	if runtime.GOOS != "windows" && runtime.GOOS != "darwin" && runtime.GOOS != "linux" {
		return "", errors.New(runtime.GOOS + " not supported")
	}
	var arch string
	if runtime.GOARCH == "386" {
		arch = "i686"
	} else {
		arch = "x86_64"
	}
	return fmt.Sprintf("https://github.com/darumatic/pyx/releases/download/python-%s/python-%s-%s-%s.tar.zst", PYTHON_VERSION, PYTHON_VERSION, arch, runtime.GOOS), nil
}
