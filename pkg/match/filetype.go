package match

import (
	"strings"

	"github.com/agrski/greg/pkg/fetch/github"
)

func FilterFiletype(allowed []string, next *github.FileInfo) bool {
	normalised := NormaliseExtension(next.Extension)
	for _, a := range allowed {
		if normalised == a {
			return true
		}
	}
	return false
}

func NormaliseExtension(ext string) string {
	trimmed := strings.TrimSpace(ext)
	return strings.TrimPrefix(trimmed, ".")
}
