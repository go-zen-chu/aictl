package git

import (
	"fmt"

	"github.com/go-git/go-git/v5"
)

type GitHandler interface {
	ChangedFiles() ([]string, error)
}

type gitHandler struct {
	repo *git.Repository
}

func NewGitHandler(repoRootPath string) (GitHandler, error) {
	repo, err := git.PlainOpen(repoRootPath)
	if err != nil {
		return nil, fmt.Errorf("open git repository: %w", err)
	}
	return &gitHandler{
		repo: repo,
	}, nil
}

func (g *gitHandler) ChangedFiles() ([]string, error) {
	wt, err := g.repo.Worktree()
	if err != nil {
		return nil, fmt.Errorf("get worktree: %w", err)
	}
	status, err := wt.Status()
	if err != nil {
		return nil, fmt.Errorf("get status: %w", err)
	}
	files := make([]string, 0, len(status))
	for file := range status {
		files = append(files, file)
	}
	return files, nil
}
