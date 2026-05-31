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

func (c *Client) ApprovePR(ctx context.Context, owner, repo string, number int, body string) error {
	return c.createPRReview(ctx, owner, repo, number, "APPROVE", body)
}

func (c *Client) RequestChangesPR(ctx context.Context, owner, repo string, number int, body string) error {
	return c.createPRReview(ctx, owner, repo, number, "REQUEST_CHANGES", body)
}

func (c *Client) MergePR(ctx context.Context, owner, repo string, number int) error {
	result, resp, err := c.client.PullRequests.Merge(ctx, owner, repo, number, "", &gh.PullRequestOptions{
		MergeMethod: "merge",
	})
	if err != nil {
		return fmt.Errorf("merge pull request: %s", formatGitHubError(resp, err))
	}
	if result != nil && result.GetMerged() {
		return nil
	}
	return fmt.Errorf("merge pull request: merge was not completed")
}

func (c *Client) CreatePRComment(ctx context.Context, owner, repo string, number int, body string) error {
	_, resp, err := c.client.Issues.CreateComment(ctx, owner, repo, number, &gh.IssueComment{
		Body: gh.Ptr(body),
	})
	if err != nil {
		return fmt.Errorf("create pull request comment: %s", formatGitHubError(resp, err))
	}
	return nil
}

func (c *Client) createPRReview(ctx context.Context, owner, repo string, number int, event, body string) error {
	req := &gh.PullRequestReviewRequest{
		Event: gh.Ptr(event),
	}
	if body != "" {
		req.Body = gh.Ptr(body)
	}

	_, resp, err := c.client.PullRequests.CreateReview(ctx, owner, repo, number, req)
	if err != nil {
		return fmt.Errorf("create pull request review (%s): %s", event, formatGitHubError(resp, err))
	}
	return nil
}

func formatGitHubError(resp *gh.Response, err error) string {
	if resp == nil {
		return err.Error()
	}
	switch resp.StatusCode {
	case 403:
		return fmt.Sprintf("permission denied (403): %v", err)
	case 405:
		return fmt.Sprintf("method not allowed (405), PR may already be merged or closed: %v", err)
	case 409:
		return fmt.Sprintf("conflict (409), branch may be out of date or protected: %v", err)
	case 422:
		return fmt.Sprintf("validation failed (422), branch protection or required checks may block this action: %v", err)
	default:
		return fmt.Sprintf("status %d: %v", resp.StatusCode, err)
	}
}
