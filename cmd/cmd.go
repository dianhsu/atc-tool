package cmd

import "github.com/docopt/docopt.go"

func Eval(opts docopt.Opts) error {
	Args = &ParsedArgs{}
	err := opts.Bind(Args)
	if err != nil {
		return err
	}
	if err := parseArgs(opts); err != nil {
		return err
	}
	return nil
}
