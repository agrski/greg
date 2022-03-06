package fetch

type TreeEntry string

const (
	TreeEntryDir  TreeEntry = "tree"
	TreeEntryFile TreeEntry = "blob"
)

type QueryParams struct {
	RepoOwner  string
	RepoName   string
	Commitish  string
	PathPrefix string
}

type branchRefQuery struct {
	Repository struct {
		DefaultBranchRef struct {
			Name string
		}
	} `graphql:"repository(owner: $owner, name: $repo)"`
}

type treeQuery struct {
	Repository struct {
		Name   string
		Object struct {
			Tree struct {
				Entries []struct {
					FileMetadata
					Object struct {
						FileContents `graphql:"... on Blob"`
					}
				}
			} `graphql:"... on Tree"`
		} `graphql:"object(expression: $commitishAndPath)"`
	} `graphql:"repository(owner: $owner, name: $repo)"`
}

type FileMetadata struct {
	Type      TreeEntry
	Name      string
	Extension string
	Path      string
}

type FileContents struct {
	IsBinary bool
	Text     string
}
