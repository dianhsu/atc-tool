package cmd

import (
	"fmt"
	"github.com/dianhsu/atc/client"
	"github.com/dianhsu/atc/config"
	"github.com/docopt/docopt-go"
	"os"
	"path/filepath"
)

type ParsedArgs struct {
	Info      client.Info
	File      string
	Specifier []string `docopt:"<specifier>"`
	Alias     string   `docopt:"<alias>"`
	Accepted  bool     `docopt:"ac"`
	All       bool     `docopt:"all"`
	Username  string   `docopt:"<username>"`
	Version   string   `docopt:"{version}"`
	Config    bool     `docopt:"config"`
	Submit    bool     `docopt:"submit"`
	List      bool     `docopt:"list"`
	Parse     bool     `docopt:"parse"`
	Gen       bool     `docopt:"gen"`
	Test      bool     `docopt:"test"`
	Watch     bool     `docopt:"watch"`
	Open      bool     `docopt:"open"`
	Stand     bool     `docopt:"stand"`
	Sid       bool     `docopt:"sid"`
	Race      bool     `docopt:"race"`
	Pull      bool     `docopt:"pull"`
}

var Args *ParsedArgs

func parseArgs(opts docopt.Opts) error {
	cfg := config.Instance
	cln := client.Instance
	path, err := os.Getwd()
	if err != nil {
		return err
	}
	if file, ok := opts["--file"].(string); ok {
		Args.File = file
	} else if file, ok := opts["<file>"].(string); ok {
		Args.File = file
	}
	if Args.Username == "" {
		Args.Username = cln.Username
	}
	info := client.Info{}
	for _, arg := range Args.Specifier {
		parsed := parseArg(arg)
		if value, ok := parsed["contestID"]; ok {
			if info.ContestID != "" && info.ContestID != value {
				return fmt.Errorf("contest ID conflicts: %v %v", info.ContestID, value)
			}
			info.ContestID = value
		}
		if value, ok := parsed["problemID"]; ok {
			if info.ProblemID != "" && info.ProblemID != value {
				return fmt.Errorf("problem ID conflicts: %v %v", info.ProblemID, value)
			}
			info.ProblemID = value
		}
		if value, ok := parsed["submissionID"]; ok {
			if info.SubmissionID != "" && info.SubmissionID != value {
				return fmt.Errorf("submission ID conflicts: %v %v", info.SubmissionID, value)
			}
			info.SubmissionID = value
		}
	}
	root := cfg.FolderName["root"]
	info.RootPath = filepath.Join(path, root)
	for {
		base := filepath.Base(path)
		if base == root {
			info.RootPath = path
			break
		}
		if filepath.Dir(path) == path {
			break
		}
		base = filepath.Dir(path)
	}
	Args.Info = info
	return nil
}

func parseArg(arg string) map[string]string {
	output := make(map[string]string)
	return output
}
