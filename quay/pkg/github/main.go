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

type (
	// Github struct maps the params we need to query Github
	Github struct {
		Organization string `json:"organization"`
		Client       *github.Client
	}

	// Repo sruct type for github repo
	Repo struct {
		RepoName        string
		RepoDescription *string
	}

	// TeamRoles struct keeps a reference to the Github team and a map of UserRoles
	TeamRoles struct {
		Team        github.Team
		UserRoles   map[string]string
		SortedUsers []string
	}
	// TeamRolesList struct keeps array of all TeamRoles
	TeamRolesList struct {
		TeamRoles []TeamRoles
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
		repos = append(repos, Repo{
			*r.FullName,
			r.Description,
		})

	}
	return repos, nil
}

// GetGithubTokenClient create a new github token
func GetGithubTokenClient(ctx context.Context, token string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}
