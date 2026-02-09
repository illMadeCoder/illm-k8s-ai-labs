package github

import (
	"context"
	"encoding/json"
	"fmt"

	"golang.org/x/oauth2"

	gh "github.com/google/go-github/v68/github"
)

// Client wraps the GitHub Contents API for committing experiment results.
type Client struct {
	client *gh.Client
	owner  string
	repo   string
	branch string
	path   string // e.g. "site/data"
}

// NewClient creates a GitHub client for committing experiment results.
// owner and repo are parsed from the "owner/repo" format.
func NewClient(token, owner, repo, branch, path string) *Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(context.Background(), ts)
	return &Client{
		client: gh.NewClient(tc),
		owner:  owner,
		repo:   repo,
		branch: branch,
		path:   path,
	}
}

// RepoPath returns "owner/repo" for logging.
func (c *Client) RepoPath() string {
	return c.owner + "/" + c.repo
}

// CommitResult commits an experiment summary JSON to the configured repo path.
// It creates or updates site/data/{expName}.json with indented JSON.
func (c *Client) CommitResult(ctx context.Context, expName string, summary any) error {
	body, err := json.MarshalIndent(summary, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal summary JSON: %w", err)
	}

	filePath := c.path + "/" + expName + ".json"
	commitMsg := fmt.Sprintf("data: Add %s experiment results", expName)

	// Check if file already exists (need SHA for updates)
	existing, _, resp, err := c.client.Repositories.GetContents(ctx, c.owner, c.repo, filePath, &gh.RepositoryContentGetOptions{
		Ref: c.branch,
	})
	if err != nil && (resp == nil || resp.StatusCode != 404) {
		return fmt.Errorf("check existing file %s: %w", filePath, err)
	}

	opts := &gh.RepositoryContentFileOptions{
		Message: &commitMsg,
		Content: body,
		Branch:  &c.branch,
	}

	// If file exists, include SHA for update
	if existing != nil {
		sha := existing.GetSHA()
		opts.SHA = &sha
		commitMsg = fmt.Sprintf("data: Update %s experiment results", expName)
		opts.Message = &commitMsg
	}

	_, _, err = c.client.Repositories.CreateFile(ctx, c.owner, c.repo, filePath, opts)
	if err != nil {
		return fmt.Errorf("commit %s: %w", filePath, err)
	}

	return nil
}
