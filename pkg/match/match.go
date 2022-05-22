package match

import (
	"github.com/agrski/greg/pkg/fetch/github"
)

type Matcher interface {
	// FIXME - move github.FileInfo -> fetch.FileInfo
	//	as we should not be relying on something so specific.
	Match(pattern string, next *github.FileInfo) (*Match, bool)
}

type Match struct {
	line uint
}
