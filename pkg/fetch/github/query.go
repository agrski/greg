package fetch

const (
	TreeEntryDir  = "tree"
	TreeEntryFile = "blob"
)

type queryParams struct {
	repoOwner  string
	repoName   string
	commitish  string
	pathPrefix string
}

type Query struct {
	Repository struct {
		Name   string
		Object struct {
			FileInfoFragment `graphql:"... on Tree"`
		} `graphql:"object(expression: $commitishAndPath)"`
	} `graphql:"repository(owner: $owner, name: $repo)"`
}

type FileInfoFragment struct {
	Entries []FileInfo
}

type FileInfo struct {
	Name      string
	Type      string
	Extension string
}
