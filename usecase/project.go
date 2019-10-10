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
		switch op.Kind {
		case domain.OpCreate:
			reqs = append(reqs, &createProjectRequest{
				owner: target.Owner,
				repo:  target.Name,
				opts: &github.ProjectOptions{
					Name: op.Project.GetName(),
					Body: op.Project.GetBody(),
				},
			})
		case domain.OpUpdate:
			columnReqs, err := u.buildProjectColumnRequests(ctx, op.Project, op.TargetProject, source, target)
			if err != nil {
				return nil, err
			}
			log.Printf("%d project column requests", len(columnReqs))

			reqs = append(reqs, columnReqs...)
		default:
			// no-op
		}
	}

	// TODO: accumulate project column, project column card

	return reqs, nil
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

func (u *Usecase) buildProjectColumnRequests(ctx context.Context, sourceProject, targetProject *github.Project, sourceRepo, targetRepo *config.Repository) ([]request, error) {
	sourceProjectColumns, err := u.sourceService.SlurpProjectColumns(ctx, sourceProject.GetID())
	if err != nil {
		return nil, fmt.Errorf("failed to fetch project columns on %s/%s id=%d: %w", sourceRepo.Owner, sourceRepo.Name, sourceProject.GetID(), err)
	}
	targetProjectColumns, err := u.targetService.SlurpProjectColumns(ctx, targetProject.GetID())
	if err != nil {
		return nil, fmt.Errorf("failed to fetch project columns on %s/%s id=%d: %w", targetRepo.Owner, targetRepo.Name, targetProject.GetID(), err)
	}

	reqs := []request{}
	ops := domain.NewProjectColumnOpsList(sourceProjectColumns, targetProjectColumns, sourceProject, targetProject)
	for _, op := range ops {
		switch op.Kind {
		case domain.OpCreate:
			req := &createProjectColumnRequest{
				projectID: op.Project.GetID(),
				opts: &github.ProjectColumnOptions{
					Name: op.ProjectColumn.GetName(),
				},
			}
			reqs = append(reqs, req)
		default:
			// no-op
		}
	}

	return reqs, nil
}

type createProjectColumnRequest struct {
	projectID int64
	opts      *github.ProjectColumnOptions
}

func (r *createProjectColumnRequest) Do(ctx context.Context, ghClient *github.Client) error {
	log.Printf("create project column (%q) on project.ID=%d", r.opts.Name, r.projectID)
	_, _, err := ghClient.Projects.CreateProjectColumn(ctx, r.projectID, r.opts)
	if err != nil {
		return err
	}
	return nil
}
