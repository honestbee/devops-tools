package github

import (
	"context"
	"os"
	"testing"
)

func TestListRepos(t *testing.T) {
	ctx := context.Background()
	githubClient := NewClient(ctx, os.Getenv("GITHUB_TOKEN"))
	testCases := []struct {
		c    Config
		want error
	}{
		{
			Config{
				Organization: "honestee",
				Client:       githubClient,
				Repo: RepoConfig{
					Type: "private",
				},
			},
			ErrNotFound,
		},
		{
			Config{
				Organization: "honestbee",
				Client:       githubClient,
				Repo: RepoConfig{
					Type: "private",
				},
			},
			nil,
		},
	}

	for _, testCase := range testCases {
		_, got := testCase.c.ListRepos()

		if got != testCase.want {
			t.Errorf("got %v want %v", got, testCase.want)
		}
	}

}
