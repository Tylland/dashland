package components

import (
	"github.com/tylland/dashland/internal/common"
	"github.com/tylland/dashland/internal/ecs"
	"github.com/tylland/dashland/internal/game/types"
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
	w.StartSection(Section(types.ComponentTypeDoor))

	w.WriteString("Stage", c.Stage, false)
	w.WriteInt("X", c.Position.X, false)
	w.WriteInt("Y", c.Position.Y, false)
	w.WriteEnum("State", c.State, false, "state")

	return nil
}

func (c *DoorComponent) ReadFrom(r ecs.Reader) error {
	r.ConsumeSection(types.ComponentTypeDoor.Name)

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
