package cmd

import (
	"flag"
	"fmt"
	"os"
	"testing"
)

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Fatalf("%s != %s", a, b)
	}
}

func TestArgs(t *testing.T) {
	{
		os.Args = []string{"pyx", "--version"}
		p := &Args{}
		fs := flag.NewFlagSet("test", flag.ContinueOnError)

		args, err := p.Parse(fs)
		if err != nil {
			t.Error(err)
		}

		assertEqual(t, args.version, true)
		assertEqual(t, args.help, false)
	}

	{
		os.Args = []string{"pyx", "--help"}
		p := &Args{}
		fs := flag.NewFlagSet("test", flag.ContinueOnError)

		args, err := p.Parse(fs)
		if err != nil {
			t.Error(err)
		}

		assertEqual(t, args.help, true)
		assertEqual(t, args.version, false)
	}

	{
		os.Args = []string{"pyx", "--branch=master", "darumatic/pyx", "scripts/hello.py", "--help"}
		p := &Args{}
		fs := flag.NewFlagSet("test", flag.ContinueOnError)

		args, err := p.Parse(fs)
		if err != nil {
			t.Error(err)
		}

		fmt.Printf("%+v\n", *args)

		assertEqual(t, args.help, false)
		assertEqual(t, args.version, false)
		assertEqual(t, args.branch, "master")
		assertEqual(t, args.repo, "darumatic/pyx")
		assertEqual(t, args.script, "scripts/hello.py")
		assertEqual(t, len(args.scriptArgs), 1)
	}
}
