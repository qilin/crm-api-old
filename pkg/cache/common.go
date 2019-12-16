package cache

import "github.com/pkg/errors"

const (
	Prefix       = "pkg.cache"
	UnmarshalKey = "cache"
)

var ErrNotSupported = errors.WithMessage(errors.New("store not supported"), Prefix)
