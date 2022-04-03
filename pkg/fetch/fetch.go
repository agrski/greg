package fetch

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
	Next() interface{}
}

func New(l Location) Fetcher {
	return nil
}
