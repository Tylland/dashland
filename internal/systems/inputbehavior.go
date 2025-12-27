package systems

import (
	"fmt"

	"github.com/tylland/dashland/internal/common"
	"github.com/tylland/dashland/internal/components"
	"github.com/tylland/dashland/internal/ecs"
	"github.com/tylland/dashland/internal/game"
)

const walkSpeed = 128 // pixels per second

type InputBehavior struct {
	stage *game.Stage
}

func NewInputBehavior(stage *game.Stage) *InputBehavior {
	return &InputBehavior{
		stage: stage,
	}
}

func (s *InputBehavior) Update(world *ecs.World, deltaTime float32) {
	for _, entity := range world.Entities() {
		input := ecs.GetComponent[components.InputComponent](entity)
		position := ecs.GetComponent[components.PositionComponent](entity)
		step := ecs.GetComponent[components.BlockStep](entity)

		if input != nil && position != nil && step != nil {
			s.nextStep(input, position, step)
		}
	}
}

func (s *InputBehavior) nextStep(input *components.InputComponent, position *components.PositionComponent, step *components.BlockStep) {

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
			step.Move(direction, walkSpeed)
		}
	}
}
