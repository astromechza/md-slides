package main

import (
	"flag"
	"fmt"
	"os"
)

const mainUsage = `md-slides is an html slide server based on slides defined in a markdown formatted file.

Usage:
  md-slides [subcommand] [options...]

Subcommands:
  serve     serve the slides as html
  html      render slide html to a directory

Options:
`

// These are "build-time" vars that will be filled in during the final build
var (
	version = "unknown"
	buildDate = "unknown"
)

func mainInner() error {
	// Define and parse the top level cli flags - each subcommand has their own flag set too!
	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	versionFlag := fs.Bool("version", false, "Show version information")
	fs.Usage = func() {
		_, _ = fmt.Fprint(os.Stderr, mainUsage)
		fs.PrintDefaults()
	}
	if err := fs.Parse(os.Args[1:]); err != nil {
		return err
	}

	// Handle version flag as a command (sneaky but that's kind of convention these days)
	if *versionFlag {
		fmt.Printf("Version:     %s\n", version)
		fmt.Printf("Build Date:  %s\n", buildDate)
		fmt.Printf("URL:         https://github.com/astromechza/md-slides\n")
		return nil
	}

	// First argument is the subcommand name
	if fs.NArg() == 0 {
		fs.Usage()
		_, _ = fmt.Fprintf(os.Stderr, "\n")
		return fmt.Errorf("expected a subcommand as the first argument")
	}
	subcommand := fs.Arg(0)

	// Execute the chosen subcommand or throw an appropriate message
	switch subcommand {
	case "serve":
		return SubcommandServe(fs.Args()[1:])
	case "html":
		return SubcommandHTML(fs.Args()[1:])
	default:
		return fmt.Errorf("unknown subcommand '%s' (choose from check, serve, html, pdf)", subcommand)
	}
}

func main() {
	if err := mainInner(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
