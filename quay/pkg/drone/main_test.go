package drone

import (
	"fmt"
	"testing"
)

func TestNewClient(t *testing.T) {
	client, _ := NewClient()
	repos, _ := client.RepoList()

	for _, repo := range repos {
		fmt.Println(repo)
	}
}
