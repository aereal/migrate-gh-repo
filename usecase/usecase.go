package usecase

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aereal/migrate-gh-repo/config"
	"github.com/aereal/migrate-gh-repo/domain"
	"github.com/aereal/migrate-gh-repo/external"
	"github.com/google/go-github/github"
)

func New(userResolver *domain.UserAliasResolver, sourceClient, targetClient *github.Client) (*Usecase, error) {
	if sourceClient == nil || targetClient == nil {
		return nil, fmt.Errorf("both of sourceClient and targetClient must be given")
	}
	sourceService, err := external.NewGitHubService(sourceClient)
	if err != nil {
		return nil, fmt.Errorf("failed to build GitHubService for source: %w", err)
	}
	targetService, err := external.NewGitHubService(targetClient)
	if err != nil {
		return nil, fmt.Errorf("failed to build GitHubService for target: %w", err)
	}

	return &Usecase{
		sourceClient:      sourceClient,
		targetClient:      targetClient,
		sourceService:     sourceService,
		targetService:     targetService,
		userAliasResolver: userResolver,
	}, nil
}

type Usecase struct {
	sourceClient      *github.Client
	sourceService     *external.GitHubService
	targetClient      *github.Client
	targetService     *external.GitHubService
	userAliasResolver *domain.UserAliasResolver
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

type createIssueRequest struct {
	owner    string
	repo     string
	issueReq *github.IssueRequest
}

func (r *createIssueRequest) Do(ctx context.Context, ghClient *github.Client) error {
	log.Printf(
		"create issue on %s/%s: title=%q body=%q labels=[%s] assignees=[%s] state=%q milestone.id=%d",
		r.owner, r.repo,
		r.issueReq.GetTitle(),
		r.issueReq.GetBody(),
		strings.Join(r.issueReq.GetLabels(), ", "),
		strings.Join(r.issueReq.GetAssignees(), ", "),
		r.issueReq.GetState(),
		r.issueReq.GetMilestone(),
	)
	_, _, err := ghClient.Issues.Create(ctx, r.owner, r.repo, r.issueReq)
	if err != nil {
		return err
	}
	return nil
}

type updateIssueRequest struct {
	owner       string
	repo        string
	issueNumber int
	issueReq    *github.IssueRequest
}

func (r *updateIssueRequest) Do(ctx context.Context, ghClinet *github.Client) error {
	log.Printf(
		"update issue on %s/%s#%d: title=%q body=%q labels=[%s] assignees=[%s] state=%q milestone.id=%d",
		r.owner, r.repo, r.issueNumber,
		r.issueReq.GetTitle(),
		r.issueReq.GetBody(),
		strings.Join(r.issueReq.GetLabels(), ", "),
		strings.Join(r.issueReq.GetAssignees(), ", "),
		r.issueReq.GetState(),
		r.issueReq.GetMilestone(),
	)
	if _, _, err := ghClinet.Issues.Edit(ctx, r.owner, r.repo, r.issueNumber, r.issueReq); err != nil {
		return err
	}
	return nil
}

type createIssueCommentRequest struct {
	owner       string
	repo        string
	issueNumber int
	body        string
}

func (r *createIssueCommentRequest) Do(ctx context.Context, ghClient *github.Client) error {
	issueComment := &github.IssueComment{Body: &r.body}
	log.Printf("create issue comment on %s/%s#%d issueComment=%s", r.owner, r.repo, r.issueNumber, issueComment)
	_, _, err := ghClient.Issues.CreateComment(ctx, r.owner, r.repo, r.issueNumber, issueComment)
	if err != nil {
		return err
	}
	return nil
}

func newIssueRequests(resolver *domain.UserAliasResolver, sourceRepo, targetRepo *config.Repository, op *domain.IssueOp) []request {
	switch op.Kind {
	case domain.OpCreate:
		body := fmt.Sprintf("This issue or P-R imported from %s in previous repository (%s/%s)", op.Issue.GetHTMLURL(), sourceRepo.Owner, sourceRepo.Name)
		assignees := []string{}
		for _, assignee := range op.Issue.Assignees {
			userOnTarget, _ := resolver.AssumeResolved(assignee.GetLogin())
			assignees = append(assignees, userOnTarget)
		}
		labels := []string{}
		for _, label := range op.Issue.Labels {
			labels = append(labels, label.GetName())
		}
		issueReq := &github.IssueRequest{
			Body:      &body,
			Assignees: &assignees,
			Labels:    &labels,
			Title:     op.Issue.Title,
			State:     op.Issue.State,
		}
		if op.Issue.Milestone != nil && op.Issue.Milestone.Number != nil {
			issueReq.Milestone = op.Issue.Milestone.Number
		}
		reqs := []request{&createIssueRequest{
			owner:    targetRepo.Owner,
			repo:     targetRepo.Name,
			issueReq: issueReq,
		}}
		if op.Issue.GetState() == "closed" {
			reqs = append(reqs, &updateIssueRequest{
				owner:       targetRepo.Owner,
				repo:        targetRepo.Name,
				issueNumber: op.Issue.GetNumber(),
				issueReq: &github.IssueRequest{
					State: op.Issue.State,
				},
			})
		}
		return reqs
	case domain.OpUpdate:
		body := fmt.Sprintf("This issue or P-R referenced as %s in previous repository (%s/%s)", op.Issue.GetHTMLURL(), sourceRepo.Owner, sourceRepo.Name)
		labels := []string{"migrated"}
		reqs := []request{
			&createIssueCommentRequest{
				owner:       targetRepo.Owner,
				repo:        targetRepo.Name,
				issueNumber: op.Issue.GetNumber(),
				body:        body,
			},
			&updateIssueRequest{
				owner:       targetRepo.Owner,
				repo:        targetRepo.Name,
				issueNumber: op.Issue.GetNumber(),
				issueReq: &github.IssueRequest{
					Labels: &labels,
				},
			},
		}
		return reqs
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

func (u *Usecase) buildIssueRequests(ctx context.Context, source, target *config.Repository) ([]request, error) {
	sourceIssues, err := u.sourceService.SlurpIssues(ctx, source.Owner, source.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch issues from source repository: %w", err)
	}

	targetIssues, err := u.targetService.SlurpIssues(ctx, target.Owner, target.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch issues from target repository: %w", err)
	}

	reqs := []request{}
	ops := domain.NewIssueOpsList(sourceIssues, targetIssues)
	for _, op := range ops {
		reqs = append(reqs, newIssueRequests(u.userAliasResolver, source, target, op)...)
	}
	return reqs, nil
}
