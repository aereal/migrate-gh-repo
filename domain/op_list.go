package domain

import "github.com/google/go-github/github"

func eqMilestone(l *github.Milestone, r *github.Milestone) bool {
	if l == nil || r == nil {
		return false
	}
	return l.GetTitle() == r.GetTitle()
}
