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
	bases    []string
	dryRun   bool
}

const (
	defaultBases = "main,master"
)

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
	flag.Func("bases", `base branches to update (default "`+defaultBases+`")`, func(s string) error {
		for _, v := range strings.Split(s, ",") {
			v = strings.TrimSpace(v)
			if v != "" {
				flags.bases = append(flags.bases, v)
			}
		}
		return nil
	})
	flag.BoolVar(&flags.dryRun, "dry-run", false, "dry run")
	printVer := flag.Bool("version", false, "print version")
	flag.Parse()

	if *printVer {
		printVersionAndExit()
	}

	if flags.token == "" {
		printErrorAndExit("GitHub access token is required")
	}

	switch flag.NArg() {
	case 0:
		printErrorAndExit("pass owner/repo and patterns")
	case 1:
		printErrorAndExit("pass one or more patterns")
	}

	args := flag.Args()
	ownerRepo := strings.Split(args[0], "/")

	if len(ownerRepo) != 2 {
		printErrorAndExit("invalid OWNER/REPO format")
	}

	flags.owner = ownerRepo[0]
	flags.repo = ownerRepo[1]
	flags.patterns = args[1:]

	if len(flags.bases) == 0 {
		flags.bases = strings.Split(defaultBases, ",")
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

func printErrorAndExit(format string, a ...any) {
	fmt.Fprintf(os.Stderr, "error: "+format+"\n", a...)
	os.Exit(1)
}
