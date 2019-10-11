package domain

import (
	"fmt"
)

type hasKey interface {
	Key() *Key
}

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

type OpKind string

const (
	OpCreate  = OpKind("create")
	OpUpdate  = OpKind("update")
	OpNothing = OpKind("nothing")
)

func stringify(kind OpKind, payload interface{}) string {
	return fmt.Sprintf(`{"kind":%q, "payload":%s}`, kind, payload)
}

type opMapping map[string]OpKind

func (m opMapping) get(hk hasKey) OpKind {
	return m[hk.Key().String()]
}

func (m opMapping) requestCreate(hk hasKey) {
	m[hk.Key().String()] = OpCreate
}

func (m opMapping) requestUpdate(hk hasKey) {
	m[hk.Key().String()] = OpUpdate
}

func (m opMapping) requestNothing(hk hasKey) {
	m[hk.Key().String()] = OpNothing
}
