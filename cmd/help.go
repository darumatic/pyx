package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func Help(cmd *cobra.Command) {
	fmt.Printf("Single command to run python3 script anywhere.\n\n")
	python, _ := GetPython()
	fmt.Printf("python: %s\n", python)
	ExampleUsage()
}

func ExampleUsage() {
	fmt.Println("Example usage:")
	fmt.Println("  1) Run git repository scripts")
	fmt.Println("     $ pyx https://github.com/darumatic/pyx scripts/hello.py")
	fmt.Println("     or")
	fmt.Println("     $ pyx git@github.com:darumatic/pyx.git scripts/hello.py")

	fmt.Println("     For github repositories, we could also simply use the repository name.")
	fmt.Println("     $ pyx darumatic/pyx scripts/hello.py")

	fmt.Println("  2) Run http script")
	fmt.Println("     $ pyx https://raw.githubusercontent.com/darumatic/pyx/master/scripts/hello.py")

	fmt.Println("  3) Run local script")
	fmt.Println("     $ pyx hello.py")
}
