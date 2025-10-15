package github_client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/config"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/github_client_gen"
)

type (
	githubClientImpl struct {
		client        *github_client_gen.ClientWithResponses
		authorization string
		apiVersion    string
	}

	GitHubClient interface {
		GetPullRequest(ctx context.Context, owner, repo string, pullNumber int) (*github_client_gen.PullRequestResponse, error)
		GetPullRequestFiles(ctx context.Context, owner, repo string, pullNumber int) (*github_client_gen.PullRequestFilesResponse, error)
		IsPullRequestMerged(ctx context.Context, owner, repo string, pullNumber int) (bool, error)
	}
)

func NewGitHubClient(config *config.Config) (GitHubClient, error) {
	client, err := github_client_gen.NewClientWithResponses("https://api.github.com")
	if err != nil {
		return nil, err
	}
	return &githubClientImpl{
		client:        client,
		authorization: fmt.Sprintf("Bearer %s", config.GitHubToken),
		apiVersion:    "2022-11-28",
	}, nil
}

func (c *githubClientImpl) GetPullRequest(ctx context.Context, owner, repo string, pullNumber int) (*github_client_gen.PullRequestResponse, error) {
	params := &github_client_gen.GetReposOwnerRepoPullsPullNumberParams{
		XGitHubApiVersion: c.apiVersion,
	}

	response, err := c.client.GetReposOwnerRepoPullsPullNumberWithResponse(ctx, owner, repo, pullNumber, params, func(ctx context.Context, req *http.Request) error {
		req.Header.Set("Authorization", c.authorization)
		return nil
	})
	if err != nil {
		return nil, err
	}

	switch response.HTTPResponse.StatusCode {
	case http.StatusNotFound:
		return nil, fmt.Errorf("GitHub GetPullRequest: not found")
	case http.StatusOK:
		if response.JSON200 == nil {
			return nil, fmt.Errorf("200 response is nil")
		}
		return response.JSON200, nil
	default:
		return nil, fmt.Errorf("unexpected code: %d", response.HTTPResponse.StatusCode)
	}
}

func (c *githubClientImpl) GetPullRequestFiles(ctx context.Context, owner, repo string, pullNumber int) (*github_client_gen.PullRequestFilesResponse, error) {
	params := &github_client_gen.GetReposOwnerRepoPullsPullNumberFilesParams{
		XGitHubApiVersion: c.apiVersion,
	}

	response, err := c.client.GetReposOwnerRepoPullsPullNumberFilesWithResponse(ctx, owner, repo, pullNumber, params, func(ctx context.Context, req *http.Request) error {
		req.Header.Set("Authorization", c.authorization)
		return nil
	})
	if err != nil {
		return nil, err
	}

	switch response.HTTPResponse.StatusCode {
	case http.StatusNotFound:
		return nil, fmt.Errorf("GitHub GetPullRequestFiles: not found")
	case http.StatusOK:
		if response.JSON200 == nil {
			return nil, fmt.Errorf("200 response is nil")
		}
		return response.JSON200, nil
	default:
		return nil, fmt.Errorf("unexpected code: %d", response.HTTPResponse.StatusCode)
	}
}

func (c *githubClientImpl) IsPullRequestMerged(ctx context.Context, owner, repo string, pullNumber int) (bool, error) {
	params := &github_client_gen.GetReposOwnerRepoPullsPullNumberMergeParams{
		XGitHubApiVersion: c.apiVersion,
	}

	response, err := c.client.GetReposOwnerRepoPullsPullNumberMergeWithResponse(ctx, owner, repo, pullNumber, params, func(ctx context.Context, req *http.Request) error {
		req.Header.Set("Authorization", c.authorization)
		return nil
	})
	if err != nil {
		return false, err
	}

	switch response.HTTPResponse.StatusCode {
	case http.StatusNotFound:
		return false, nil
	case http.StatusNoContent:
		return true, nil
	default:
		return false, fmt.Errorf("unexpected code: %d", response.HTTPResponse.StatusCode)
	}
}
