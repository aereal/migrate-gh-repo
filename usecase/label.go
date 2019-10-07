package usecase

import (
	"context"
	"fmt"
	"log"

	"github.com/aereal/migrate-gh-repo/config"
	"github.com/aereal/migrate-gh-repo/domain"
	"github.com/google/go-github/github"
)

func (u *Usecase) buildLabelRequests(ctx context.Context, source, target *config.Repository) ([]request, error) {
	sourceLabels, err := u.sourceService.SlurpLabels(ctx, source.Owner, source.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch labels from source repository: %w", err)
	}
	targetLabels, err := u.targetService.SlurpLabels(ctx, target.Owner, target.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch labels from target repository: %w", err)
	}

	reqs := []request{}
	ops := domain.NewLabelOpsList(sourceLabels, targetLabels)
	for _, op := range ops {
		reqs = append(reqs, newLabelRequest(target, op))
	}
	return reqs, nil
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
