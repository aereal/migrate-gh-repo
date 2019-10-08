package usecase

import (
	"context"
	"fmt"
	"log"

	"github.com/aereal/migrate-gh-repo/config"
	"github.com/aereal/migrate-gh-repo/domain"
	"github.com/google/go-github/github"
)

func (u *Usecase) buildProjectRequests(ctx context.Context, source, target *config.Repository) ([]request, error) {
	sourceProjects, err := u.sourceService.SlurpProjects(ctx, source.Owner, source.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch projects from source repository: %w", err)
	}
	targetProjects, err := u.targetService.SlurpProjects(ctx, target.Owner, target.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch projects from target repository: %w", err)
	}

	reqs := []request{}
	ops := domain.NewProjectOpsList(sourceProjects, targetProjects)
	for _, op := range ops {
		reqs = append(reqs, newProjectRequest(source, target, op)...)
	}

	// TODO: accumulate project column, project column card

	return reqs, nil
}

func newProjectRequest(sourceRepo, targetRepo *config.Repository, op *domain.ProjectOp) []request {
	switch op.Kind {
	case domain.OpCreate:
		reqs := []request{}
		reqs = append(reqs, &createProjectRequest{
			owner: targetRepo.Owner,
			repo:  targetRepo.Name,
			opts: &github.ProjectOptions{
				Name: op.Project.GetName(),
				Body: op.Project.GetBody(),
			},
		})
		return reqs
	default:
		return nil
	}
}

type createProjectRequest struct {
	owner string
	repo  string
	opts  *github.ProjectOptions
}

func (r *createProjectRequest) Do(ctx context.Context, ghClient *github.Client) error {
	log.Printf("create project (%q) on %s/%s", r.opts.Name, r.owner, r.repo)
	_, _, err := ghClient.Repositories.CreateProject(ctx, r.owner, r.repo, r.opts)
	if err != nil {
		return err
	}
	return nil
}
