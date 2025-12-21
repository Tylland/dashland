package components

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/tylland/dashland/internal/common"
)

type BlockStep struct {
	Increment common.BlockVector
	Target    rl.Vector2
	Speed     float32
}

func NewBlockStep(target rl.Vector2) *BlockStep {
	return &BlockStep{Increment: common.NewBlockVector(0, 0), Target: target}
}

func (s *BlockStep) Move(increment common.BlockVector, target rl.Vector2, speed float32) {
	s.Increment = increment
	s.Target = target
	s.Speed = speed
}

// func (s *BlockStep) Halt() {
// 	s.Increment.Clear()
// }

// func (s *BlockStep) HasIncrement() bool {
// 	return s.Increment.X != 0 || s.Increment.Y != 0
// }

// func (s *BlockStep) AtTarget() bool {
// 	return s.Increment.X != 0 || s.Increment.Y != 0
// }
