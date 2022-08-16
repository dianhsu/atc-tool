package main

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"os"
	"strings"

	"github.com/dianhsu/atc/client"
	"github.com/dianhsu/atc/cmd"
	"github.com/dianhsu/atc/config"
	docopt "github.com/docopt/docopt.go"
	"github.com/fatih/color"
	ansi "github.com/k0kubun/go-ansi"
)

const version = "v1.0.0"
const configPath = "~/.atc/config.yaml"
const sessionPath = "~/.atc/session.yaml"

func main() {
	usage := `AtCoder Tool $%version%$ (atc). https://github.com/dianhsu/atc-tool

Usage:
	atc [options]
	
Options:
	-h			Show this screen.
	--version	Show version.
`
	color.Output = ansi.NewAnsiStdout()
	usage = strings.Replace(usage, `$%version%$`, version, 1)
	opts, _ := docopt.ParseArgs(usage, os.Args[1:], fmt.Sprintf("AtCoder Tool %v (atc)", version))
	opts[`{version}`] = version

	cfgPath, _ := homedir.Expand(configPath)
	clnPath, _ := homedir.Expand(sessionPath)
	config.Init(cfgPath)
	client.Init(clnPath, config.Instance.Host, config.Instance.Proxy)

	err := cmd.Eval(opts)
	if err != nil {
		color.Red(err.Error())
	}
	color.Unset()
}
