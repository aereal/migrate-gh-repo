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

func (p *project) eq(other *project) bool {
	if p == nil || other == nil {
		return false
	}
	if !p.Key().Eq(other.Key()) {
		return false
	}
	return p.GetName() == other.GetName()
}

func NewProjectOpsList(sourceIssues, targetIssues []*github.Project) ProjectOpsList {
	if len(sourceIssues) == 0 && len(targetIssues) == 0 {
		return nil
	}

	kinds := opMapping{}
	mapping := map[string]*github.Project{}
	for _, s := range sourceIssues {
		src := &project{s}
		kinds.requestCreate(src)
		for _, t := range targetIssues {
			target := &project{t}
			// tell update (creating columns) if target project has same name
			if src.Key().Eq(target.Key()) {
				kinds.requestUpdate(src)
				mapping[src.Key().String()] = t
			}
		}
	}

	ops := []*ProjectOp{}
	for _, s := range sourceIssues {
		src := &project{s}
		switch kinds.get(src) {
		case OpCreate:
			ops = append(ops, &ProjectOp{
				Kind:    OpCreate,
				Project: s,
			})
		case OpUpdate:
			op := &ProjectOp{
				Kind:    OpUpdate,
				Project: s,
			}
			if target, ok := mapping[src.Key().String()]; ok {
				op.TargetProject = target
			}
			ops = append(ops, op)
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
	Kind          OpKind
	Project       *github.Project
	TargetProject *github.Project // maybe nil
}

func (op *ProjectOp) String() string {
	return stringify(op.Kind, op.Project)
}
