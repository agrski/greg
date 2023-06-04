package match

import (
	"github.com/rs/zerolog"

	"github.com/agrski/greg/pkg/types"
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

func New(
	logger zerolog.Logger,
	caseInsensitive bool,
	allowedFiletypes []types.FileExtension,
) *filteringMatcher {
	em := newExactMatcher(logger, caseInsensitive)
	logger = logger.With().Str("source", "FilteringMatcher").Logger()

	return &filteringMatcher{
		matcher:   em,
		filetypes: allowedFiletypes,
		logger:    logger,
	}
}

func (fm *filteringMatcher) Match(pattern string, next *types.FileInfo) (*Match, bool) {
	ok := FilterFiletype(fm.filetypes, next)
	if !ok {
		return nil, false
	}

	return fm.matcher.Match(pattern, next)
}
