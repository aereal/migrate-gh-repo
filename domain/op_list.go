package domain

import (
	"fmt"

	"github.com/google/go-github/github"
)

type Key struct {
	kind string
	repr string
}

func (k *Key) String() string {
	return fmt.Sprintf("kind=%s;repr=%s", k.kind, k.repr)
}

func (k *Key) Eq(other *Key) bool {
	if k == nil || other == nil {
		return false
	}
	return k.String() == other.String()
}

type Equalable interface {
	Key() *Key
	Eq(other Equalable) bool
}

type milestone struct {
	*github.Milestone
}

func (m *milestone) Key() *Key {
	if m == nil {
		return nil
	}
	return &Key{kind: "milestone", repr: m.GetTitle()}
}

func (m *milestone) Eq(other Equalable) bool {
	if m == nil || other == nil {
		return false
	}
	if !m.Key().Eq(other.Key()) {
		return false
	}
	if otherMilestone, ok := other.(*milestone); ok {
		return m.String() == otherMilestone.String()
	}
	return false
}

func NewMilestoneOpsList(sourceMilestones, targetMilestones []*github.Milestone) MilestoneOpsList {
	if len(sourceMilestones) == 0 && len(targetMilestones) == 0 {
		return nil
	}

	kinds := map[string]OpKind{}
	for _, src := range sourceMilestones {
		srcm := &milestone{src}
		kinds[srcm.Key().String()] = OpCreate
		for _, tgt := range targetMilestones {
			tgtm := &milestone{tgt}
			if srcm.Key().Eq(tgtm.Key()) {
				if srcm.Eq(tgtm) { // completely equal
					kinds[srcm.Key().String()] = OpNothing
				} else {
					kinds[srcm.Key().String()] = OpUpdate
				}
			}
		}
	}

	ops := []*MilestoneOp{}
	for _, src := range sourceMilestones {
		srcm := &milestone{src}
		switch kinds[srcm.Key().String()] {
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
