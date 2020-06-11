package main

import (
	"context"
	"net/url"

	"github.com/google/go-github/v21/github"
	"golang.org/x/oauth2"
)

const (
	defaultGithubAPIURL = "https://api.github.com"
)

type githubClient struct {
	apiURL string
	token  string
}

func (c *githubClient) getGithubAPIURL() (*url.URL, error) {
	u := defaultGithubAPIURL
	if c.apiURL != "" {
		u = c.apiURL
	}
	apiURL, err := url.Parse(u + "/")
	if err != nil {
		return nil, err
	}
	return apiURL, nil
}

func (c *githubClient) makeClient(ctx context.Context) (*github.Client, error) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: c.token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	apiURL, err := c.getGithubAPIURL()
	if err != nil {
		return nil, err
	}

	client.BaseURL = apiURL

	return client, nil
}

func (c *githubClient) searchGithubIssues(query string) ([]github.Issue, error) {
	ctx := context.Background()
	client, err := c.makeClient(ctx)
	if err != nil {
		return nil, err
	}

	i, _, err := client.Search.Issues(
		ctx,
		query,
		&github.SearchOptions{Sort: "created", Order: "asc"},
	)
	if err != nil {
		return nil, err
	}

	return i.Issues, nil
}
