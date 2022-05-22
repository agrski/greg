package match

import (
	"github.com/agrski/greg/pkg/fetch/github"
	"github.com/rs/zerolog"
)

type ExactMatcher struct {
	logger zerolog.Logger
}

var _ Matcher = (*ExactMatcher)(nil)

func newExactMatcher(logger zerolog.Logger) *ExactMatcher {
	logger = logger.With().Str("source", "ExactMatcher").Logger()
	return &ExactMatcher{
		logger: logger,
	}
}

func (em *ExactMatcher) Match(pattern string, next *github.FileInfo) (*Match, bool) {
	return nil, false
}
