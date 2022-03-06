package fetch

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hasura/go-graphql-client"
	"golang.org/x/oauth2"
)

const (
	apiUrl                 = "https://api.github.com/graphql"
	defaultQueryTimeout    = 30 * time.Second
	treeResultsCapacity    = 100
	treesRemainingCapacity = 10_000 // Max subtrees we support for any node; TODO - make configurable
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

func (g *GitHub) GetFiles(params QueryParams) (<-chan *FileInfo, func(), error) {
	results := make(chan *FileInfo, treeResultsCapacity)
	remaining := make(chan string, treesRemainingCapacity)
	cancel := make(chan struct{})
	canceller := func() {
		close(cancel)
	}

	g.ensureCommitish(&params)

	// Bootstrap loop with root of query
	remaining <- params.PathPrefix

	go func() {
		defer close(results)
		defer close(remaining)

		for {
			select {
			case path := <-remaining:
				params.PathPrefix = path
				variables := g.paramsToVariables(params)

				tree, err := g.getTree(variables)
				if err != nil {
					// TODO - log warning or exit?
					return
				}

				g.parseTree(tree, results, remaining, cancel)
			case <-cancel:
				return
			default:
				return
			}
		}
	}()

	return results, canceller, nil
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

func (g *GitHub) parseTree(
	tree *treeQuery,
	results chan<- *FileInfo,
	remaining chan<- string,
	cancel <-chan struct{},
) {
	root := tree.Repository.Object.Tree

	for _, e := range root.Entries {
		select {
		case <-cancel:
			return
		default:
			switch e.Type {
			case TreeEntryDir:
				remaining <- e.Path
			case TreeEntryFile:
				f := &FileInfo{
					FileMetadata: e.FileMetadata,
					FileContents: e.Object.FileContents,
				}
				results <- f
			default:
				// TODO - log error
				continue
			}
		}
	}
}
