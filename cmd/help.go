package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"strings"
)

func Help(cmd *cobra.Command) {
	path := AppHomeDir()
	var scripts []string
	var sb strings.Builder
	folders := Repository()
	for _, folder := range folders {
		files, _ := ioutil.ReadDir(path + "/" + folder + "/cmd")
		for _, file := range files {
			if strings.HasSuffix(file.Name(), ".py") {
				if !Find(scripts, file.Name()) {
					scripts = append(scripts, file.Name())
					commandName := strings.ReplaceAll(file.Name(), ".py", "")
					s := fmt.Sprintf("%-10v", commandName)
					sb.WriteString("  " + s)
					sb.WriteString("  " + folder + "/cmd/" + file.Name())
					sb.WriteString("\n")
				}
			}
		}
	}
	fmt.Printf("dev makes delivering python scripts easier.\n")
	python, _ := GetPython()
	version, _ := GetPythonVersion(python)
	fmt.Printf("python %s, version %s.\n\n", python, version)
	usage := cmd.UsageString()
	usage = strings.ReplaceAll(usage, "dev version", "dev version\n"+sb.String())
	fmt.Printf(usage)
}

func MakeHelp() *cobra.Command {
	var command = &cobra.Command{
		Use:     "help",
		Short:   "help",
		Long:    `help document of dev cli`,
		Example: `  dev help`,
	}

	command.RunE = func(command *cobra.Command, args []string) error {
		Help(command.Parent())
		return nil
	}

	return command
}
