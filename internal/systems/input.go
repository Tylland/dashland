package systems

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/tylland/dashland/internal/components"
	"github.com/tylland/dashland/internal/ecs"
)

type InputSystem struct {
}

func NewInputSystem() *InputSystem {
	return &InputSystem{}
}

func (s *InputSystem) Update(world *ecs.World, deltaTime float32) {
	for _, entity := range world.Entities() {
		input := ecs.GetComponent[components.InputComponent](entity)

		if input != nil {
			s.updateInput(input)
		}

	}
}

func (s *InputSystem) updateInput(input *components.InputComponent) {
	input.RightKeyPressed = rl.IsKeyDown(rl.KeyRight) || rl.IsKeyDown(rl.KeyD)
	input.LeftKeyPressed = rl.IsKeyDown(rl.KeyLeft) || rl.IsKeyDown(rl.KeyA)
	input.DownKeyPressed = rl.IsKeyDown(rl.KeyDown) || rl.IsKeyDown(rl.KeyS)
	input.UpKeyPressed = rl.IsKeyDown(rl.KeyUp) || rl.IsKeyDown(rl.KeyW)
}
