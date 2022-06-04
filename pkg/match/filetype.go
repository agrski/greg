package match

import (
	"strings"

	"github.com/agrski/greg/pkg/types"
)

func FilterFiletype(allowed []string, next *types.FileInfo) bool {
	if len(allowed) == 0 {
		return true
	}

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
