package match

import (
	"strings"

	"github.com/agrski/greg/pkg/types"
)

func FilterFiletype(allowed []types.FileExtension, next *types.FileInfo) bool {
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

func NormaliseExtension(ext types.FileExtension) types.FileExtension {
	trimmed := strings.TrimSpace(string(ext))
	withoutDot := strings.TrimPrefix(trimmed, ".")
	return types.FileExtension(withoutDot)
}
