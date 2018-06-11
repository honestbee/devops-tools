package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
	"github.com/urfave/cli"
	"golang.org/x/oauth2"
)

type Github struct {
}

func NewGithubClient(c *cli.Context) (context.Context, *github.Client) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: c.String("github-token")})
	tc := oauth2.NewClient(ctx, ts)
	return ctx, github.NewClient(tc)
}

func (g Github) AddUser(c *cli.Context) {
	ctx, client := NewGithubClient(c)
	_, _, err := client.Organizations.EditOrgMembership(ctx, c.Args().Get(0), "honestbee", &github.Membership{})
	if err != nil {
		fmt.Println(err)
	}
}

func (g Github) ListUsers(c *cli.Context) {
	ctx, client := NewGithubClient(c)
	users, _, err := client.Organizations.ListMembers(ctx, "honestbee", &github.ListMembersOptions{
		ListOptions: github.ListOptions{
			Page:    1,
			PerPage: 1000,
		},
	})
	if err != nil {
		fmt.Println(err)
	}

	for _, user := range users {
		fmt.Println(*user.Login)
	}
}

func (g Github) DeleteUser(c *cli.Context) {
	ctx, client := NewGithubClient(c)
	_, err := client.Organizations.RemoveOrgMembership(ctx, c.Args().Get(0), "honestbee")
	if err != nil {
		fmt.Println(err)
	}
}

func ListUserTeams(c *cli.Context, ctx context.Context, client *github.Client) []*github.Team {
	// https://godoc.org/github.com/google/go-github/github#OrganizationsService.ListUserTeams
	var teamList []*github.Team
	teams, _, err := client.Organizations.ListTeams(ctx, "honestbee", &github.ListOptions{})
	for _, team := range teams {
		isTeamMember, _, err := client.Organizations.IsTeamMember(ctx, team.GetID(), c.Args().Get(0))
		if err != nil {
			fmt.Println(err)
		}
		if isTeamMember {
			fmt.Println(*team.Name)
			teamList = append(teamList, team)
		}
	}
	if err != nil {
		fmt.Println(err)
	}
	return teamList
}

func (g Github) RemoveUserFromTeams(c *cli.Context) {
	// https://godoc.org/github.com/google/go-github/github#OrganizationsService.RemoveTeamMembership
	ctx, client := NewGithubClient(c)
	teams := ListUserTeams(c, ctx, client)
	for _, team := range teams {
		_, err := client.Organizations.RemoveTeamMembership(ctx, team.GetID(), c.Args().Get(0))
		if err != nil {
			fmt.Println(err)
		}
	}
}
