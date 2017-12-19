package main

import (
	"context"

	"github.com/google/go-github/github"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

type (
	// GitHub struct maps the params we need to query GitHub
	GitHub struct {
		Organization string `json:"organization"`
		Token        string `json:"token"`
	}
)

// ListTeams will list all teams for the organization
func (gh *GitHub) ListTeams() ([]*github.Team, error) {
	ctx := context.Background()
	client := getGitHubTokenClient(ctx, gh.Token)
	opt := &github.ListOptions{PerPage: 100}

	var allTeams []*github.Team
	for {
		teams, resp, err := client.Organizations.ListTeams(ctx, gh.Organization, opt)
		if err != nil {
			return nil, err
		}
		log.Debugf("GitHub Response Code: %v\n", resp.StatusCode)
		log.Debugf("GitHub Rate Remaining: %v\n", resp.Rate.Remaining)
		allTeams = append(allTeams, teams...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return allTeams, nil
}

func getGitHubTokenClient(ctx context.Context, token string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}
