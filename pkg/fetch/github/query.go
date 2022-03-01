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
			TreeLevel3 `graphql:"... on Tree"`
		} `graphql:"object(expression: $commitishAndPath)"`
	} `graphql:"repository(owner: $owner, name: $repo)"`
}

type TreeLevel3 struct {
	Entries []struct {
		FileInfo
		Object struct {
			TreeLevel2 `graphql:"... on Tree"`
		}
	}
}

type TreeLevel2 struct {
	Entries []struct {
		FileInfo
		Object struct {
			TreeLevel1 `graphql:"... on Tree"`
		}
	}
}

type TreeLevel1 struct {
	Entries []struct {
		FileInfo
		Object struct {
			TreeLevel0 `graphql:"... on Tree"`
		}
	}
}

type TreeLevel0 struct {
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
