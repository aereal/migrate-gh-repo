package domain

import (
	"fmt"
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

type OpKind string

const (
	OpCreate  = OpKind("create")
	OpUpdate  = OpKind("update")
	OpNothing = OpKind("nothing")
)

func stringify(kind OpKind, payload interface{}) string {
	return fmt.Sprintf(`{"kind":%q, "payload":%s}`, kind, payload)
}
