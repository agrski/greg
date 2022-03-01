package fetch

type TreeEntry string

const (
	TreeEntryDir  TreeEntry = "tree"
	TreeEntryFile TreeEntry = "blob"
)

type QueryParams struct {
	repoOwner  string
	repoName   string
	commitish  string
	pathPrefix string
}

type Query struct {
	Repository struct {
		Name   string
		Object struct {
			Tree struct {
				Entries []struct {
					FileInfo
					Object struct {
						Tree struct {
							Entries []struct {
								FileInfo
								Object struct {
									Leaf `graphql:"... on Tree"`
								}
							}
						} `graphql:"... on Tree"`
					}
				}
			} `graphql:"... on Tree"`
		} `graphql:"object(expression: $commitishAndPath)"`
	} `graphql:"repository(owner: $owner, name: $repo)"`
}

type Leaf struct {
	Entries []struct {
		FileInfo
	}
}

type FileInfo struct {
	Type      TreeEntry
	Name      string
	Extension string
	Path      string
}
