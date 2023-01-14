package main

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/google/go-github/v49/github"
	"github.com/winebarrel/pru"
)

func init() {
	log.SetFlags(0)
}

func main() {
	flags := parseFlags()
	ctx := context.Background()
	client := pru.NewGitHubClient(ctx, flags.token)
	pulls, err := pru.ListOpenPullRequests(ctx, client, flags.owner, flags.repo)

	if err != nil {
		log.Fatal(err)
	}

	for _, pr := range pulls {
		files, err := pru.ListPullRequestFiles(ctx, client, pr)

		if err != nil {
			log.Fatal(err)
		}

		for _, pat := range flags.patterns {
			ok, err := match(pat, files)

			if err != nil {
				log.Fatal(err)
			}

			if ok {
				if flags.dryRun {
					log.Printf("update %s (dry run)\n", *pr.HTMLURL)
					break
				}

				log.Printf("update %s\n", *pr.HTMLURL)
				err := pru.UpdatePullRequestBranch(ctx, client, pr)

				if err == nil {
					break
				}

				errResp := &github.ErrorResponse{}

				if !errors.As(err, &errResp) {
					log.Fatal(err)
				}

				if strings.Contains(errResp.Error(), "422 merge conflict") {
					log.Printf("warning: %s", errResp)
					break
				}

				log.Fatal(errResp)
			}
		}
	}
}

func match(pattern string, files []string) (bool, error) {
	for _, f := range files {
		ok, err := doublestar.Match(pattern, f)

		if err != nil {
			return false, err
		}

		if ok {
			return true, nil
		}
	}

	return false, nil
}
