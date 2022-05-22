package match

import (
	"github.com/agrski/greg/pkg/fetch/github"
	"github.com/rs/zerolog"
)

type Matcher interface {
	// FIXME - move github.FileInfo -> fetch.FileInfo
	//	as we should not be relying on something so specific.
	Match(pattern string, next *github.FileInfo) (*Match, bool)
}

type Match struct {
	line uint
}

type FilteringMatcher struct {
	matcher   Matcher
	filetypes []string
	logger    zerolog.Logger
}

var _ Matcher = (*FilteringMatcher)(nil)

func New(logger zerolog.Logger, allowedFiletypes []string) *FilteringMatcher {
	em := newExactMatcher(logger)
	logger = logger.With().Str("source", "FilteringMatcher").Logger()

	return &FilteringMatcher{
		matcher:   em,
		filetypes: allowedFiletypes,
		logger:    logger,
	}
}

func (fm *FilteringMatcher) Match(pattern string, next *github.FileInfo) (*Match, bool) {
	if ok := FilterFiletype(fm.filetypes, next); !ok {
		return nil, false
	}
	return fm.matcher.Match(pattern, next)
}
