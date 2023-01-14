package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	version string
)

type Flags struct {
	token    string
	owner    string
	repo     string
	patterns []string
}

func init() {
	cmdLine := flag.NewFlagSet(filepath.Base(os.Args[0]), flag.ExitOnError)

	cmdLine.Usage = func() {
		fmt.Fprintf(cmdLine.Output(), "Usage: %s [OPTION] OWNER/REPO PATTERNS...\n", cmdLine.Name())
		cmdLine.PrintDefaults()
	}

	flag.CommandLine = cmdLine
}

func parseFlags() *Flags {
	flags := &Flags{}
	flag.StringVar(&flags.token, "token", os.Getenv("GITHUB_TOKEN"), "GitHub access token. use $GITHUB_TOKEN env")
	printVer := flag.Bool("version", false, "print version")
	flag.Parse()

	if *printVer {
		printVersionAndExit()
	}

	if flags.token == "" {
		printErrorAndExit("GitHub access token is required")
	}

	if flag.NArg() == 0 {
		printUsageAndExit()
	}

	args := flag.Args()
	ownerRepo := strings.Split(args[0], "/")

	if len(ownerRepo) != 2 {
		printErrorAndExit("invalid OWNER/REPO format")
	}

	flags.owner = ownerRepo[0]
	flags.repo = ownerRepo[1]

	if len(args) >= 2 {
		flags.patterns = args[1:]
	}

	return flags
}

func printVersionAndExit() {
	v := version

	if v == "" {
		v = "<nil>"
	}

	fmt.Fprintln(flag.CommandLine.Output(), v)
	os.Exit(0)
}

func printUsageAndExit() {
	flag.CommandLine.Usage()
	os.Exit(1)
}

func printErrorAndExit(format string, a ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", a...)
	os.Exit(1)
}
