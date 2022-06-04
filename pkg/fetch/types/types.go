package types

type Fetcher interface {
	Start() error
	Stop() error
	Next() (interface{}, bool) // TODO - formalise this return param for interop with matcher
}
