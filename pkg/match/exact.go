package match

import (
	"bufio"
	"strings"

	"github.com/agrski/greg/pkg/fetch/github"
	"github.com/rs/zerolog"
)

type exactMatcher struct {
	logger zerolog.Logger
}

var _ Matcher = (*exactMatcher)(nil)

func newExactMatcher(logger zerolog.Logger) *exactMatcher {
	logger = logger.With().Str("source", "ExactMatcher").Logger()
	return &exactMatcher{
		logger: logger,
	}
}

func (em *exactMatcher) Match(pattern string, next *github.FileInfo) (*Match, bool) {
	logger := em.logger.With().Str("func", "Match").Logger()

	if next.IsBinary {
		logger.Debug().Str("filename", next.Path).Msg("rejecting binary file")
		return nil, false
	}

	match := &Match{
		Lines: []FilePosition{},
	}
	lineReader := bufio.NewScanner(
		strings.NewReader(next.Text),
	)
	row := uint(0)

	for lineReader.Scan() {
		row++
		column := strings.Index(next.Text, pattern)
		if column >= 0 {
			match.Lines = append(
				match.Lines,
				FilePosition{Line: row, Column: 1 + uint(column)},
			)
		}
	}

	if err := lineReader.Err(); err != nil {
		return nil, false
	}

	if len(match.Lines) == 0 {
		return nil, false
	}

	return match, true
}
