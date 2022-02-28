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

const (
	TreeEntryDir  = "tree"
	TreeEntryFile = "blob"
)

type gitHub struct {
	client *graphql.Client
}

type queryParams struct {
	repoOwner  string
	repoName   string
	commitish  string
	pathPrefix string
}

func NewGitHub(accessToken string) *gitHub {
	// TODO - refactor OAuth handling entirely outside this package
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: accessToken,
		},
	)

	authClient := oauth2.NewClient(context.Background(), tokenSource)
	client := graphql.NewClient(apiUrl, authClient)

	return &gitHub{
		client: client,
	}
}

func (g *gitHub) makeBaseQuery(params queryParams) (*Query, error) {
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

type Query struct {
	Repository struct {
		Name   string
		Object struct {
			FileInfoFragment `graphql:"... on Tree"`
		} `graphql:"object(expression: $commitishAndPath)"`
	} `graphql:"repository(owner: $owner, name: $repo)"`
}

type FileInfoFragment struct {
	Entries []FileInfo
}

type FileInfo struct {
	Name      string
	Type      string
	Extension string
}
