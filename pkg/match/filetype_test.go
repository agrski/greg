package match

import (
	"testing"

	"github.com/agrski/greg/pkg/types"
	"github.com/stretchr/testify/require"
)

func TestNormaliseExtension(t *testing.T) {
	type test struct {
		name      string
		extension types.FileExtension
		expected  types.FileExtension
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

func TestFilterFiletype(t *testing.T) {
	type test struct {
		name      string
		allowed   []types.FileExtension
		extension types.FileExtension
		expected  bool
	}

	tests := []test{
		{
			name:      "should allow everything when allowed is nil",
			allowed:   nil,
			extension: "md",
			expected:  true,
		},
		{
			name:      "should allow everything when allowed is empty",
			allowed:   []types.FileExtension{},
			extension: "md",
			expected:  true,
		},
		{
			name:      "single-element filter matches",
			allowed:   []types.FileExtension{"md"},
			extension: "md",
			expected:  true,
		},
		{
			name:      "multi-element filter matches",
			allowed:   []types.FileExtension{"a", "b", "md", "c"},
			extension: "md",
			expected:  true,
		},
		{
			name:      "single-element filter does not match",
			allowed:   []types.FileExtension{"go"},
			extension: "md",
			expected:  false,
		},
		{
			name:      "multi-element filter does not match",
			allowed:   []types.FileExtension{"go", "py", "sh"},
			extension: "md",
			expected:  false,
		},
		{
			name:      "should not match when extension is prefix of an allowed file-type",
			allowed:   []types.FileExtension{"pyc"},
			extension: "py",
			expected:  false,
		},
		{
			name:      "should not match when extension is substring of an allowed file-type",
			allowed:   []types.FileExtension{"numpy"},
			extension: "py",
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fileInfo := &types.FileInfo{}
			fileInfo.Extension = tt.extension

			actual := FilterFiletype(tt.allowed, fileInfo)

			require.Equal(t, tt.expected, actual)
		})
	}
}
