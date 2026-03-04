package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type GHRepo struct {
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

func parseRepoInput(input string) (string, error) {
	input = strings.TrimSpace(input)

	re := regexp.MustCompile(`^(?:https?://)?(?:github\.com/)?([a-zA-Z0-9_.-]+/[a-zA-Z0-9_.-]+)(?:\.git)?(?:/.*)?$`)
	matches := re.FindStringSubmatch(input)
	if matches == nil {
		return "", fmt.Errorf("invalid repository format: %q\n\nSupported formats:\n  https://github.com/creator/repo\n  github.com/creator/repo\n  creator/repo", input)
	}
	return matches[1], nil
}

func fetchRepo(creatorRepo string) (*GHRepo, error) {
	url := "https://api.github.com/repos/" + creatorRepo

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("repository %q not found (404)", creatorRepo)
	}
	if resp.StatusCode == 403 {
		return nil, fmt.Errorf("access forbidden (403) — you may have hit the rate limit; set GITHUB_TOKEN env variable to increase it")
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	var repo GHRepo
	if err := json.Unmarshal(body, &repo); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &repo, nil
}

func (ghr *GHRepo) String() string {
	sep := strings.Repeat("─", 40)

	license := "—"
	if ghr.License != nil {
		license = ghr.License.Name
	}

	return fmt.Sprintf(
		"\n  %s/%s\n  %s\n\n %s\n  ★ %-6d  ⑂ %-6d  %s\n  issues %-4d  size %d KB  %s\n  created %s\n %s\n\n",
		ghr.Creator.Login, ghr.Title,
		strOrDash(ghr.Desc),
		sep,
		ghr.StarsCnt, ghr.ForksCnt, strOrDash(ghr.Lang),
		ghr.IssuesCnt, ghr.Size, license,
		formatTime(ghr.CreatedAt)[:10],
		sep,
	)
}

func strOrDash(s string) string {
	if s == "" {
		return "—"
	}
	return s
}

func formatTime(t string) string {
	t = strings.ReplaceAll(t, "T", " ")
	t = strings.TrimSuffix(t, "Z")
	return t
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: gitcli <repo>\n\nExamples:\n gitcli https://github.com/creator/repo\n gitcli github.com/creator/repo\n gitcli creator/repo")
		os.Exit(1)
	}

	repoLink, err := parseRepoInput(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	ghr, err := fetchRepo(repoLink)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	fmt.Print(ghr)
}
