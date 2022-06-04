package fetch

import (
	"github.com/agrski/greg/pkg/fetch/github"
	"github.com/agrski/greg/pkg/fetch/types"
	"github.com/rs/zerolog"
	"golang.org/x/oauth2"
)

/*
	Strategy pattern that select from supported providers.
	Currently the CLI knows about supported git providers,
	but really that information should come from this package.
*/

func New(
	logger zerolog.Logger,
	location types.Location,
	tokenSource oauth2.TokenSource,
) types.Fetcher {
	githubFetcher := github.New(
		logger,
		location,
		tokenSource,
	)
	return githubFetcher
}
