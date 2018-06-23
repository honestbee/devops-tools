package github

import (
	"context"
	"errors"

	"github.com/google/go-github/github"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

// ErrNotFound 404 not found
var ErrNotFound = errors.New("404 Not found")

// Client alias for github client
type Client *github.Client

type (
	// Github struct maps the params we need to query Github
	Github struct {
		Organization string `json:"organization"`
		Client       Client
	}

	// Repo sruct type for github repo
	Repo struct {
		RepoName        string
		RepoDescription string
	}
)

// ListRepos will list all repos for the organization
func (gh *Github) ListRepos() ([]Repo, error) {
	ctx := context.Background()
	var repos []Repo
	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	// get all pages of results
	var allRepos []*github.Repository
	for {
		repos, resp, err := gh.Client.Repositories.ListByOrg(ctx, gh.Organization, opt)
		if err != nil {
			return nil, ErrNotFound
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
		log.Debugf("Github Rate Remaining: %v\n", resp.Rate.Remaining)
	}

	for _, r := range allRepos {
		repoName := *r.FullName
		repoDescription := ""
		if r.Description != nil {
			repoDescription = *r.Description
		}
		repos = append(repos, Repo{
			repoName,
			repoDescription,
		})

	}
	return repos, nil
}

// GetGithubTokenClient create a new github token
func GetGithubTokenClient(ctx context.Context, token string) Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return Client(github.NewClient(tc))
}
