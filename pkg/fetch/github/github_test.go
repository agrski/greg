//go:build !integration

package github

import (
	"testing"

	"github.com/agrski/greg/pkg/types"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestParseTree(t *testing.T) {
	type treeEntry struct {
		fileMetadata
		Object struct {
			fileContents "graphql:\"... on Blob\""
		}
	}

	type test struct {
		name              string
		entries           []entry
		expectedResults   []*types.FileInfo
		expectedRemaining []string
	}

	tests := []test{
		{
			name:              "empty root dir",
			entries:           []entry{},
			expectedResults:   []*types.FileInfo{},
			expectedRemaining: []string{},
		},
		{
			name: "one empty directory",
			entries: []entry{
				{
					fileMetadata{
						Type: TreeEntryDir,
						Name: "dir1",
						Path: "dir1",
					},
					entryObject{
						fileContents{},
					},
				},
			},
			expectedResults:   []*types.FileInfo{},
			expectedRemaining: []string{"dir1"},
		},
		{
			name: "one file in root dir",
			entries: []entry{
				{
					fileMetadata{
						Type:      TreeEntryFile,
						Name:      "file1.txt",
						Extension: ".txt",
					},
					entryObject{
						fileContents{
							IsBinary: false,
							Text:     "some text",
						},
					},
				},
			},
			expectedResults: []*types.FileInfo{
				{
					Path:      "file1.txt",
					Extension: ".txt",
					IsBinary:  false,
					Text:      "some text",
				},
			},
			expectedRemaining: []string{},
		},
		{
			name: "files and nested dirs",
			entries: []entry{
				{
					fileMetadata{
						Type:      TreeEntryFile,
						Name:      "file1.txt",
						Extension: ".txt",
					},
					entryObject{
						fileContents{
							IsBinary: false,
							Text:     "some text",
						},
					},
				},
				{
					fileMetadata{
						Type: TreeEntryDir,
						Name: "dir1",
						Path: "dir1",
					},
					entryObject{},
				},
			},
			expectedResults: []*types.FileInfo{
				{
					Path:      "file1.txt",
					Extension: ".txt",
					IsBinary:  false,
					Text:      "some text",
				},
			},
			expectedRemaining: []string{"dir1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree := &treeQuery{}
			tree.repository.Name = "some repo"
			tree.repository.Object.Tree.Entries = tt.entries

			results := make(chan *types.FileInfo, 100)
			remaining := make(chan string, 100)
			cancel := make(chan struct{}, 1)

			g := GitHub{
				logger: zerolog.Nop(),
			}

			g.parseTree(tree, results, remaining, cancel)

			close(results)
			close(remaining)

			actualResults := make([]*types.FileInfo, 0)
			for f := range results {
				actualResults = append(actualResults, f)
			}

			actualRemaining := make([]string, 0)
			for r := range remaining {
				actualRemaining = append(actualRemaining, r)
			}

			require.ElementsMatch(t, tt.expectedResults, actualResults)
			require.ElementsMatch(t, tt.expectedRemaining, actualRemaining)
		})
	}
}
