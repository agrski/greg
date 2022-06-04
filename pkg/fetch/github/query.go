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
	repository `graphql:"repository(owner: $owner, name: $repo)"`
}

type repository struct {
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
	fileMetadata
	Object entryObject
}

type entryObject struct {
	fileContents `graphql:"... on Blob"`
}

type fileMetadata struct {
	Type      TreeEntry
	Name      string
	Extension string
	Path      string
}

type fileContents struct {
	IsBinary bool
	Text     string
}
