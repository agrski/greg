package match

import (
	"github.com/agrski/greg/pkg/types"
	"github.com/rs/zerolog"
)

type Matcher interface {
	Match(pattern string, next *types.FileInfo) (*Match, bool)
}

type Match struct {
	Positions []*FilePosition
}

type FilePosition struct {
	Line        uint
	ColumnStart uint
	ColumnEnd   uint
	Text        string
}

type filteringMatcher struct {
	matcher   Matcher
	filetypes []types.FileExtension
	logger    zerolog.Logger
}

var _ Matcher = (*filteringMatcher)(nil)

func New(logger zerolog.Logger, allowedFiletypes []types.FileExtension) *filteringMatcher {
	em := newExactMatcher(logger)
	logger = logger.With().Str("source", "FilteringMatcher").Logger()

	return &filteringMatcher{
		matcher:   em,
		filetypes: allowedFiletypes,
		logger:    logger,
	}
}

func (fm *filteringMatcher) Match(pattern string, next *types.FileInfo) (*Match, bool) {
	if ok := FilterFiletype(fm.filetypes, next); !ok {
		return nil, false
	}
	return fm.matcher.Match(pattern, next)
}
