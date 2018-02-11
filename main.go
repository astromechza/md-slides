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
	serve - serve the slides as html
`

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
