package github

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
	Repository `graphql:"repository(owner: $owner, name: $repo)"`
}

type Repository struct {
	Name   string
	Object repositoryObject `graphql:"object(expression: $commitishAndPath)"`
}

type repositoryObject struct {
	Tree tree `graphql:"... on Tree"`
}

type tree struct {
	Entries []entry
}

type entry struct {
	FileMetadata
	Object entryObject
}

type entryObject struct {
	FileContents `graphql:"... on Blob"`
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
