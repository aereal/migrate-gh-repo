package domain

import (
	"fmt"

	"github.com/google/go-github/github"
)

type projectColumn struct {
	*github.ProjectColumn
}

func (c *projectColumn) Key() *Key {
	if c == nil {
		return nil
	}
	return &Key{
		kind: "project_column",
		repr: c.GetName(),
	}
}

func (c *projectColumn) eq(other *projectColumn) bool {
	if c == nil || other == nil {
		return false
	}
	if !c.Key().Eq(other.Key()) {
		return false
	}
	return c.GetName() == other.GetName()
}

func NewProjectColumnOpsList(sourceColumns, targetColumns []*github.ProjectColumn) ProjectColumnOpsList {
	if len(sourceColumns) == 0 && len(targetColumns) == 0 {
		return nil
	}

	kinds := opMapping{}
	mapping := map[string]*github.ProjectColumn{}
	for _, s := range sourceColumns {
		src := &projectColumn{s}
		kinds.requestCreate(src)
		for _, t := range targetColumns {
			target := &projectColumn{t}
			if src.Key().Eq(target.Key()) {
				kinds.requestUpdate(src)
				mapping[src.Key().String()] = t
			}
		}
	}

	ops := []*ProjectColumnOp{}
	for _, s := range sourceColumns {
		src := &projectColumn{s}
		switch kinds.get(src) {
		case OpCreate:
			ops = append(ops, &ProjectColumnOp{
				Kind:          OpCreate,
				ProjectColumn: s,
			})
		case OpUpdate:
			op := &ProjectColumnOp{
				Kind:          OpUpdate,
				ProjectColumn: s,
			}
			ops = append(ops, op)
		default:
		}
	}
	return ops
}

type ProjectColumnOpsList []*ProjectColumnOp

func (l ProjectColumnOpsList) String() string {
	s := "["
	for _, op := range l {
		s += fmt.Sprintf("%s, ", op)
	}
	s += "]"
	return s
}

type ProjectColumnOp struct {
	Kind          OpKind
	ProjectColumn *github.ProjectColumn
}

func (op *ProjectColumnOp) String() string {
	return stringify(op.Kind, op.ProjectColumn)
}
