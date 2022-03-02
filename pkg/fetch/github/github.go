package fetch

import (
	"context"
	"fmt"
	"strings"
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

type graphqlVariables map[string]interface{}

type FileInfo struct {
	FileMetadata
	FileContents
}

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

func (g *GitHub) ListFiles(params QueryParams) ([]FileInfo, error) {
	g.ensureCommitish(&params)
	variables := g.paramsToVariables(params)

	_, err := g.getTree(variables)
	if err != nil {
		return nil, err
	}

	// TODO - parse tree listing into list of file details/contents

	return nil, nil
}

func (g *GitHub) ensureCommitish(params *QueryParams) error {
	if strings.TrimSpace(params.Commitish) != "" {
		return nil
	}

	c, err := g.getDefaultBranchRef(params.RepoOwner, params.RepoName)
	if err != nil {
		return err
	}

	params.Commitish = c
	return nil
}

func (g *GitHub) getDefaultBranchRef(owner string, repo string) (string, error) {
	q := &branchRefQuery{}
	variables := graphqlVariables{
		"owner": graphql.String(owner),
		"repo":  graphql.String(repo),
	}
	ctx, _ := context.WithTimeout(context.Background(), defaultQueryTimeout)

	err := g.client.Query(ctx, q, variables)
	if err != nil {
		return "", err
	}

	return q.Repository.DefaultBranchRef.Name, nil
}

func (g *GitHub) paramsToVariables(params QueryParams) graphqlVariables {
	rootExpression := g.makeRootPathExpression(params.Commitish, params.PathPrefix)

	variables := graphqlVariables{
		"owner":            graphql.String(params.RepoOwner),
		"repo":             graphql.String(params.RepoName),
		"commitishAndPath": graphql.String(rootExpression),
	}

	return variables
}

func (g *GitHub) makeRootPathExpression(commitish string, path string) string {
	return fmt.Sprintf("%s:%s", commitish, path)
}

func (g *GitHub) getTree(variables graphqlVariables) (*treeQuery, error) {
	query := &treeQuery{}
	ctx, _ := context.WithTimeout(context.Background(), defaultQueryTimeout)
	err := g.client.Query(ctx, query, variables)

	return query, err
}
