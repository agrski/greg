package types

import (
	common "github.com/agrski/greg/pkg/types"
)

type HostName string
type OrganisationName string
type RepositoryName string

type Location struct {
	Host         HostName
	Organisation OrganisationName
	Repository   RepositoryName
}

type Fetcher interface {
	Start() error
	Stop() error
	Next() (*common.FileInfo, bool)
}
