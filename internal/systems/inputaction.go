package systems

import (
	"fmt"

	"github.com/tylland/dashland/internal/common"
	"github.com/tylland/dashland/internal/components"
	"github.com/tylland/dashland/internal/ecs"
	"github.com/tylland/dashland/internal/game"
)

const walkSpeed = 128 // pixels per second

type InputActionSystem struct {
	stage *game.Stage
}

func NewInputActionSystem(stage *game.Stage) *InputActionSystem {
	return &InputActionSystem{
		stage: stage,
	}
}

func (s *InputActionSystem) Update(world *ecs.World, deltaTime float32) {
	for _, entity := range world.Entities() {
		comps := world.GetComponents(entity)

		input := ecs.GetComponent[components.InputComponent](comps)
		position := ecs.GetComponent[components.PositionComponent](comps)
		step := ecs.GetComponent[components.BlockStep](comps)

		if input != nil && position != nil && step != nil {
			s.nextStep(input, position, step)
		}
	}
}

func (s *InputActionSystem) nextStep(input *components.InputComponent, position *components.PositionComponent, step *components.BlockStep) {

	if position.Vector2 == step.Target {
		var direction common.BlockVector = common.BlockVector{}

		if input.RightKeyPressed {
			fmt.Print("Move player right!")
			direction = common.DirectionRight
		} else if input.LeftKeyPressed {
			fmt.Print("Move player left!")
			direction = common.DirectionLeft
		} else if input.DownKeyPressed {
			fmt.Print("Move player down!")
			direction = common.DirectionDown
		} else if input.UpKeyPressed {
			fmt.Print("Move player up!")
			direction = common.DirectionUp
		}

		if !direction.IsZero() {
			targetBlockPos := position.CurrentBlockPosition.Add(direction)

			step.Move(direction, s.stage.GetPosition(targetBlockPos), walkSpeed)
		}
	}
}
