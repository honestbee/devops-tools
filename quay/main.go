package main

import (
	"context"
	"os"

	"github.com/honestbee/devops-tools/quay/pkg/github"
	"github.com/honestbee/devops-tools/quay/pkg/quay"
)

func getGithubRepos() ([]github.Repo, error) {
	githubToken := os.Getenv("GITHUB_TOKEN")
	ctx := context.Background()
	githubClient := github.NewClient(ctx, githubToken)

	githubConfig := github.Config{
		Organization: "honestbee",
		Client:       githubClient,
		Repo: github.RepoConfig{
			Type: "private",
		},
	}

	repos, err := githubConfig.ListRepos()
	return repos, err

}

func createQuayRepos(githubRepos []github.Repo) []quay.RepositoryOutput {
	for _, githubRepo := range githubRepos {
		quay.
			strings.Split(githubRepo, "/")[1]
	}

}
