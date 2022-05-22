package match

import (
	"testing"

	"github.com/agrski/greg/pkg/fetch/github"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestMatch(t *testing.T) {
	type test struct {
		name       string
		isBinary   bool
		text       string
		pattern    string
		expected   *Match
		expectedOk bool
	}

	tests := []test{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fileInfo := &github.FileInfo{}
			fileInfo.IsBinary = tt.isBinary
			fileInfo.Text = tt.text

			matcher := newExactMatcher(zerolog.Nop())

			actual, ok := matcher.Match(tt.pattern, fileInfo)

			require.Equal(t, tt.expectedOk, ok)
			require.Equal(t, actual, tt.expected)
		})
	}
}
