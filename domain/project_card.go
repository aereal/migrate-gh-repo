package domain

import (
	"fmt"

	"github.com/google/go-github/github"
)

type projectCard struct {
	*github.ProjectCard
}

func (c *projectCard) Key() *Key {
	if c == nil {
		return nil
	}
	return &Key{
		kind: "project_card",
		repr: fmt.Sprintf("id=%d", c.GetID()),
	}
}

func (c *projectCard) eq(other *projectCard) bool {
	if c == nil || other == nil {
		return false
	}
	if !c.Key().Eq(other.Key()) {
		return false
	}
	return c.GetArchived() == other.GetArchived() && c.GetNote() == other.GetNote()
}

func NewProjectCardOpsList(sourceCards, targetCards []*github.ProjectCard, sourceColumn, targetColumn *github.ProjectColumn) ProjectCardOpsList {
	if len(sourceCards) == 0 && len(targetCards) == 0 {
		return nil
	}

	kinds := opMapping{}
	for _, s := range sourceCards {
		src := &projectCard{s}
		kinds.requestCreate(src)
		for _, t := range targetCards {
			target := &projectCard{t}
			if src.Key().Eq(target.Key()) {
				kinds.requestNothing(src)
			}
		}
	}

	ops := []*ProjectCardOp{}
	for _, s := range sourceCards {
		src := &projectCard{s}
		switch kinds.get(src) {
		case OpCreate:
			ops = append(ops, &ProjectCardOp{
				Kind:          OpCreate,
				ProjectCard:   s,
				ProjectColumn: targetColumn,
			})
		default:
		}
	}
	return ops
}

type ProjectCardOpsList []*ProjectCardOp

func (l ProjectCardOpsList) String() string {
	s := "["
	for _, op := range l {
		s += fmt.Sprintf("%s, ", op)
	}
	s += "]"
	return s
}

type ProjectCardOp struct {
	Kind          OpKind
	ProjectCard   *github.ProjectCard
	ProjectColumn *github.ProjectColumn
}

func (op *ProjectCardOp) String() string {
	return stringify(op.Kind, op.ProjectCard)
}
