package main

import (
	"context"
	"fmt"

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
	// TeamRoles struct keeps a reference to the GitHub team and a map of UserRoles
	TeamRoles struct {
		Team      github.Team
		UserRoles map[string]string
	}
)

// ListRepos will list all repos for the organization
func (gh *GitHub) ListRepos() error {
	ctx := context.Background()
	client := getGitHubTokenClient(ctx, gh.Token)

	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	// get all pages of results
	var allRepos []*github.Repository
	for {
		repos, resp, err := client.Repositories.ListByOrg(ctx, gh.Organization, opt)
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
		fmt.Printf("%v: %v\n",
			repoName,
			repoDescription,
		)
	}
	return nil
}

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

// GetTeamRoles returns a map of team members (Login) to role
func (gh *GitHub) GetTeamRoles(t *github.Team) (TeamRoles, error) {
	tm := TeamRoles{
		Team:      *t,
		UserRoles: make(map[string]string),
	}
	var m []*github.User
	m, _ = gh.listTeamMembersByRole(*t.ID, "maintainer")
	for _, u := range m {
		tm.UserRoles[*u.Login] = "maintainer"
	}
	m, _ = gh.listTeamMembersByRole(*t.ID, "member")
	for _, u := range m {
		tm.UserRoles[*u.Login] = "member"
	}
	return tm, nil
}

// listTeamMembersByRole will get TeamMembers by Role
// role values are "all", "member", "maintainer". Default is "all".
func (gh *GitHub) listTeamMembersByRole(team int, role string) ([]*github.User, error) {
	ctx := context.Background()
	client := getGitHubTokenClient(ctx, gh.Token)
	opt := &github.OrganizationListTeamMembersOptions{
		Role:        role,
		ListOptions: github.ListOptions{PerPage: 100},
	}

	var allUsers []*github.User
	for {
		users, resp, err := client.Organizations.ListTeamMembers(ctx, team, opt)
		if err != nil {
			return nil, err
		}
		log.Debugf("GitHub Response Code: %v\n", resp.StatusCode)
		log.Debugf("GitHub Rate Remaining: %v\n", resp.Rate.Remaining)
		allUsers = append(allUsers, users...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return allUsers, nil
}

func getGitHubTokenClient(ctx context.Context, token string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}
