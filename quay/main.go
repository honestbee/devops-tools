package main

import (
	"context"
	"fmt"
	"os"
	"strings"

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

func createQuayRepos(githubRepos []github.Repo) ([]quay.RepositoryOutput, error) {
	var quayRepoOuputs []quay.RepositoryOutput

	for _, githubRepo := range githubRepos {
		quayRepoInput := quay.RepositoryInput{
			Namespace:   "honestbee",
			Visibility:  "private",
			Repository:  strings.Split(githubRepo.RepoName, "/")[1],
			Description: githubRepo.RepoDescription,
		}

		quayRepoOutput, err := quayRepoInput.CreateRepository()
		if err != nil {
			fmt.Printf("Error creating %v : %v", quayRepoInput, err)
		}
		quayRepoOuputs = append(quayRepoOuputs, quayRepoOutput)
	}
	return quayRepoOuputs, nil

}

func main() {
	githubRepos, err := getGithubRepos()
	if err != nil {
		panic(err)
	}
	quayRepos, err := createQuayRepos(githubRepos)
	if err != nil {
		panic(err)
	}
	for _, quayRepo := range quayRepos {
		fmt.Printf("%v\n", quayRepo.Name)
	}
}
