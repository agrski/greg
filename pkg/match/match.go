package match

import (
	"github.com/agrski/greg/pkg/fetch/github"
)

type Matcher interface {
	Match(pattern string, next *github.FileInfo) (*Match, bool)
}

type Match struct {
	line uint
}
