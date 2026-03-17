package domain

import "errors"

type GHRepo struct {
	Title     string `json:"title"`
	Desc      string `json:"description"`
	Creator   string `json:"creator"`
	StarsCnt  int32  `json:"stars_cnt"`
	ForksCnt  int32  `json:"forks_cnt"`
	IssuesCnt int32  `json:"issues_cnt"`
	Lang      string `json:"lang"`
	Size      int32  `json:"size"`
	CreatedAt string `json:"created_at"`
	License   string `json:"license"`
}

var ErrNotFound = errors.New("repository not found")
