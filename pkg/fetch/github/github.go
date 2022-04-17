package github

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
	client      *graphql.Client
	queryParams QueryParams
	results     <-chan *FileInfo
	cancel      func()
}

func New(q QueryParams, accessToken string) *GitHub {
	// TODO - refactor OAuth handling entirely outside this package
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: accessToken,
		},
	)

	authClient := oauth2.NewClient(context.Background(), tokenSource)
	client := graphql.NewClient(apiUrl, authClient)

	return &GitHub{
		client:      client,
		queryParams: q,
	}
}

func (g *GitHub) Start() error {
	results, cancel := g.getFiles()
	g.results = results
	g.cancel = cancel
	return nil
}

func (g *GitHub) Stop() error {
	g.cancel()
	return nil
}

func (g *GitHub) Next() (interface{}, bool) {
	next := <-g.results
	if next == nil {
		return nil, false
	} else {
		return next, true
	}
}

func (g *GitHub) getFiles() (<-chan *FileInfo, func()) {
	results := make(chan *FileInfo, treeResultsCapacity)
	remaining := make(chan string, treesRemainingCapacity)
	cancel := make(chan struct{})
	canceller := func() {
		close(cancel)
	}

	g.ensureCommitish()

	// Bootstrap loop with root of query
	remaining <- g.queryParams.PathPrefix

	go func() {
		defer close(results)
		defer close(remaining)

		for {
			select {
			case path := <-remaining:
				g.queryParams.PathPrefix = path
				variables := g.paramsToVariables()

				tree, err := g.getTree(variables)
				if err != nil {
					fmt.Printf("unable to fetch from GitHub: %v", err)
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

	return results, canceller
}

func (g *GitHub) ensureCommitish() error {
	if strings.TrimSpace(g.queryParams.Commitish) != "" {
		return nil
	}

	c, err := g.getDefaultBranchRef()
	if err != nil {
		return err
	}

	g.queryParams.Commitish = c
	return nil
}

func (g *GitHub) getDefaultBranchRef() (string, error) {
	q := &branchRefQuery{}
	variables := graphqlVariables{
		"owner": graphql.String(g.queryParams.RepoOwner),
		"repo":  graphql.String(g.queryParams.RepoName),
	}
	ctx, cancel := context.WithTimeout(context.Background(), defaultQueryTimeout)
	defer cancel()

	err := g.client.Query(ctx, q, variables)
	if err != nil {
		return "", err
	}

	return q.Repository.DefaultBranchRef.Name, nil
}

func (g *GitHub) paramsToVariables() graphqlVariables {
	rootExpression := g.makeRootPathExpression(g.queryParams.Commitish, g.queryParams.PathPrefix)

	variables := graphqlVariables{
		"owner":            graphql.String(g.queryParams.RepoOwner),
		"repo":             graphql.String(g.queryParams.RepoName),
		"commitishAndPath": graphql.String(rootExpression),
	}

	return variables
}

func (g *GitHub) makeRootPathExpression(commitish string, path string) string {
	return fmt.Sprintf("%s:%s", commitish, path)
}

func (g *GitHub) getTree(variables graphqlVariables) (*treeQuery, error) {
	query := &treeQuery{}
	ctx, cancel := context.WithTimeout(context.Background(), defaultQueryTimeout)
	defer cancel()
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
