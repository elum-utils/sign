package utils

import (
	"crypto/sha256"
	"sync"
)

var (
	BufCanonicalPool = sync.Pool{
		New: func() any {
			b := make([]byte, 0, 2048)
			return &b
		},
	}

	TmpBufPool = sync.Pool{
		New: func() any {
			b := make([]byte, 0, 256)
			return &b
		},
	}

	Base64BufPool = sync.Pool{
		New: func() any {
			b := make([]byte, 43)
			return &b
		},
	}

	Sha256SumBufPool = sync.Pool{
		New: func() any {
			b := make([]byte, sha256.Size)
			return &b
		},
	}

	KVPool = sync.Pool{
		New: func() any {
			s := make(KVSlice, 0, 56)
			return &s
		},
	}
)
