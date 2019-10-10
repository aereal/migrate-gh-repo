package usecase

import (
	"context"
	"fmt"
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
		sourceClient:         sourceClient,
		targetClient:         targetClient,
		sourceService:        sourceService,
		targetService:        targetService,
		userAliasResolver:    userResolver,
		issueNumberIDMapping: map[int]int64{},
	}, nil
}

type Usecase struct {
	sourceClient         *github.Client
	sourceService        *external.GitHubService
	targetClient         *github.Client
	targetService        *external.GitHubService
	userAliasResolver    *domain.UserAliasResolver
	issueNumberIDMapping map[int]int64 // number -> id
}

type request interface {
	Do(ctx context.Context, ghClient *github.Client) error
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

	projectReqs, err := u.buildProjectRequests(ctx, source, target, u.issueNumberIDMapping)
	if err != nil {
		return nil, err
	}
	reqs = append(reqs, projectReqs...)

	return reqs, nil
}
