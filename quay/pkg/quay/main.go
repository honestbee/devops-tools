package quay

import (
	"encoding/json"
	"os"
	"path"

	"github.com/koudaiii/qucli/quay"
	"github.com/koudaiii/qucli/utils"
)

const (
	// Hostname default
	Hostname string = "quay.io"
)

// Config alias
type Config struct {
	Hostname string `json:"hostname"`
	APIToken string `json:"api_token"`
}

var config *Config

// RepositoryInput struct defines repository api request
type RepositoryInput struct {
	Namespace   string `json:"namespace"`
	Visibility  string `json:"visibility"`
	Repository  string `json:"repository"`
	Description string `json:"description"`
}

// RepositoryOutput struct defines repository api response
type RepositoryOutput quay.QuayRepository

func init() {
	config = &Config{
		Hostname: "quay.io",
		APIToken: os.Getenv("QUAY_API_TOKEN"),
	}
}

// CreateRepository function
func (ri *RepositoryInput) CreateRepository() (RepositoryOutput, error) {
	var ro RepositoryOutput
	req, err := json.Marshal(ri)

	u := quay.QuayURLParse(config.Hostname)
	u.Path = path.Join(u.Path, "repository")

	body, err := utils.HttpPost(u.String(), config.APIToken, req)
	if err != nil {
		return ro, err
	}

	if err := json.Unmarshal([]byte(body), &ro); err != nil {
		return ro, err
	}

	return ro, nil
}
