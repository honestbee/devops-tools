package main

import (
	"context"
	"fmt"
	"os"

	gh "github.com/honestbee/devops-tools/quay/pkg/github"
)

func main() {
	githubToken := os.Getenv("GITHUB_TOKEN")
	ctx := context.Background()
	ghClient := gh.GetGithubTokenClient(ctx, githubToken)

	gitStruct := gh.Github{
		Organization: "honestbee",
		Client:       ghClient,
	}

	repos, err := gitStruct.ListRepos()

	for _, repo := range repos {
		fmt.Printf("repo name: %v - repo description: %v", repo.RepoName, repo.RepoDescription)
	}

	if err != nil {
		fmt.Println(err)
	}
}
