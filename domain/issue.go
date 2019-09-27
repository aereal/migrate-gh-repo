package domain

import (
	"fmt"

	"github.com/google/go-github/github"
)

type issue struct {
	*github.Issue
}

func (i *issue) hasMigrated() bool {
	if i == nil {
		return false
	}
	for _, label := range i.Labels {
		if label.GetName() == "migrated" {
			return true
		}
	}
	return false
}

func (i *issue) Key() *Key {
	if i == nil {
		return nil
	}
	return &Key{kind: "issue", repr: fmt.Sprintf("%d", i.GetNumber())}
}

func (i *issue) Eq(other Equalable) bool {
	if i == nil || other == nil {
		return false
	}
	if !i.Key().Eq(other.Key()) {
		return false
	}
	if otherIssue, ok := other.(*issue); ok {
		return i.GetTitle() == otherIssue.GetTitle()
	}
	return false
}

func NewIssueOpsList(sourceIssues, targetIssues []*github.Issue) IssueOpsList {
	if len(sourceIssues) == 0 && len(targetIssues) == 0 {
		return nil
	}

	kinds := map[string]OpKind{}
	for _, s := range sourceIssues {
		src := &issue{s}
		defaultKind := OpCreate
		kinds[src.Key().String()] = defaultKind
		for _, t := range targetIssues {
			target := &issue{t}
			if src.Key().Eq(target.Key()) {
				if target.hasMigrated() || src.Eq(target) { // completely equal
					kinds[src.Key().String()] = OpNothing
				} else {
					kinds[src.Key().String()] = OpUpdate
				}
			}
		}
	}

	ops := []*IssueOp{}
	for _, s := range sourceIssues {
		src := &issue{s}
		switch kinds[src.Key().String()] {
		case OpCreate:
			ops = append(ops, &IssueOp{
				Kind:  OpCreate,
				Issue: s,
			})
		case OpUpdate:
			ops = append(ops, &IssueOp{
				Kind:  OpUpdate,
				Issue: s,
			})
		default:
		}
	}
	return ops
}

type IssueOpsList []*IssueOp

func (il IssueOpsList) String() string {
	s := "["
	for _, op := range il {
		s += fmt.Sprintf("%s, ", op)
	}
	s += "]"
	return s
}

type IssueOp struct {
	Kind  OpKind
	Issue *github.Issue
}

func (op *IssueOp) String() string {
	return stringify(op.Kind, op.Issue)
}
