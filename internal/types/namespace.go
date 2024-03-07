package types

import (
	"fmt"
	"os"
	"sync"
)

type Namespace string

var (
	once sync.Once
	ns   Namespace
)

func (n Namespace) String() string {
	return string(n)
}

func SetNamespace() error {
	var onceerr error
	once.Do(func() {
		nsb, err := os.ReadFile(nsInferPath)
		if err != nil {
			onceerr = fmt.Errorf("unable to infer pod namespace from path %q: %w", nsInferPath, err)
		}
		ns = Namespace(string(nsb))
	})

	return onceerr
}
