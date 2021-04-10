package cmd

import (
	"flag"
	"os"
)

type Args struct {
	version    bool
	help       bool
	repo       string
	branch     string
	script     string
	scriptArgs []string
}

func (args *Args) Parse(fs *flag.FlagSet) (*Args, error) {
	fs.BoolVar(&args.version, "version", false, "pyx version")
	fs.BoolVar(&args.help, "help", false, "pyx help")
	fs.StringVar(&args.branch, "branch", "", "Git branch")

	err := fs.Parse(os.Args[1:])
	if err != nil {
		return nil, err
	}

	var scriptArgs []string

	for i := 0; i < fs.NArg(); i++ {
		if i == 0 {
			args.repo = fs.Arg(i)
		} else if i == 1 {
			args.script = fs.Arg(i)
		} else {
			scriptArgs = append(scriptArgs, fs.Arg(i))
		}
	}
	args.scriptArgs = scriptArgs
	return args, nil
}
