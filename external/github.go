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
