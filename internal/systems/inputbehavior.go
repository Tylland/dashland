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
			s.nextStep(entity, input, position, step)
		}
	}
}

func (s *InputBehavior) nextStep(entity *ecs.Entity, input *components.InputComponent, position *components.PositionComponent, step *components.BlockStep) {

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

		var animation string = "idle"

		if !direction.IsZero() {
			step.Move(direction, walkSpeed)

			switch direction {
			case common.DirectionRight:
				animation = "walk_right"
			case common.DirectionLeft:
				animation = "walk_left"
			case common.DirectionUp:
				animation = "walk_up"
			case common.DirectionDown:
				animation = "walk_down"
			default:
				animation = "idle"
			}
		} else {
			switch step.Direction {
			case common.DirectionRight:
				animation = "stand_right"
			case common.DirectionLeft:
				animation = "stand_left"
			case common.DirectionUp:
				animation = "stand_up"
			case common.DirectionDown:
				animation = "stand_down"
			default:
				animation = "idle"
			}
		}

		// set walk animation based on direction
		anim := ecs.GetComponent[components.AnimationComponent](entity)
		sprite := ecs.GetComponent[components.SpriteComponent](entity)

		if anim != nil && sprite != nil {
			anim.Current = animation
			anim.ApplyAnimation(sprite.Sprite)
		}

	}
}
