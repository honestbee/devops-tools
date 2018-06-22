package github

import (
	"context"
	"os"
	"testing"
)

func TestListRepos(t *testing.T) {
	ctx := context.Background()
	githubClient := GetGithubTokenClient(ctx, os.Getenv("GITHUB_TOKEN"))
	githubTests := []struct {
		gh   Github
		want error
	}{
		{
			Github{
				Organization: "honestee",
				Client:       githubClient,
			},
			ErrNotFound,
		},
		{
			Github{
				Organization: "honestbee",
				Client:       githubClient,
			},
			nil,
		},
	}

	for _, githubTest := range githubTests {
		_, got := githubTest.gh.ListRepos()

		if got != githubTest.want {
			t.Errorf("got %v want %v", got, githubTest.want)
		}
	}

}
