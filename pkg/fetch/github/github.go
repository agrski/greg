package fetch

import (
	"context"
	"time"

	"github.com/hasura/go-graphql-client"
	"golang.org/x/oauth2"
)

/*
	Candidate GraphQL clients:
	* https://github.com/hasura/go-graphql-client
		- Generic usage (not GitHub-specific)
		- Struct-based rather than string-based
	* https://github.com/shurcooL/githubv4
		- GitHub-specific
		- Struct-based
		- Supports pagination
*/

const (
	apiUrl              = "https://api.github.com/graphql"
	defaultQueryTimeout = 30 * time.Second
)

type GitHub struct {
	client *graphql.Client
}

func NewGitHub(accessToken string) *GitHub {
	// TODO - refactor OAuth handling entirely outside this package
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: accessToken,
		},
	)

	authClient := oauth2.NewClient(context.Background(), tokenSource)
	client := graphql.NewClient(apiUrl, authClient)

	return &GitHub{
		client: client,
	}
}

func (g *GitHub) MakeBaseQuery(params QueryParams) (*Query, error) {
	query := &Query{}

	variables := map[string]interface{}{
		"owner":            graphql.String(params.repoOwner),
		"repo":             graphql.String(params.repoName),
		"commitishAndPath": graphql.String(params.commitish + ":" + params.pathPrefix),
	}

	ctx, _ := context.WithTimeout(context.Background(), defaultQueryTimeout)
	err := g.client.Query(ctx, query, variables)

	return query, err
}
