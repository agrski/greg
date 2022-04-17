package fetch

import (
	"github.com/agrski/gitfind/pkg/fetch/github"
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

func New(l Location, accessToken string) Fetcher {
	githubFetcher := github.New(
		github.QueryParams{
			RepoOwner: string(l.Organisation),
			RepoName:  string(l.Repository),
		},
		accessToken,
	)
	return githubFetcher
}
