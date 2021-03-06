package domain

import (
	"fmt"

	"github.com/google/go-github/github"
)

type label struct {
	*github.Label
}

func (l *label) Key() *Key {
	if l == nil {
		return nil
	}
	return &Key{kind: "label", repr: l.GetName()}
}

func (l *label) eq(other *label) bool {
	if l == nil || other == nil {
		return false
	}
	if !l.Key().Eq(other.Key()) {
		return false
	}
	return l.GetColor() == other.GetColor() && l.GetDescription() == other.GetDescription()
}

type LabelOpsList []*LabelOp

func (l LabelOpsList) String() string {
	s := "["
	for _, op := range l {
		s += fmt.Sprintf("%s, ", op)
	}
	s += "]"
	return s
}

type LabelOp struct {
	Kind  OpKind
	Label *github.Label
}

func (op *LabelOp) String() string {
	return stringify(op.Kind, op.Label)
}

func NewLabelOpsList(sourceLabels, targetLabels []*github.Label) LabelOpsList {
	if len(sourceLabels) == 0 && len(targetLabels) == 0 {
		return nil
	}

	kinds := opMapping{}
	for _, src := range sourceLabels {
		srcm := &label{src}
		kinds.requestCreate(srcm)
		for _, tgt := range targetLabels {
			tgtm := &label{tgt}
			if srcm.Key().Eq(tgtm.Key()) {
				if srcm.eq(tgtm) { // completely equal
					kinds.requestNothing(srcm)
				} else {
					kinds.requestNothing(srcm)
				}
			}
		}
	}

	ops := []*LabelOp{}
	for _, src := range sourceLabels {
		srcm := &label{src}
		switch kinds.get(srcm) {
		case OpCreate:
			ops = append(ops, &LabelOp{
				Kind:  OpCreate,
				Label: src,
			})
		case OpUpdate:
			ops = append(ops, &LabelOp{
				Kind:  OpUpdate,
				Label: src,
			})
		default:
		}
	}
	return LabelOpsList(ops)
}
