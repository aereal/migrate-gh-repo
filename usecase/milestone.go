package usecase

import (
	"context"
	"fmt"
	"log"

	"github.com/aereal/migrate-gh-repo/config"
	"github.com/aereal/migrate-gh-repo/domain"
	"github.com/google/go-github/github"
)

func (u *Usecase) buildMilestoneRequests(ctx context.Context, source, target *config.Repository) ([]request, error) {
	sourceMilestones, err := u.sourceService.SlurpMilestones(ctx, source.Owner, source.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch milestones from source repository: %w", err)
	}
	targetMilestones, err := u.targetService.SlurpMilestones(ctx, target.Owner, target.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch milestones from target repository: %w", err)
	}

	reqs := []request{}
	ops := domain.NewMilestoneOpsList(sourceMilestones, targetMilestones)
	for _, op := range ops {
		reqs = append(reqs, newMilestoneRequest(target, op))
	}
	return reqs, nil
}

type createMilestoneRequest struct {
	owner     string
	repo      string
	milestone *github.Milestone
}

func (r *createMilestoneRequest) Do(ctx context.Context, ghClient *github.Client) error {
	_, resp, err := ghClient.Issues.CreateMilestone(ctx, r.owner, r.repo, r.milestone)
	if err != nil {
		return err
	}
	log.Printf("create milestone owner=%s repo=%s statusCode=%d milestone=%s", r.owner, r.repo, resp.StatusCode, r.milestone)
	return nil
}

type updateMilestoneRequest struct {
	owner     string
	repo      string
	number    int
	milestone *github.Milestone
}

func (r *updateMilestoneRequest) Do(ctx context.Context, ghClient *github.Client) error {
	log.Printf("update milestone number=%d owner=%s repo=%s milestone=%s", r.number, r.owner, r.repo, r.milestone)
	_, _, err := ghClient.Issues.EditMilestone(ctx, r.owner, r.repo, r.number, r.milestone)
	if err != nil {
		return err
	}
	return nil
}

func newMilestoneRequest(repo *config.Repository, op *domain.MilestoneOp) request {
	switch op.Kind {
	case domain.OpCreate:
		return &createMilestoneRequest{owner: repo.Owner, repo: repo.Name, milestone: &github.Milestone{
			State:       op.Milestone.State,
			Title:       op.Milestone.Title,
			Description: op.Milestone.Description,
			DueOn:       op.Milestone.DueOn,
		}}
	case domain.OpUpdate:
		return &updateMilestoneRequest{owner: repo.Owner, repo: repo.Name, number: op.Milestone.GetNumber(), milestone: &github.Milestone{
			State:       op.Milestone.State,
			Title:       op.Milestone.Title,
			Description: op.Milestone.Description,
			DueOn:       op.Milestone.DueOn,
		}}
	default:
		return nil
	}
}
