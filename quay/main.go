package main

import (
	"context"
	"fmt"
	"os"

	drone_go "github.com/drone/drone-go/drone"
	"github.com/honestbee/devops-tools/quay/pkg/drone"
	"github.com/honestbee/devops-tools/quay/pkg/github"
	"github.com/honestbee/devops-tools/quay/pkg/quay"
	"github.com/tuannvm/tools/pkg/utils"
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

func createQuayRepos(droneRepos []*drone_go.Repo) ([]quay.RepositoryOutput, error) {
	var quayRepoOuputs []quay.RepositoryOutput

	for _, droneRepo := range droneRepos {
		quayRepoInput := quay.RepositoryInput{
			Namespace:   "honestbee",
			Visibility:  droneRepo.Visibility,
			Repository:  droneRepo.Name,
			Description: "",
		}

		quayRepoOutput, err := quayRepoInput.CreateRepository()
		if err != nil {
			fmt.Printf("Error creating %v : %v", quayRepoInput, err)
		}
		quayRepoOuputs = append(quayRepoOuputs, quayRepoOutput)
	}
	return quayRepoOuputs, nil

}

func droneRegistryCreate(c drone.Client, hostname, username, password string, repo *drone_go.Repo) error {
	registry := &drone_go.Registry{
		Address:  hostname,
		Username: username,
		Password: password,
	}
	_, err := c.RegistryCreate(repo.Owner, repo.Name, registry)
	if err != nil {
		return err
	}
	return nil
}

func saveToCsv(droneRepos []drone_go.Repo) error {
	file, err := os.Create("repos.csv")
	if err != nil {
		panic(err)
	}

	defer file.Close()

	data := [][]string{
		{"Name", "Created"},
	}

	for _, droneRepo := range droneRepos {
		data = append(data, []string{droneRepo.FullName, "true"})
	}

	utils.SaveToCsv(file, data)
	return nil
}

func main() {

	droneClient, _ := drone.NewClient()
	droneRepos, _ := droneClient.RepoList()

	hostname := "quay.io"
	username := os.Getenv("DRONE_USERNAME")
	password := os.Getenv("DRONE_PASSWORD")

	for _, droneRepo := range droneRepos {
		err := droneRegistryCreate(droneClient, hostname, username, password, droneRepo)
		if err != nil {
			fmt.Println(err)
		}
	}
	//	quayRepos, err := createQuayRepos(droneRepos)
	//	if err != nil {
	//		panic(err)
	//	}
	//	for _, quayRepo := range quayRepos {
	//		fmt.Printf("%v\n", quayRepo.Name)
	//	}
}
