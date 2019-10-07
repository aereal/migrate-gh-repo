package usecase

import (
	"context"
	"log"

	"github.com/google/go-github/github"
)

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
