package github

import (
	"context"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var ghClient *client

type client struct {
	restClient *github.Client
}

type GHClientSpec interface {
	GetLatestRelease(owner string, repo string) (*github.RepositoryRelease, error)
}

func GHClient(token string) *client {
	if ghClient != nil {
		return ghClient
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: token,
	})
	tc := oauth2.NewClient(ctx, ts)

	ghClient = &client{
		restClient: github.NewClient(tc),
	}

	return ghClient
}

func (c *client) GetLatestRelease(owner string, repo string) (*github.RepositoryRelease, error) {
	ctx := context.Background()
	release, _, err := c.restClient.Repositories.GetLatestRelease(ctx, owner, repo)

	return release, err
}
