package usecase

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/aereal/migrate-gh-repo/config"
	"github.com/aereal/migrate-gh-repo/domain"
	"github.com/google/go-github/github"
)

func (u *Usecase) buildProjectRequests(ctx context.Context, source, target *config.Repository, issueMapping map[int]int64) ([]request, error) {
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
			columnReqs, err := u.buildProjectColumnRequests(ctx, op.Project, op.TargetProject, source, target, issueMapping)
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

func (u *Usecase) buildProjectColumnRequests(ctx context.Context, sourceProject, targetProject *github.Project, sourceRepo, targetRepo *config.Repository, issueMapping map[int]int64) ([]request, error) {
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
		case domain.OpUpdate:
			cardReqs, err := u.buildProjectCardRequests(ctx, op.ProjectColumn, op.TargetProjectColumn, sourceRepo, targetRepo, issueMapping)
			if err != nil {
				return nil, err
			}
			log.Printf("%d card reqs", len(cardReqs))

			reqs = append(reqs, cardReqs...)
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

type createProjectCardRequest struct {
	columnID int64
	opts     *github.ProjectCardOptions
}

func (r *createProjectCardRequest) Do(ctx context.Context, ghClient *github.Client) error {
	log.Printf("create project card (opts=%#v) on projectColumn.ID=%d", r.opts, r.columnID)
	_, _, err := ghClient.Projects.CreateProjectCard(ctx, r.columnID, r.opts)
	if err != nil {
		return err
	}
	return nil
}

func (u *Usecase) buildProjectCardRequests(ctx context.Context, sourceColumn, targetColumn *github.ProjectColumn, sourceRepo, targetRepo *config.Repository, issueMapping map[int]int64) ([]request, error) {
	sourceCards, err := u.sourceService.SlurpProjectCards(ctx, sourceColumn.GetID())
	if err != nil {
		return nil, fmt.Errorf("failed to fetch project cards on %s/%s columnId=%d: %w", sourceRepo.Owner, sourceRepo.Name, sourceColumn.GetID(), err)
	}
	targetCards, err := u.targetService.SlurpProjectCards(ctx, targetColumn.GetID())
	if err != nil {
		return nil, fmt.Errorf("failed to fetch project cards on %s/%s columnId=%d: %w", targetRepo.Owner, targetRepo.Name, targetColumn.GetID(), err)
	}
	log.Printf("%d source cards on column %q %d", len(sourceCards), sourceColumn.GetName(), sourceColumn.GetID())
	log.Printf("%d target cards on column %q %d", len(targetCards), targetColumn.GetName(), targetColumn.GetID())

	reqs := []request{}
	for _, op := range domain.NewProjectCardOpsList(sourceCards, targetCards, sourceColumn, targetColumn) {
		switch op.Kind {
		case domain.OpCreate:
			opts := &github.ProjectCardOptions{Note: op.ProjectCard.GetNote()}
			if opts.Note == "" {
				issues := "/issues/"
				contentURL := op.ProjectCard.GetContentURL() // e.g. https://api.github.com/repos/api-playground/projects-test/issues/3
				idx := strings.Index(contentURL, issues)
				if idx == -1 {
					log.Printf("! card (id=%d) invalid contentURL: %q", op.ProjectCard.GetID(), contentURL)
					continue
				}
				offset := idx + len(issues)
				repr := contentURL[offset:]
				issueNum, err := strconv.Atoi(repr)
				if err != nil {
					log.Printf("! card (id=%d) invalid contentURL: %q", op.ProjectCard.GetID(), contentURL)
					continue
				}
				id, ok := issueMapping[issueNum]
				if !ok {
					return nil, fmt.Errorf("no issue mapping found for number=%d", issueNum)
				}
				opts.ContentID = id
				opts.ContentType = "Issue"
			}
			req := &createProjectCardRequest{
				columnID: op.ProjectColumn.GetID(),
				opts:     opts,
			}
			reqs = append(reqs, req)
		default:
			// no-op
		}
	}

	return reqs, nil
}
