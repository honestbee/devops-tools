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

// Query will query GitHub
func (g *GitHub) Query() error {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: g.Token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 30},
	}

	// get all pages of results
	var allRepos []*github.Repository
	for {
		repos, resp, err := client.Repositories.ListByOrg(ctx, g.Organization, opt)
		if err != nil {
			return err
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
		log.Debugf("GitHub Rate Remaining: %v\n", resp.Rate.Remaining)
	}

	for _, r := range allRepos {
		repoName := *r.FullName
		repoDescription := ""
		if r.Description != nil {
			repoDescription = *r.Description
		}
		log.Infof("%v: %v\n",
			repoName,
			repoDescription,
		)
	}
	return nil
}
