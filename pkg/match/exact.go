package match

import (
	"github.com/agrski/greg/pkg/fetch/github"
	"github.com/rs/zerolog"
)

type ExactMatcher struct {
	logger zerolog.Logger
}

var _ Matcher = (*ExactMatcher)(nil)
var _ FiletypeFilter = (*ExactMatcher)(nil)

func (em *ExactMatcher) Match(pattern string, next *github.FileInfo) (*Match, bool) {
	return nil, false
}

func (em *ExactMatcher) Filter(allowed []string, next *github.FileInfo) bool {
	return false
}
