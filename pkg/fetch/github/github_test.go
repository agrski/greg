package github

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseTree(t *testing.T) {
	type treeEntry struct {
		FileMetadata
		Object struct {
			FileContents "graphql:\"... on Blob\""
		}
	}

	type test struct {
		name    string
		entries []struct {
			FileMetadata
			Object struct {
				FileContents "graphql:\"... on Blob\""
			}
		}
		expectedResults   []*FileInfo
		expectedRemaining []string
	}

	tests := []test{
		{
			name: "empty root dir",
			entries: []struct {
				FileMetadata
				Object struct {
					FileContents "graphql:\"... on Blob\""
				}
			}{},
			expectedResults:   []*FileInfo{},
			expectedRemaining: []string{},
		},
		{
			name: "one empty directory",
			entries: []struct {
				FileMetadata
				Object struct {
					FileContents "graphql:\"... on Blob\""
				}
			}{
				{
					FileMetadata{
						Type: TreeEntryDir,
						Name: "dir1",
						Path: "dir1",
					},
					struct {
						FileContents "graphql:\"... on Blob\""
					}{
						FileContents{},
					},
				},
			},
			expectedResults:   []*FileInfo{},
			expectedRemaining: []string{"dir1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree := &treeQuery{}
			tree.Repository.Name = "some repo"
			tree.Repository.Object.Tree.Entries = tt.entries

			results := make(chan *FileInfo, 100)
			remaining := make(chan string, 100)
			cancel := make(chan struct{}, 1)

			parseTree(tree, results, remaining, cancel)

			close(results)
			close(remaining)

			actualResults := make([]*FileInfo, 0)
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
