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

type filteringMatcher struct {
	matcher   Matcher
	filetypes []string
	logger    zerolog.Logger
}

var _ Matcher = (*filteringMatcher)(nil)

func New(logger zerolog.Logger, allowedFiletypes []string) *filteringMatcher {
	em := newExactMatcher(logger)
	logger = logger.With().Str("source", "FilteringMatcher").Logger()

	return &filteringMatcher{
		matcher:   em,
		filetypes: allowedFiletypes,
		logger:    logger,
	}
}

func (fm *filteringMatcher) Match(pattern string, next *github.FileInfo) (*Match, bool) {
	if ok := FilterFiletype(fm.filetypes, next); !ok {
		return nil, false
	}
	return fm.matcher.Match(pattern, next)
}
