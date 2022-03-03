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

// TODO
// The `treeParser` interface would benefit from the introduction
// of generics in Go 1.18.
// A type list for the subtree (TreeLevelX) structs would
// mitigate the excrutiating duplication of the `parse` methods.
type treeParser interface {
	parse(fs *[]*FileInfo)
}

type treeQuery struct {
	Repository struct {
		Name   string
		Object struct {
			Tree TreeLevel3 `graphql:"... on Tree"`
		} `graphql:"object(expression: $commitishAndPath)"`
	} `graphql:"repository(owner: $owner, name: $repo)"`
}

func (t *treeQuery) parse(fs *[]*FileInfo) {
	t.Repository.Object.Tree.parse(fs)
}

var _ treeParser = (*treeQuery)(nil)

type TreeLevel3 struct {
	Entries []struct {
		FileMetadata
		Object struct {
			Tree         TreeLevel2 `graphql:"... on Tree"`
			FileContents `graphql:"... on Blob"`
		}
	}
}

func (t TreeLevel3) parse(fs *[]*FileInfo) {
	for _, e := range t.Entries {
		switch e.Type {
		case TreeEntryFile:
			f := &FileInfo{
				FileMetadata{
					Type:      e.Type,
					Name:      e.Name,
					Extension: e.Extension,
					Path:      e.Path,
				},
				FileContents{
					IsBinary: e.Object.IsBinary,
					Text:     e.Object.Text,
				},
			}
			*fs = append(*fs, f)
		case TreeEntryDir:
			e.Object.Tree.parse(fs)
		default:
			// TODO - log error
			return
		}
	}
}

var _ treeParser = (*TreeLevel3)(nil)

type TreeLevel2 struct {
	Entries []struct {
		FileMetadata
		Object struct {
			Tree         TreeLevel1 `graphql:"... on Tree"`
			FileContents `graphql:"... on Blob"`
		}
	}
}

func (t *TreeLevel2) parse(fs *[]*FileInfo) {
	for _, e := range t.Entries {
		switch e.Type {
		case TreeEntryFile:
			f := &FileInfo{
				FileMetadata{
					Type:      e.Type,
					Name:      e.Name,
					Extension: e.Extension,
					Path:      e.Path,
				},
				FileContents{
					IsBinary: e.Object.IsBinary,
					Text:     e.Object.Text,
				},
			}
			*fs = append(*fs, f)
		case TreeEntryDir:
			e.Object.Tree.parse(fs)
		default:
			// TODO - log error
			return
		}
	}
}

var _ treeParser = (*TreeLevel2)(nil)

type TreeLevel1 struct {
	Entries []struct {
		FileMetadata
		Object struct {
			Tree         TreeLevel0 `graphql:"... on Tree"`
			FileContents `graphql:"... on Blob"`
		}
	}
}

func (t *TreeLevel1) parse(fs *[]*FileInfo) {
	for _, e := range t.Entries {
		switch e.Type {
		case TreeEntryFile:
			f := &FileInfo{
				FileMetadata{
					Type:      e.Type,
					Name:      e.Name,
					Extension: e.Extension,
					Path:      e.Path,
				},
				FileContents{
					IsBinary: e.Object.IsBinary,
					Text:     e.Object.Text,
				},
			}
			*fs = append(*fs, f)
		case TreeEntryDir:
			e.Object.Tree.parse(fs)
		default:
			// TODO - log error
			return
		}
	}
}

var _ treeParser = (*TreeLevel1)(nil)

type TreeLevel0 struct {
	Entries []struct {
		FileMetadata
		Object struct {
			FileContents `graphql:"... on Blob"`
		}
	}
}

func (t *TreeLevel0) parse(fs *[]*FileInfo) {
	for _, e := range t.Entries {
		switch e.Type {
		case TreeEntryFile:
			f := &FileInfo{
				FileMetadata{
					Type:      e.Type,
					Name:      e.Name,
					Extension: e.Extension,
					Path:      e.Path,
				},
				FileContents{
					IsBinary: e.Object.IsBinary,
					Text:     e.Object.Text,
				},
			}
			*fs = append(*fs, f)
		case TreeEntryDir:
			// TODO - return entries for further queries
			continue
		default:
			// TODO - log error
			return
		}
	}
}

var _ treeParser = (*TreeLevel0)(nil)

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
