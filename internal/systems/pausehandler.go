package systems

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/tylland/dashland/internal/ecs"
)

type PauseSystem struct {
}

func NewPauseSystem() *PauseSystem {
	return &PauseSystem{}
}

func (s *PauseSystem) Update(world *ecs.World, deltaTime float32) {
	if rl.IsKeyReleased(rl.KeyP) || rl.IsKeyReleased(rl.KeyPause) {
		world.SetPaused(!world.Paused())
	}

}

func (s *PauseSystem) RunWhenPaused() bool {
	return true
}
