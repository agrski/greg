package github

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hasura/go-graphql-client"
	"github.com/rs/zerolog"
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

// GitHub instances retrieve the files present as of some commit in a GitHub repository.
// An instance should only be used once, as it stores intermediate state internally.
// A stopped instance cannot be restarted cleanly; instead, create a fresh instance.
type GitHub struct {
	client      *graphql.Client
	queryParams QueryParams
	logger      zerolog.Logger
	results     <-chan *FileInfo
	cancel      func()
}

func New(l zerolog.Logger, q QueryParams, tokenSource oauth2.TokenSource) *GitHub {
	authClient := oauth2.NewClient(context.Background(), tokenSource)
	client := graphql.NewClient(apiUrl, authClient)
	logger := l.With().Str("source", "GitHub").Logger()

	return &GitHub{
		logger:      logger,
		client:      client,
		queryParams: q,
	}
}

func (g *GitHub) Start() error {
	g.logger.
		Info().
		Str("func", "Start").
		Dur("query timeout", defaultQueryTimeout).
		Int("fetch capacity", treesRemainingCapacity).
		Int("result capacity", treeResultsCapacity).
		Str("org", g.queryParams.RepoOwner).
		Str("repo", g.queryParams.RepoName).
		Msg("starting GitHub fetcher")

	results, cancel := g.getFiles()
	g.results = results
	g.cancel = cancel
	return nil
}

func (g *GitHub) Stop() error {
	g.logger.Info().Str("func", "Stop").Msg("stopping GitHub fetcher")

	g.cancel()
	return nil
}

func (g *GitHub) Next() (interface{}, bool) {
	logger := g.logger.With().Str("func", "Next").Logger()
	next := <-g.results
	if next == nil {
		logger.Debug().Msg("no more results")
		return nil, false
	} else {
		logger.Debug().Msg("providing next result")
		return next, true
	}
}

func (g *GitHub) getFiles() (<-chan *FileInfo, func()) {
	logger := g.logger.With().Str("func", "getFiles").Logger()

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
					logger.Error().Err(err).Msg("unable to fetch from GitHub")
					return
				}

				parseTree(tree, results, remaining, cancel)
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
	rootExpression := g.makeRootPathExpression()

	variables := graphqlVariables{
		"owner":            graphql.String(g.queryParams.RepoOwner),
		"repo":             graphql.String(g.queryParams.RepoName),
		"commitishAndPath": graphql.String(rootExpression),
	}

	return variables
}

func (g *GitHub) makeRootPathExpression() string {
	return fmt.Sprintf("%s:%s", g.queryParams.Commitish, g.queryParams.PathPrefix)
}

func (g *GitHub) getTree(variables graphqlVariables) (*treeQuery, error) {
	g.logger.
		Info().
		Str("func", "getTree").
		Interface("commit and path", variables["commitishAndPath"]).
		Send()
	ctx, cancel := context.WithTimeout(context.Background(), defaultQueryTimeout)
	defer cancel()

	query := &treeQuery{}
	err := g.client.Query(ctx, query, variables)

	return query, err
}

func parseTree(
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
