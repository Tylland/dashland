package components

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/tylland/dashland/internal/common"
	"github.com/tylland/dashland/internal/ecs"
	"github.com/tylland/dashland/internal/game/types"
)

type BlockStep struct {
	Increment         common.BlockVector
	Direction         common.BlockVector
	previousDirection common.BlockVector
	Target            rl.Vector2
	Speed             float32
}

func NewBlockStep(target rl.Vector2) *BlockStep {
	return &BlockStep{Increment: common.NewBlockVector(0, 0), Target: target}
}

func (s *BlockStep) Move(increment common.BlockVector, speed float32) {
	s.Increment = increment
	s.SetDirection(increment)
	s.Speed = speed
}

func (s *BlockStep) SetDirection(dir common.BlockVector) {
	s.previousDirection = s.Direction
	s.Direction = dir
}

func (s *BlockStep) Cancel() {
	s.Increment.Clear()
	//s.Direction = s.previousDirection
}

func (s *BlockStep) Commit(target rl.Vector2) {
	s.Increment.Clear()
	s.Target = target
}

func (c *BlockStep) WriteTo(w ecs.Writer) error {
	w.StartSection(Section(types.ComponentTypeBlockStep))

	w.WriteInt("Inc X", c.Increment.X, false)
	w.WriteInt("Inc Y", c.Increment.Y, false)

	w.WriteInt("Dir X", c.Direction.X, false)
	w.WriteInt("Dir Y", c.Direction.Y, false)

	w.WriteInt("Prev Dir X", c.previousDirection.X, false)
	w.WriteInt("Prev Dir Y", c.previousDirection.Y, false)

	w.WriteFloat("Target X", c.Target.X, false)
	w.WriteFloat("Target Y", c.Target.Y, false)

	w.WriteFloat("Speed", c.Speed, false)

	return nil
}

func (c *BlockStep) ReadFrom(r ecs.Reader) error {
	r.ConsumeSection(types.ComponentTypeBlockStep.Name)
	r.ConsumeLine()
	r.ConsumeLine()

	r.ConsumeLine()
	r.ConsumeLine()

	r.ConsumeLine()
	r.ConsumeLine()

	ecs.WithChangedValue(r.ReadFloat, func(x float32) {
		c.Target.X = x
	})

	ecs.WithChangedValue(r.ReadFloat, func(y float32) {
		c.Target.Y = y
	})

	ecs.WithChangedValue(r.ReadFloat, func(speed float32) {
		c.Speed = speed
	})

	return nil
}
