package github

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	gh "github.com/google/go-github/v69/github"
	"github.com/ZGjie20/PR-check-r/ai-pr-review/internal/model"
	"golang.org/x/oauth2"
)

var prURLPattern = regexp.MustCompile(`(?i)github\.com/([^/]+)/([^/]+)/pull/(\d+)`)

type Client struct {
	client *gh.Client
}

func NewClient(token string) *Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	httpClient := oauth2.NewClient(context.Background(), ts)
	return &Client{client: gh.NewClient(httpClient)}
}

func ParsePRURL(rawURL string) (owner, repo string, number int, err error) {
	rawURL = strings.TrimSpace(rawURL)
	matches := prURLPattern.FindStringSubmatch(rawURL)
	if len(matches) != 4 {
		return "", "", 0, fmt.Errorf("invalid GitHub PR URL: %s", rawURL)
	}

	num, err := strconv.Atoi(matches[3])
	if err != nil {
		return "", "", 0, fmt.Errorf("invalid PR number in URL: %w", err)
	}

	return matches[1], strings.TrimSuffix(matches[2], "/"), num, nil
}

type FetchResult struct {
	PR       *model.PullRequest
	Commits  []string
}

func (c *Client) FetchPR(ctx context.Context, owner, repo string, number int) (*FetchResult, error) {
	pr, _, err := c.client.PullRequests.Get(ctx, owner, repo, number)
	if err != nil {
		return nil, fmt.Errorf("get pull request: %w", err)
	}

	files, err := c.listChangedFiles(ctx, owner, repo, number)
	if err != nil {
		return nil, err
	}

	diff, err := c.fetchRawDiff(ctx, owner, repo, number)
	if err != nil {
		return nil, err
	}

	commits, err := c.listCommitMessages(ctx, owner, repo, number)
	if err != nil {
		return nil, err
	}

	author := ""
	if pr.User != nil && pr.User.Login != nil {
		author = *pr.User.Login
	}

	title := ""
	if pr.Title != nil {
		title = *pr.Title
	}

	return &FetchResult{
		PR: &model.PullRequest{
			Repo:   fmt.Sprintf("%s/%s", owner, repo),
			Number: number,
			Title:  title,
			Author: author,
			Files:  files,
			Diff:   diff,
		},
		Commits: commits,
	}, nil
}

func (c *Client) listChangedFiles(ctx context.Context, owner, repo string, number int) ([]string, error) {
	var paths []string
	opts := &gh.ListOptions{PerPage: 100}

	for {
		files, resp, err := c.client.PullRequests.ListFiles(ctx, owner, repo, number, opts)
		if err != nil {
			return nil, fmt.Errorf("list pull request files: %w", err)
		}

		for _, f := range files {
			if f.Filename != nil {
				paths = append(paths, *f.Filename)
			}
		}

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return paths, nil
}

func (c *Client) fetchRawDiff(ctx context.Context, owner, repo string, number int) (string, error) {
	raw, _, err := c.client.PullRequests.GetRaw(ctx, owner, repo, number, gh.RawOptions{
		Type: gh.Diff,
	})
	if err != nil {
		return "", fmt.Errorf("get pull request diff: %w", err)
	}
	return string(raw), nil
}

func (c *Client) listCommitMessages(ctx context.Context, owner, repo string, number int) ([]string, error) {
	var messages []string
	opts := &gh.ListOptions{PerPage: 100}

	for {
		commits, resp, err := c.client.PullRequests.ListCommits(ctx, owner, repo, number, opts)
		if err != nil {
			return nil, fmt.Errorf("list pull request commits: %w", err)
		}

		for _, commit := range commits {
			if commit.Commit != nil && commit.Commit.Message != nil {
				msg := strings.SplitN(*commit.Commit.Message, "\n", 2)[0]
				messages = append(messages, msg)
			}
		}

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return messages, nil
}
