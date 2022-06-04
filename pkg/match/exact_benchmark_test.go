package match

import (
	"math/rand"
	"strings"
	"testing"

	"github.com/agrski/greg/pkg/fetch/github"
	"github.com/rs/zerolog"
)

var results int

// Exactly 64 characters to select from makes looping easier.
// 64 == 2e6 and we select from 0 to 63 inclusive.
//
// See https://stackoverflow.com/a/31832326
const (
	selectableChars      = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ \n\t,.[]{}()-"
	selectionBits        = 6
	selectionMask        = 1<<selectionBits - 1
	charsPerRandomNumber = 63 / selectionBits
)

func makeTextOfLength(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)

	cache := rand.Int63()
	remaining := charsPerRandomNumber

	for i := 0; i < n; i++ {
		if remaining == 0 {
			cache = rand.Int63()
			remaining = charsPerRandomNumber
		}

		idx := (cache & selectionMask)
		selection := selectableChars[idx]
		sb.WriteByte(selection)

		cache >>= selectionBits
		remaining--
	}

	return sb.String()
}

func BenchmarkExactMatcher(b *testing.B) {
	matcher := newExactMatcher(zerolog.Nop())
	pattern := ""
	fileInfo := &github.FileInfo{}
	fileInfo.IsBinary = false
	fileInfo.Text = ""

	for i := 0; i < b.N; i++ {
		matches, ok := matcher.Match(pattern, &github.FileInfo{})
		if ok {
			results = len(matches.Lines)
		}
	}
}