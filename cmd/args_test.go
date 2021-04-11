package cmd

import (
	"flag"
	"fmt"
	"os"
	"testing"
)

func TestArgs(t *testing.T) {
	{
		os.Args = []string{"pyx", "--version"}
		p := &Args{}
		fs := flag.NewFlagSet("test", flag.ContinueOnError)

		args, err := p.Parse(fs)
		if err != nil {
			t.Error(err)
		}

		AssertEqual(t, args.version, true)
		AssertEqual(t, args.help, false)
	}

	{
		os.Args = []string{"pyx", "--help"}
		p := &Args{}
		fs := flag.NewFlagSet("test", flag.ContinueOnError)

		args, err := p.Parse(fs)
		if err != nil {
			t.Error(err)
		}

		AssertEqual(t, args.help, true)
		AssertEqual(t, args.version, false)
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

		AssertEqual(t, args.help, false)
		AssertEqual(t, args.version, false)
		AssertEqual(t, args.branch, "master")
		AssertEqual(t, args.repo, "darumatic/pyx")
		AssertEqual(t, args.script, "scripts/hello.py")
		AssertEqual(t, len(args.scriptArgs), 1)
	}
}
