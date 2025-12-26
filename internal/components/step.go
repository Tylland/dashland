package components

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/tylland/dashland/internal/common"
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
