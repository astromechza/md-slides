package main

import (
	"flag"
	"fmt"
	"os"
)

const mainUsage = `%s is a html slide server based on slides defined in a markdown formatted file.

Usage:
  md-slides [subcommand] [options...]

Subcommands:
  check     validate that the slides are parsable
  serve     serve the slides as html
  html      render slide html to a directory
  pdf       export the slides to a pdf

`


var version = "unknown"
var buildDate = "unknown"

func mainInner() error {
	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	versionFlag := fs.Bool("version", false, "Print version information")
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, mainUsage, os.Args[0])
		fs.PrintDefaults()
	}
	if err := fs.Parse(os.Args[1:]); err != nil {
		return err
	}

	if *versionFlag {
		fmt.Printf("Version:     %s\n", version)
		fmt.Printf("Build Date:  %s\n", buildDate)
		fmt.Printf("URL:         https://github.com/astromechza/md-slides\n")
		return nil
	}

	if fs.NArg() == 0 {
		fs.Usage()
		fmt.Fprintf(os.Stderr, "\n")
		return fmt.Errorf("expected subcommand as first argument")
	}

	subcommand := fs.Arg(0)
	switch subcommand {
	case "check":
		return SubcommandCheck(fs.Args()[1:])
	case "serve":
		return SubcommandServe(fs.Args()[1:])
	case "html":
		return SubcommandHTML(fs.Args()[1:])
	case "pdf":
		return SubcommandPDF(fs.Args()[1:])
	default:
		return fmt.Errorf("unknown subcommand '%s' (choose from check, server, pdf)", subcommand)
	}
}

func main() {
	if err := mainInner(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
