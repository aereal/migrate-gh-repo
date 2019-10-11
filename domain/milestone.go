package domain

import (
	"fmt"

	"github.com/google/go-github/github"
)

type milestone struct {
	*github.Milestone
}

func (m *milestone) Key() *Key {
	if m == nil {
		return nil
	}
	return &Key{kind: "milestone", repr: m.GetTitle()}
}

func (m *milestone) eq(other *milestone) bool {
	if m == nil || other == nil {
		return false
	}
	if !m.Key().Eq(other.Key()) {
		return false
	}
	return m.GetDescription() == other.GetDescription() && m.GetDueOn() == other.GetDueOn()
}

func NewMilestoneOpsList(sourceMilestones, targetMilestones []*github.Milestone) MilestoneOpsList {
	if len(sourceMilestones) == 0 && len(targetMilestones) == 0 {
		return nil
	}

	kinds := opMapping{}
	for _, src := range sourceMilestones {
		srcm := &milestone{src}
		kinds.requestCreate(srcm)
		for _, tgt := range targetMilestones {
			tgtm := &milestone{tgt}
			if srcm.Key().Eq(tgtm.Key()) {
				if srcm.eq(tgtm) { // completely equal
					kinds.requestNothing(srcm)
				} else {
					kinds.requestUpdate(srcm)
				}
			}
		}
	}

	ops := []*MilestoneOp{}
	for _, src := range sourceMilestones {
		srcm := &milestone{src}
		switch kinds.get(srcm) {
		case OpCreate:
			ops = append(ops, &MilestoneOp{
				Kind:      OpCreate,
				Milestone: src,
			})
		case OpUpdate:
			ops = append(ops, &MilestoneOp{
				Kind:      OpUpdate,
				Milestone: src,
			})
		default:
		}
	}
	return MilestoneOpsList(ops)
}

type MilestoneOpsList []*MilestoneOp

func (l MilestoneOpsList) String() string {
	s := "["
	for _, op := range l {
		s += fmt.Sprintf("%s, ", op)
	}
	s += "]"
	return s
}

type MilestoneOp struct {
	Kind      OpKind
	Milestone *github.Milestone
}

func (op *MilestoneOp) String() string {
	return stringify(op.Kind, op.Milestone)
}
