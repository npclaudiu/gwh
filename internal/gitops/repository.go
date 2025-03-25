package gitops

import (
	"github.com/go-git/go-git/v5"
)

func OpenRepository(path string) (*git.Repository, error) {
	return git.PlainOpen(path)
}
