package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"

	"github.com/magalab/seekai-cli/cmd"
)

var version = "dev"

func main() {
	var cli cmd.CLI
	parser, err := kong.New(&cli,
		kong.Name("seekai"),
		kong.Description("Seek DB AI model, endpoint, and AI function CLI."),
		kong.Vars{"version": version},
		kong.UsageOnError(),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	ctx, err := parser.Parse(os.Args[1:])
	if err != nil {
		parser.FatalIfErrorf(err)
	}

	if err := ctx.Run(&cli.Globals); err != nil {
		cmd.PrintError(err)
		os.Exit(1)
	}
}
