package domain

import (
	"fmt"

	"github.com/google/go-github/github"
)

type project struct {
	*github.Project
}

func (p *project) Key() *Key {
	if p == nil {
		return nil
	}
	return &Key{
		kind: "project",
		repr: p.GetName(),
	}
}

func (p *project) Eq(other Equalable) bool {
	if p == nil || other == nil {
		return false
	}
	if !p.Key().Eq(other.Key()) {
		return false
	}
	if otherProject, ok := other.(*project); ok {
		return p.GetName() == otherProject.GetName() // TODO
	}
	return false
}

func NewProjectOpsList(sourceIssues, targetIssues []*github.Project) ProjectOpsList {
	if len(sourceIssues) == 0 && len(targetIssues) == 0 {
		return nil
	}

	kinds := map[string]OpKind{}
	for _, s := range sourceIssues {
		src := &project{s}
		defaultKind := OpCreate
		kinds[src.Key().String()] = defaultKind
		for _, t := range targetIssues {
			target := &project{t}
			// project does not support update
			if src.Eq(target) {
				kinds[src.Key().String()] = OpNothing
			}
		}
	}

	ops := []*ProjectOp{}
	for _, s := range sourceIssues {
		src := &project{s}
		switch kinds[src.Key().String()] {
		case OpCreate:
			ops = append(ops, &ProjectOp{
				Kind:    OpCreate,
				Project: s,
			})
		default:
		}
	}
	return ops
}

type ProjectOpsList []*ProjectOp

func (l ProjectOpsList) String() string {
	s := "["
	for _, op := range l {
		s += fmt.Sprintf("%s, ", op)
	}
	s += "]"
	return s
}

type ProjectOp struct {
	Kind    OpKind
	Project *github.Project
}

func (op *ProjectOp) String() string {
	return stringify(op.Kind, op.Project)
}
