package main

import (
	"flag"
	"fmt"
	"os"
)

func mainInner() error {
	flag.Parse()

	flag.Usage = func() {
		flag.PrintDefaults()
		fmt.Println()
	}

	if flag.NArg() == 0 {
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
