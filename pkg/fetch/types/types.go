package types

import (
	common "github.com/agrski/greg/pkg/types"
)

type Fetcher interface {
	Start() error
	Stop() error
	Next() (*common.FileInfo, bool)
}
