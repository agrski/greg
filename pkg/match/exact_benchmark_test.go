package match

import (
	"math/rand"
	"strings"
	"testing"

	"github.com/agrski/greg/pkg/types"
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

func benchmarkExactMatcher(b *testing.B, patternSize int, textSize int, caseInsensitive bool) {
	matcher := newExactMatcher(zerolog.Nop(), caseInsensitive)
	pattern := makeTextOfLength(patternSize)
	fileInfo := &types.FileInfo{}
	fileInfo.IsBinary = false
	fileInfo.Text = makeTextOfLength(textSize)

	for i := 0; i < b.N; i++ {
		matches, ok := matcher.Match(pattern, fileInfo)
		if ok {
			results = len(matches.Positions)
		}
	}
}

func BenchmarkExactMatcher_Pattern10_Text100(b *testing.B) {
	benchmarkExactMatcher(b, 10, 100, false)
}
func BenchmarkExactMatcher_Pattern10_Text100_CaseInsensitive(b *testing.B) {
	benchmarkExactMatcher(b, 10, 100, true)
}

func BenchmarkExactMatcher_Pattern10_Text1_000(b *testing.B) {
	benchmarkExactMatcher(b, 10, 1_000, false)
}
func BenchmarkExactMatcher_Pattern100_Text1_000(b *testing.B) {
	benchmarkExactMatcher(b, 100, 1_000, false)
}
func BenchmarkExactMatcher_Pattern10_Text1_000_CaseInsensitive(b *testing.B) {
	benchmarkExactMatcher(b, 10, 1_000, true)
}
func BenchmarkExactMatcher_Pattern100_Text1_000_CaseInsensitive(b *testing.B) {
	benchmarkExactMatcher(b, 100, 1_000, true)
}

func BenchmarkExactMatcher_Pattern10_Text10_000(b *testing.B) {
	benchmarkExactMatcher(b, 10, 10_000, false)
}
func BenchmarkExactMatcher_Pattern100_Text10_000(b *testing.B) {
	benchmarkExactMatcher(b, 100, 10_000, false)
}
func BenchmarkExactMatcher_Pattern1_000_Text10_000(b *testing.B) {
	benchmarkExactMatcher(b, 1_000, 10_000, false)
}
func BenchmarkExactMatcher_Pattern10_Text10_000_CaseInsensitive(b *testing.B) {
	benchmarkExactMatcher(b, 10, 10_000, true)
}
func BenchmarkExactMatcher_Pattern100_Text10_000_CaseInsensitive(b *testing.B) {
	benchmarkExactMatcher(b, 100, 10_000, true)
}
func BenchmarkExactMatcher_Pattern1_000_Text10_000_CaseInsensitive(b *testing.B) {
	benchmarkExactMatcher(b, 1_000, 10_000, true)
}

func BenchmarkExactMatcher_Pattern10_Text100_000(b *testing.B) {
	benchmarkExactMatcher(b, 10, 100_000, false)
}
func BenchmarkExactMatcher_Pattern100_Text100_000(b *testing.B) {
	benchmarkExactMatcher(b, 100, 100_000, false)
}
func BenchmarkExactMatcher_Pattern1_000_Text100_000(b *testing.B) {
	benchmarkExactMatcher(b, 1_000, 100_000, false)
}
func BenchmarkExactMatcher_Pattern10_000_Text100_000(b *testing.B) {
	benchmarkExactMatcher(b, 10_000, 100_000, false)
}
func BenchmarkExactMatcher_Pattern10_Text100_000_CaseInsensitive(b *testing.B) {
	benchmarkExactMatcher(b, 10, 100_000, true)
}
func BenchmarkExactMatcher_Pattern100_Text100_000_CaseInsensitive(b *testing.B) {
	benchmarkExactMatcher(b, 100, 100_000, true)
}
func BenchmarkExactMatcher_Pattern1_000_Text100_000_CaseInsensitive(b *testing.B) {
	benchmarkExactMatcher(b, 1_000, 100_000, true)
}
func BenchmarkExactMatcher_Pattern10_000_Text100_000_CaseInsensitive(b *testing.B) {
	benchmarkExactMatcher(b, 10_000, 100_000, true)
}
