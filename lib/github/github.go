package github

import (
	"context"
	"time"

	"golang.org/x/oauth2"

	"github.com/shurcooL/githubv4"

	"github.com/philips-labs/tabia/lib/github/graphql"
)

type Client struct {
	*githubv4.Client
}

func NewClientWithTokenAuth(token string) *Client {
	src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	httpClient := oauth2.NewClient(context.Background(), src)

	client := githubv4.NewClient(httpClient)

	return &Client{client}
}

//go:generate stringer -type=Visibility

// Visibility indicates repository visibility
type Visibility int

const (
	// Public repositories are publicly visible
	Public Visibility = iota
	// Internal repositories are only visible to organization members
	Internal
	// Private repositories are only visible to authorized users
	Private
)

type Repository struct {
	Name       string     `json:"name,omitempty"`
	ID         string     `json:"id,omitempty"`
	URL        string     `json:"url,omitempty"`
	SSHURL     string     `json:"ssh_url,omitempty"`
	Owner      string     `json:"owner,omitempty"`
	Visibility Visibility `json:"is_private,omitempty"`
	CreatedAt  time.Time  `json:"created_at,omitempty"`
	UpdatedAt  time.Time  `json:"updated_at,omitempty"`
	PushedAt   time.Time  `json:"pushed_at,omitempty"`
}

func (c *Client) FetchOrganziationRepositories(ctx context.Context, owner string) ([]Repository, error) {
	var q struct {
		Organization graphql.Organization `graphql:"organization(login: $owner)"`
	}

	variables := map[string]interface{}{
		"owner":      githubv4.String(owner),
		"repoCursor": (*githubv4.String)(nil),
	}

	var repositories []graphql.Repository
	for {
		err := c.Query(ctx, &q, variables)
		if err != nil {
			return nil, err
		}

		repositories = append(repositories, q.Organization.Repositories.Nodes...)
		if !q.Organization.Repositories.PageInfo.HasNextPage {
			break
		}
		variables["repoCursor"] = githubv4.NewString(q.Organization.Repositories.PageInfo.EndCursor)
	}

	return Map(repositories), nil
}

func Map(repositories []graphql.Repository) []Repository {
	repos := make([]Repository, len(repositories))
	for i, repo := range repositories {
		repos[i] = Repository{
			Name:      repo.Name,
			ID:        repo.ID,
			URL:       repo.URL,
			SSHURL:    repo.SSHURL,
			Owner:     repo.Owner.Login,
			CreatedAt: repo.CreatedAt,
			UpdatedAt: repo.UpdatedAt,
			PushedAt:  repo.PushedAt,
		}

		if repo.IsPrivate {
			repos[i].Visibility = Private
		}
	}

	return repos
}
