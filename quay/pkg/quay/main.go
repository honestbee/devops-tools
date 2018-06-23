package quay

import (
	"github.com/koudaiii/qucli/quay"
)

const (
	// Hostname default
	Hostname string = "quay.io"
)

// RepositoryInput struct defines repository api request
type RepositoryInput struct {
	Namespace  string
	Visibility string
	Name       string
	RepoKind   string
}

// RepositoryOutput struct defines repository api response
type RepositoryOutput quay.QuayRepository

// CreateRepository function
func (ri *RepositoryInput) CreateRepository() (RepositoryOutput, error) {
	var ro RepositoryOutput
	result, err := quay.CreateRepository(ri.Namespace, ri.Name, ri.Visibility, Hostname)
	if err != nil {
		return ro, err
	}
	return RepositoryOutput(result), nil
}
