package components

import "github.com/tylland/dashland/internal/common"

type DoorState int

const (
	DoorHidden DoorState = iota
	DoorClosed
	DoorOpen
)

type DoorComponent struct {
	Stage    string
	Position common.BlockPosition
	State    DoorState
	SourceX  float32 // base X in the tileset for the closed state
	OpenX    float32 // base X in the tileset for the open state
}

func NewDoorComponent(stage string, pos common.BlockPosition, state DoorState) *DoorComponent {
	return &DoorComponent{Stage: stage, Position: pos, State: state}
}
