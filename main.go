package main

import (
	"flag"
	"fmt"
	"os"
)

const mainUsage = `md-slides is a html slide server based on slides defined in a markdown formatted file.

Usage:
	md-slides [subcommand] [options...]

Subcommands:
	serve     serve the slides as html
	pdf       render slides to pdf
	version   print version information
`

var commitHash = "unknown"
var buildDate = "unknown"
var gitVersion = "unknown"

func mainInner() error {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, mainUsage)
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		fmt.Fprintf(os.Stderr, "\n")
		return fmt.Errorf("expected subcommand as first argument")
	}

	subcommand := flag.Arg(0)
	switch subcommand {
	case "serve":
		return Serve(flag.Args()[1:])
	case "pdf":
		return PDF(flag.Args()[1:])
	case "version":
		fmt.Printf("Version: %s\n", gitVersion)
		fmt.Printf("Hash:    %s\n", commitHash)
		fmt.Printf("Date:    %s\n", buildDate)
		fmt.Printf("URL:     https://github.com/AstromechZA/md-slides\n")
		return nil
	default:
		return fmt.Errorf("unknown subcommand '%s'", subcommand)
	}
}

func main() {
	if err := mainInner(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
