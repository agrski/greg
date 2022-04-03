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
	Next() interface{} // TODO - formalise this return param for interop with matcher
}

func New(l Location) Fetcher {
	// TODO - support (GitHub) token/file
	githubFetcher := github.New(l, "")
	return githubFetcher
}
