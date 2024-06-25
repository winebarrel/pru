package pru

import (
	"context"
	"slices"

	"github.com/google/go-github/v49/github"
	"golang.org/x/oauth2"
)

func NewGitHubClient(ctx context.Context, token string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

func ListOpenPullRequests(ctx context.Context, client *github.Client, owner string, repo string, ignoreLabels []string) ([]*github.PullRequest, error) {
	opt := &github.PullRequestListOptions{
		State: "open",
	}

	var allPulls []*github.PullRequest

	for {
		pulls, resp, err := client.PullRequests.List(ctx, owner, repo, opt)

		if err != nil {
			return nil, err
		}

	PULLS:
		for _, p := range pulls {
			for _, label := range p.Labels {
				slices.Contains(ignoreLabels, *label.Name)
				continue PULLS
			}

			allPulls = append(allPulls, p)
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	return allPulls, nil
}

func ListPullRequestFiles(ctx context.Context, client *github.Client, pull *github.PullRequest) ([]string, error) {
	repo := pull.Head.Repo
	opt := &github.ListOptions{}
	var allFiles []string

	for {
		files, resp, err := client.PullRequests.ListFiles(ctx, *repo.Owner.Login, *repo.Name, *pull.Number, opt)

		if err != nil {
			return nil, err
		}

		for _, f := range files {
			allFiles = append(allFiles, *f.Filename)
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	return allFiles, nil
}

func UpdatePullRequestBranch(ctx context.Context, client *github.Client, pull *github.PullRequest) error {
	repo := pull.Head.Repo

	opt := &github.PullRequestBranchUpdateOptions{
		ExpectedHeadSHA: pull.Head.SHA,
	}

	_, _, err := client.PullRequests.UpdateBranch(ctx, *repo.Owner.Login, *repo.Name, *pull.Number, opt)

	if err != nil {
		if _, ok := err.(*github.AcceptedError); !ok {
			return err
		}
	}

	return nil
}
