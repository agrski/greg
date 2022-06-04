package types

type FileExtension string

type FileInfo struct {
	Path      string
	Extension FileExtension
	IsBinary  bool
	Text      string
}
