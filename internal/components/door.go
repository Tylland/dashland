package components

import "github.com/tylland/dashland/internal/common"

type DoorComponent struct {
	Stage    string
	Position common.BlockPosition
}

func NewDoorComponent(stage string, pos common.BlockPosition) *DoorComponent {
	return &DoorComponent{Stage: stage, Position: pos}
}
