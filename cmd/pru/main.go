package main

import (
	"context"
	"fmt"
	"log"

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

		if len(flags.patterns) == 0 {
			update(ctx, client, pr)
		} else {
			for _, pat := range flags.patterns {
				if ok := match(pat, files); ok {
					update(ctx, client, pr)
				}
			}
		}
	}
}

func match(pattern string, files []string) bool {
	for _, f := range files {
		ok, err := doublestar.Match(pattern, f)

		if err != nil {
			log.Fatal(err)
		}

		if ok {
			return true
		}
	}

	return false
}

func update(ctx context.Context, client *github.Client, pr *github.PullRequest) {
	fmt.Printf("update %s\n", *pr.HTMLURL)
	err := pru.UpdatePullRequestBranch(ctx, client, pr)

	if err != nil {
		log.Fatal(err)
	}
}
