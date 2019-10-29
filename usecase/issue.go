package usecase

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/aereal/migrate-gh-repo/config"
	"github.com/aereal/migrate-gh-repo/domain"
	"github.com/google/go-github/github"
)

func (u *Usecase) buildIssueRequests(ctx context.Context, source, target *config.Repository) ([]request, error) {
	sourceIssues, err := u.sourceService.SlurpIssues(ctx, source.Owner, source.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch issues from source repository: %w", err)
	}

	targetIssues, err := u.targetService.SlurpIssues(ctx, target.Owner, target.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch issues from target repository: %w", err)
	}
	for _, issue := range targetIssues {
		u.issueNumberIDMapping[issue.GetNumber()] = issue.GetID()
	}

	reqs := []request{}
	ops := domain.NewIssueOpsList(sourceIssues, targetIssues)
	for _, op := range ops {
		reqs = append(reqs, newIssueRequests(u.userAliasResolver, source, target, u.skipUsers, op)...)
	}
	return reqs, nil
}

func newIssueRequests(resolver *domain.UserAliasResolver, sourceRepo, targetRepo *config.Repository, skipUsers []string, op *domain.IssueOp) []request {
	switch op.Kind {
	case domain.OpCreate:
		body := fmt.Sprintf("This issue or P-R imported from %s in previous repository (%s/%s)", op.Issue.GetHTMLURL(), sourceRepo.Owner, sourceRepo.Name)
		assignees := []string{}
		for _, u := range op.Issue.Assignees {
			if contains(skipUsers, u.GetLogin()) {
				continue
			}
			userOnTarget, _ := resolver.AssumeResolved(u.GetLogin())
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
		log.Printf("update issue")
		body := fmt.Sprintf("This issue or P-R referenced as %s in previous repository (%s/%s)", op.Issue.GetHTMLURL(), sourceRepo.Owner, sourceRepo.Name)
		labels := []string{"migrated"}
		assignees := []string{}
		for _, u := range op.Issue.Assignees {
			if contains(skipUsers, u.GetLogin()) {
				continue
			}
			userOnTarget, _ := resolver.AssumeResolved(u.GetLogin())
			assignees = append(assignees, userOnTarget)
		}
		for _, l := range op.Issue.Labels {
			labels = append(labels, l.GetName())
		}
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
					Labels:    &labels,
					Assignees: &assignees,
				},
			},
		}
		return reqs
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
