package codecov

import (
	"context"
	"net/http"
)

type (
	// ListRepositoriesResponse contains the response of a ListRepositories call
	ListRepositoriesResponse struct {
		Repos []Repository
		Meta  Meta
	}

	// GetRepositoryResponse contains the response of a GetRepository call
	GetRepositoryResponse struct {
		Repo Repository
		Meta Meta
	}

	// Repository is a tracked Codecov VCS repository
	Repository struct {
		Name             string  // The name of the repository
		Language         string  // The primary programming language of the repository
		Activated        bool    // If this repository is active in Codecov
		Deleted          bool    // If the repository has been deleted
		Private          bool    // If this is a private repository
		Updatestamp      *Time   // The last time this data was updated
		Branch           string  // The branch being monitored
		Coverage         float64 // The current coverage amount
		RepoID           string  // The internal Codecov ID for the repository
		UsingIntegration bool    `json:"using_integration"`
	}
)

// ListRepositories will list all repositories for the given owner
// https://docs.codecov.io/reference#repositories
func (c *Client) ListRepositories(ctx context.Context, owner string) (response ListRepositoriesResponse, _ error) {
	request, err := http.NewRequestWithContext(ctx, "GET", "/gh/"+owner, nil)
	if err != nil {
		return response, err
	}

	err = c.doRequest(request, &response)
	return response, err
}

// GetRepository will get a single repository for the given owner
// https://docs.codecov.io/reference#repositories
func (c *Client) GetRepository(ctx context.Context, owner string, repo string) (response GetRepositoryResponse, _ error) {
	request, err := http.NewRequestWithContext(ctx, "GET", "/gh/"+owner+"/"+repo, nil)
	if err != nil {
		return response, err
	}

	err = c.doRequest(request, &response)
	return response, err
}
