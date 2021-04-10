package main

import (
	"os"
	"pyx/cmd"
)

func main() {
	pyx := cmd.Pyx{}
	code := pyx.Run()
	os.Exit(code)
}
