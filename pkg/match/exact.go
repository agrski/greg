package match

import (
	"bufio"
	"strings"

	"github.com/agrski/greg/pkg/types"
	"github.com/rs/zerolog"
)

type exactMatcher struct {
	caseInsensitive bool
	logger          zerolog.Logger
}

var _ Matcher = (*exactMatcher)(nil)

func newExactMatcher(logger zerolog.Logger, caseInsensitive bool) *exactMatcher {
	logger = logger.With().Str("source", "ExactMatcher").Logger()
	return &exactMatcher{
		caseInsensitive: caseInsensitive,
		logger:          logger,
	}
}

func (em *exactMatcher) Match(pattern string, next *types.FileInfo) (*Match, bool) {
	logger := em.logger.With().Str("func", "Match").Logger()

	if next.IsBinary {
		logger.Debug().Str("filename", next.Path).Msg("rejecting binary file")
		return nil, false
	}

	match := &Match{}
	lineReader := bufio.NewScanner(
		strings.NewReader(next.Text),
	)

	for row := 0; lineReader.Scan(); row++ {
		line := lineReader.Text()
		matchColumns := em.matchLine(pattern, line)
		for _, column := range matchColumns {
			match.Positions = append(
				match.Positions,
				&FilePosition{
					Line:        uint(row),
					ColumnStart: column,
					ColumnEnd:   column + uint(len(pattern)),
					Text:        line,
				},
			)
		}
	}

	if err := lineReader.Err(); err != nil {
		return nil, false
	}

	if len(match.Positions) == 0 {
		return nil, false
	}

	return match, true
}

func (em *exactMatcher) matchLine(pattern string, line string) []uint {
	if em.caseInsensitive {
		pattern = strings.ToLower(pattern)
		line = strings.ToLower(line)
	}

	column := 0
	matchColumns := []uint{}

	for {
		offset := strings.Index(line, pattern)
		if offset == -1 {
			break
		} else {
			column += offset
			matchColumns = append(matchColumns, uint(column))

			column += len(pattern)
			line = line[offset+len(pattern):]
		}
	}

	return matchColumns
}
