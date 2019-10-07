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
	opts := &github.MilestoneListOptions{State: "all", ListOptions: github.ListOptions{PerPage: 100}}
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
	opts := &github.ListOptions{PerPage: 100}
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
	opts := &github.IssueListByRepoOptions{State: "all", Direction: "asc", ListOptions: github.ListOptions{PerPage: 100}}
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

func (s *GitHubService) SlurpIssueComments(ctx context.Context, owner, repo string, issueNumber int) ([]*github.IssueComment, error) {
	opts := &github.IssueListCommentsOptions{ListOptions: github.ListOptions{PerPage: 100}}
	issueComments := []*github.IssueComment{}
	for {
		comments, resp, err := s.client.Issues.ListComments(ctx, owner, repo, issueNumber, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to list issue comments: %w", err)
		}
		issueComments = append(issueComments, comments...)
		opts.Page = resp.NextPage
		if resp.NextPage == 0 {
			break
		}
	}
	return issueComments, nil
}

func (s *GitHubService) SlurpProjects(ctx context.Context, owner, repo string) ([]*github.Project, error) {
	opts := &github.ProjectListOptions{State: "all", ListOptions: github.ListOptions{PerPage: 100}}
	projects := []*github.Project{}
	for {
		pjs, resp, err := s.client.Repositories.ListProjects(ctx, owner, repo, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to list projects: %w", err)
		}
		projects = append(projects, pjs...)
		opts.Page = resp.NextPage
		if resp.NextPage == 0 {
			break
		}
	}
	return projects, nil
}

func (s *GitHubService) SlurpProjectColumns(ctx context.Context, projectID int64) ([]*github.ProjectColumn, error) {
	opts := &github.ListOptions{PerPage: 100}
	columns := []*github.ProjectColumn{}
	for {
		cols, resp, err := s.client.Projects.ListProjectColumns(ctx, projectID, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to list project columns: %w", err)
		}
		columns = append(columns, cols...)
		opts.Page = resp.NextPage
		if resp.NextPage == 0 {
			break
		}
	}
	return columns, nil
}
