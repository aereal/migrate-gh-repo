package domain

import (
	"fmt"

	"github.com/google/go-github/github"
)

func eqMilestone(l *github.Milestone, r *github.Milestone) bool {
	if l == nil || r == nil {
		return false
	}
	return milestoneID(l) == milestoneID(r)
}

func milestoneID(milestone *github.Milestone) string {
	return milestone.GetTitle()
}

func NewMilestoneOpsList(sourceMilestones, targetMilestones []*github.Milestone) MilestoneOpsList {
	if len(sourceMilestones) == 0 && len(targetMilestones) == 0 {
		return nil
	}

	kinds := map[string]OpKind{}
	for _, src := range sourceMilestones {
		id := milestoneID(src)
		kinds[id] = OpCreate
		for _, tgt := range targetMilestones {
			if eqMilestone(src, tgt) {
				if src.String() == tgt.String() { // completely equal
					kinds[id] = OpNothing
				} else {
					kinds[id] = OpUpdate
				}
			}
		}
	}

	ops := []*MilestoneOp{}
	for _, src := range sourceMilestones {
		switch kinds[milestoneID(src)] {
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

type OpKind string

const (
	OpCreate  = OpKind("create")
	OpUpdate  = OpKind("update")
	OpNothing = OpKind("nothing")
)

type MilestoneOpsList []*MilestoneOp

func (l MilestoneOpsList) String() string {
	s := "["
	for _, op := range l {
		s += op.String()
		s += ", "
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

func stringify(kind OpKind, payload interface{}) string {
	return fmt.Sprintf(`{"kind":%q, "payload":%s}`, kind, payload)
}
