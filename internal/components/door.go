package components

import (
	"github.com/tylland/dashland/internal/common"
	"github.com/tylland/dashland/internal/ecs"
)

type DoorState uint32

const (
	DoorHidden DoorState = iota
	DoorClosed
	DoorOpen
)

func (s DoorState) String() string {
	switch s {
	case DoorHidden:
		return "Hidden"
	case DoorClosed:
		return "Closed"
	case DoorOpen:
		return "Open"
	default:
		return "Unknown"
	}
}

type DoorComponent struct {
	Stage    string
	Position common.BlockPosition
	State    DoorState
}

func NewDoorComponent(stage string, pos common.BlockPosition, state DoorState) *DoorComponent {
	return &DoorComponent{Stage: stage, Position: pos, State: state}
}

func (c *DoorComponent) WriteTo(w ecs.Writer) error {
	w.WriteLabel(ecs.ComponentName(c))
	w.WriteString(c.Stage, false)
	w.WriteInt(c.Position.X, false)
	w.WriteInt(c.Position.Y, false)
	w.WriteEnum(c.State, false, "state")

	return nil
}

func (c *DoorComponent) ReadFrom(r ecs.Reader) error {
	r.ConsumeLabel(ecs.ComponentName(c))

	ecs.WithChangedValue(r.ReadString, func(stage string) {
		c.Stage = stage
	})

	ecs.WithChangedValue(r.ReadInt, func(x int) {
		c.Position.X = x
	})

	ecs.WithChangedValue(r.ReadInt, func(y int) {
		c.Position.Y = y
	})

	ecs.WithChangedValue(r.ReadEnum, func(val uint32) {
		c.State = DoorState(val)
	})

	return nil
}
