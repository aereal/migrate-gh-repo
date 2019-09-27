package usecase

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aereal/migrate-gh-repo/config"
	"github.com/aereal/migrate-gh-repo/domain"
	"github.com/aereal/migrate-gh-repo/external"
	"github.com/google/go-github/github"
)

func New(sourceClient, targetClient *github.Client) (*Usecase, error) {
	if sourceClient == nil || targetClient == nil {
		return nil, fmt.Errorf("both of sourceClient and targetClient must be given")
	}
	return &Usecase{
		sourceClient: sourceClient,
		targetClient: targetClient,
	}, nil
}

type Usecase struct {
	sourceClient *github.Client
	targetClient *github.Client
}

type request interface {
	Do(ctx context.Context, ghClient *github.Client) error
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

type createLabelRequest struct {
	owner string
	repo  string
	label *github.Label
}

func (r *createLabelRequest) Do(ctx context.Context, ghClient *github.Client) error {
	_, resp, err := ghClient.Issues.CreateLabel(ctx, r.owner, r.repo, r.label)
	if err != nil {
		return err
	}
	log.Printf("create label owner=%s repo=%s statusCode=%d label=%s", r.owner, r.repo, resp.StatusCode, r.label)
	return nil
}

type updateLabelRequest struct {
	owner string
	repo  string
	name  string
	label *github.Label
}

func (r *updateLabelRequest) Do(ctx context.Context, ghClient *github.Client) error {
	log.Printf("update label name=%s owner=%s repo=%s label=%s", r.name, r.owner, r.repo, r.label)
	_, _, err := ghClient.Issues.EditLabel(ctx, r.owner, r.repo, r.name, r.label)
	if err != nil {
		return err
	}
	return nil
}

func newLabelRequest(repo *config.Repository, op *domain.LabelOp) request {
	switch op.Kind {
	case domain.OpCreate:
		return &createLabelRequest{owner: repo.Owner, repo: repo.Name, label: &github.Label{
			Name:        op.Label.Name,
			Color:       op.Label.Color,
			Description: op.Label.Description,
		}}
	case domain.OpUpdate:
		return &updateLabelRequest{owner: repo.Owner, repo: repo.Name, name: op.Label.GetName(), label: &github.Label{
			Color:       op.Label.Color,
			Description: op.Label.Description,
		}}
	default:
		return nil
	}
}

func (u *Usecase) Migrate(ctx context.Context, source, target *config.Repository) error {
	if source == nil || target == nil {
		return fmt.Errorf("Both of from/to repository must be given")
	}

	reqs, err := u.buildRequests(ctx, source, target)
	if err != nil {
		return err
	}
	interval := time.Second * 1
	tried := 0
	intervalCount := 10
	for _, r := range reqs {
		if err := r.Do(ctx, u.targetClient); err != nil {
			return err
		}
		tried++
		if tried >= intervalCount {
			time.Sleep(interval)
			tried = 0
		}
	}
	return nil
}

func (u *Usecase) buildRequests(ctx context.Context, source, target *config.Repository) ([]request, error) {
	reqs := []request{}

	milestoneReqs, err := u.buildMilestoneRequests(ctx, source, target)
	if err != nil {
		return nil, err
	}
	reqs = append(reqs, milestoneReqs...)

	labelReqs, err := u.buildLabelRequests(ctx, source, target)
	if err != nil {
		return nil, err
	}
	reqs = append(reqs, labelReqs...)

	issueReqs, err := u.buildIssueRequests(ctx, source, target)
	if err != nil {
		return nil, err
	}
	reqs = append(reqs, issueReqs...)
	return reqs, nil
}

func (u *Usecase) buildMilestoneRequests(ctx context.Context, source, target *config.Repository) ([]request, error) {
	sourceService, err := external.NewGitHubService(u.sourceClient)
	if err != nil {
		return nil, err
	}
	targetService, err := external.NewGitHubService(u.targetClient)
	if err != nil {
		return nil, err
	}

	sourceMilestones, err := sourceService.SlurpMilestones(ctx, source.Owner, source.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch milestones from source repository: %w", err)
	}
	log.Printf("source milestones = %#v", sourceMilestones)
	targetMilestones, err := targetService.SlurpMilestones(ctx, target.Owner, target.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch milestones from target repository: %w", err)
	}
	log.Printf("target milestones = %#v", targetMilestones)

	reqs := []request{}
	ops := domain.NewMilestoneOpsList(sourceMilestones, targetMilestones)
	for _, op := range ops {
		reqs = append(reqs, newMilestoneRequest(target, op))
	}
	return reqs, nil
}

func (u *Usecase) buildLabelRequests(ctx context.Context, source, target *config.Repository) ([]request, error) {
	sourceService, err := external.NewGitHubService(u.sourceClient)
	if err != nil {
		return nil, err
	}
	targetService, err := external.NewGitHubService(u.targetClient)
	if err != nil {
		return nil, err
	}

	sourceLabels, err := sourceService.SlurpLabels(ctx, source.Owner, source.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch labels from source repository: %w", err)
	}
	log.Printf("source labels = %#v", sourceLabels)
	targetLabels, err := targetService.SlurpLabels(ctx, target.Owner, target.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch labels from target repository: %w", err)
	}
	log.Printf("target labels = %#v", targetLabels)

	reqs := []request{}
	ops := domain.NewLabelOpsList(sourceLabels, targetLabels)
	for _, op := range ops {
		reqs = append(reqs, newLabelRequest(target, op))
	}
	return reqs, nil
}

func (u *Usecase) buildIssueRequests(ctx context.Context, source, target *config.Repository) ([]request, error) {
	reqs := []request{}
	return reqs, nil
}
