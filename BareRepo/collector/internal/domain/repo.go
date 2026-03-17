package domain

import "errors"

type GHRepo struct {
	Title     string
	Desc      string
	Creator   string
	StarsCnt  int
	ForksCnt  int
	IssuesCnt int
	Lang      string
	Size      int
	CreatedAt string
	License   string
}

var (
	ErrNotFound  = errors.New("repository not found")
	ErrForbidden = errors.New("access forbidden")
)
