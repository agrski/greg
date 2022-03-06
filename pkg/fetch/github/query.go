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
			Tree treeLevel1 `graphql:"... on Tree"`
		} `graphql:"object(expression: $commitishAndPath)"`
	} `graphql:"repository(owner: $owner, name: $repo)"`
}

type treeLevel1 struct {
	Entries []struct {
		FileMetadata
		Object struct {
			FileContents `graphql:"... on Blob"`
		}
	}
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

// TODO - support both tree listing (to discover files to filter)
//	AND file querying, which is a simpler form with a path in `expression`
