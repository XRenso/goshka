package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/XRenso/barerepo/collector/internal/domain"
)

// Client is the adapter that talks to the GitHub REST API
type Client struct {
	httpClient *http.Client
}

const requestTimeout = 10 * time.Second

func NewClient() *Client {
	return &Client{httpClient: &http.Client{Timeout: requestTimeout}}
}

// ghRepo mirrors the GitHub API response
type ghRepo struct {
	Title   string `json:"name"`
	Desc    string `json:"description"`
	Creator struct {
		Login string `json:"login"`
	} `json:"owner"`
	StarsCnt  int    `json:"stargazers_count"`
	ForksCnt  int    `json:"forks_count"`
	IssuesCnt int    `json:"open_issues_count"`
	Lang      string `json:"language"`
	Size      int    `json:"size"`
	CreatedAt string `json:"created_at"`
	License   *struct {
		Name string `json:"name"`
	} `json:"license"`
}

func (c *Client) FetchRepo(ctx context.Context, owner, repo string) (*domain.GHRepo, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		return nil, fmt.Errorf("%w: %s/%s", domain.ErrNotFound, owner, repo)
	case http.StatusForbidden:
		return nil, fmt.Errorf("%w: set GITHUB_TOKEN to increase rate limit", domain.ErrForbidden)
	default:
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	var gh ghRepo
	if err := json.Unmarshal(body, &gh); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	license := ""
	if gh.License != nil {
		license = gh.License.Name
	}

	return &domain.GHRepo{
		Title:     gh.Title,
		Desc:      gh.Desc,
		Creator:   gh.Creator.Login,
		StarsCnt:  gh.StarsCnt,
		ForksCnt:  gh.ForksCnt,
		IssuesCnt: gh.IssuesCnt,
		Lang:      gh.Lang,
		Size:      gh.Size,
		CreatedAt: gh.CreatedAt,
		License:   license,
	}, nil
}
