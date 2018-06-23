package main

import (
	"context"
	"fmt"
	"os"

	github "github.com/honestbee/devops-tools/quay/pkg/github"
)

func main() {
	githubToken := os.Getenv("GITHUB_TOKEN")
	ctx := context.Background()
	githubClient := github.GetGithubTokenClient(ctx, githubToken)
	getGithubRepos("honestbee", githubClient)
}

func getGithubRepos(org string, client github.Client) error {
	gitStruct := github.Github{
		Organization: org,
		Client:       client,
	}

	repos, err := gitStruct.ListRepos()

	for _, repo := range repos {
		fmt.Printf("repo name: %v - repo description: %v", repo.RepoName, repo.RepoDescription)
	}

	if err != nil {
		return err
	}
	return nil
}
