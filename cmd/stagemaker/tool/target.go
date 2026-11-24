package tool

import (
	"github.com/tylland/dashland/internal/common"
)

type TargetKind uint8

const (
	TargetBlock TargetKind = iota + 1
	TargetEntity
)

type Target struct {
	Kind     TargetKind
	Position common.BlockPosition
}

func SameTarget(a, b *Target) bool {
	return a == nil && b == nil || a != nil && b != nil && *a == *b
}
