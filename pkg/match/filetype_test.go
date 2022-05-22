package match

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNormaliseExtension(t *testing.T) {
	type test struct {
		name      string
		extension string
		expected  string
	}

	tests := []test{
		{
			name:      "empty extension",
			extension: "",
			expected:  "",
		},
		{
			name:      "only whitespace",
			extension: "\t  \n",
			expected:  "",
		},
		{
			name:      "no leading dot",
			extension: "md",
			expected:  "md",
		},
		{
			name:      "with leading dot",
			extension: ".md",
			expected:  "md",
		},
		{
			name:      "with leading dot and trailing whitespace",
			extension: ".md\t",
			expected:  "md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := NormaliseExtension(tt.extension)
			require.Equal(t, tt.expected, actual)
		})
	}
}
