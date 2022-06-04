package github

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hasura/go-graphql-client"
	"github.com/rs/zerolog"
	"golang.org/x/oauth2"

	fetchTypes "github.com/agrski/greg/pkg/fetch/types"
	"github.com/agrski/greg/pkg/types"
)

const (
	apiUrl                 = "https://api.github.com/graphql"
	defaultQueryTimeout    = 30 * time.Second
	treeResultsCapacity    = 100
	treesRemainingCapacity = 10_000 // Max subtrees we support for any node; TODO - make configurable
)

type graphqlVariables map[string]interface{}

// GitHub instances retrieve the files present as of some commit in a GitHub repository.
// An instance should only be used once, as it stores intermediate state internally.
// A stopped instance cannot be restarted cleanly; instead, create a fresh instance.
type GitHub struct {
	client      *graphql.Client
	queryParams QueryParams
	logger      zerolog.Logger
	results     <-chan *types.FileInfo
	cancel      func()
}

var _ fetchTypes.Fetcher = (*GitHub)(nil)

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
		Debug().
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
	g.logger.Debug().Str("func", "Stop").Msg("stopping GitHub fetcher")

	g.cancel()
	return nil
}

func (g *GitHub) Next() (*types.FileInfo, bool) {
	logger := g.logger.With().Str("func", "Next").Logger()
	next := <-g.results
	if next == nil {
		logger.Trace().Msg("no more results")
		return nil, false
	} else {
		logger.Trace().Msg("providing next result")
		return next, true
	}
}

func (g *GitHub) getFiles() (<-chan *types.FileInfo, func()) {
	logger := g.logger.With().Str("func", "getFiles").Logger()

	results := make(chan *types.FileInfo, treeResultsCapacity)
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
		Debug().
		Str("func", "getTree").
		Interface("commit and path", variables["commitishAndPath"]).
		Send()
	ctx, cancel := context.WithTimeout(context.Background(), defaultQueryTimeout)
	defer cancel()

	query := &treeQuery{}
	err := g.client.Query(ctx, query, variables)

	return query, err
}

func (g *GitHub) parseTree(
	tree *treeQuery,
	results chan<- *types.FileInfo,
	remaining chan<- string,
	cancel <-chan struct{},
) {
	logger := g.logger.With().Str("func", "parseTree").Logger()
	root := tree.repository.Object.Tree

	for _, e := range root.Entries {
		select {
		case <-cancel:
			return
		default:
			switch e.Type {
			case TreeEntryDir:
				remaining <- e.Path
			case TreeEntryFile:
				f := &types.FileInfo{
					Path:      e.Path,
					Extension: e.Extension,
					IsBinary:  e.Object.IsBinary,
					Text:      e.Object.Text,
				}
				results <- f
			default:
				logger.Warn().Str("type", string(e.Type)).Msg("unknown entry type")
				continue
			}
		}
	}
}
