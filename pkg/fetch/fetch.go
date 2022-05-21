package fetch

import (
	"github.com/agrski/greg/pkg/fetch/github"
	"github.com/rs/zerolog"
	"golang.org/x/oauth2"
)

/*
	Strategy pattern that select from supported providers.
	Currently the CLI knows about supported git providers,
	but really that information should come from this package.
*/

type HostName string
type OrganisationName string
type RepositoryName string

type Location struct {
	Host         HostName
	Organisation OrganisationName
	Repository   RepositoryName
}

type Fetcher interface {
	Start() error
	Stop() error
	Next() (interface{}, bool) // TODO - formalise this return param for interop with matcher
}

func New(logger zerolog.Logger, l Location, tokenSource oauth2.TokenSource) Fetcher {
	githubFetcher := github.New(
		logger,
		github.QueryParams{
			RepoOwner: string(l.Organisation),
			RepoName:  string(l.Repository),
		},
		tokenSource,
	)
	return githubFetcher
}
