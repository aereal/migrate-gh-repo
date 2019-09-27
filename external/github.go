package external

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/go-github/github"
)

func NewGitHubService(client *github.Client) (*GitHubService, error) {
	if client == nil {
		return nil, errors.New("client (*github.Client) must be given")
	}
	return &GitHubService{client: client}, nil
}

type GitHubService struct {
	client *github.Client
}

func (s *GitHubService) SlurpMilestones(ctx context.Context, owner, repo string) ([]*github.Milestone, error) {
	opts := &github.MilestoneListOptions{State: "all"}
	milestones := []*github.Milestone{}
	for {
		ms, resp, err := s.client.Issues.ListMilestones(ctx, owner, repo, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to list milestones: %w", err)
		}
		milestones = append(milestones, ms...)
		opts.Page = resp.NextPage
		if resp.NextPage == 0 {
			break
		}
	}
	return milestones, nil
}

func (s *GitHubService) SlurpLabels(ctx context.Context, owner, repo string) ([]*github.Label, error) {
	opts := &github.ListOptions{}
	labels := []*github.Label{}
	for {
		ls, resp, err := s.client.Issues.ListLabels(ctx, owner, repo, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to list labels: %w", err)
		}
		labels = append(labels, ls...)
		opts.Page = resp.NextPage
		if resp.NextPage == 0 {
			break
		}
	}
	return labels, nil
}

func (s *GitHubService) SlurpIssues(ctx context.Context, owner, repo string) ([]*github.Issue, error) {
	opts := &github.IssueListByRepoOptions{State: "all", Direction: "asc"}
	issues := []*github.Issue{}
	for {
		is, resp, err := s.client.Issues.ListByRepo(ctx, owner, repo, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to list issues: %w", err)
		}
		issues = append(issues, is...)
		opts.Page = resp.NextPage
		if resp.NextPage == 0 {
			break
		}
	}
	return issues, nil
}

func (s *GitHubService) SlurpPullRequests(ctx context.Context, owner, repo string) ([]*github.PullRequest, error) {
	opts := &github.PullRequestListOptions{State: "all"}
	pullRequests := []*github.PullRequest{}
	for {
		prs, resp, err := s.client.PullRequests.List(ctx, owner, repo, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to list pull requests: %w", err)
		}
		pullRequests = append(pullRequests, prs...)
		opts.Page = resp.NextPage
		if resp.NextPage == 0 {
			break
		}
	}
	return pullRequests, nil
}
