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
	// Config struct maps the params we need to query Config
	Config struct {
		Organization string `json:"organization"`
		Client       Client
		Repo         RepoConfig
	}

	// RepoConfig struct define repo attributes
	RepoConfig struct {
		Type string `json:"private"`
	}

	// Repo sruct type for github repo
	Repo struct {
		RepoName        string
		RepoDescription string
	}
)

// ListRepos will list all repos for the organization
func (c *Config) ListRepos() ([]Repo, error) {
	ctx := context.Background()
	var repos []Repo
	opt := &github.RepositoryListByOrgOptions{
		// https://github.com/google/go-github/blob/master/github/repos.go#L195
		Type:        c.Repo.Type,
		ListOptions: github.ListOptions{PerPage: 100},
	}

	// get all pages of results
	var allRepos []*github.Repository
	for {
		repos, resp, err := c.Client.Repositories.ListByOrg(ctx, c.Organization, opt)
		if err != nil {
			return nil, ErrNotFound
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
		log.Debugf("Config Rate Remaining: %v\n", resp.Rate.Remaining)
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

// NewClient create a new github client
func NewClient(ctx context.Context, token string) Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return Client(github.NewClient(tc))
}
