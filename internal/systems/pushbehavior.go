package systems

import (
	"github.com/tylland/dashland/internal/components"
	"github.com/tylland/dashland/internal/ecs"
	"github.com/tylland/dashland/internal/game"
)

type PushBehavior struct {
	stage *game.Stage
}

func NewPushBehavior(stage *game.Stage) *PushBehavior {
	return &PushBehavior{stage: stage}
}

func (p *PushBehavior) Update(world *ecs.World, deltaTime float32) {
	player := world.GetEntity("player")
	if player != nil {
		position := ecs.GetComponent[components.PositionComponent](player)
		step := ecs.GetComponent[components.BlockStep](player)

		if position != nil && step != nil {
			p.handlePlayerPush(world, position, step)
		}
	}
}

func (s *PushBehavior) handlePlayerPush(world *ecs.World, playerPosition *components.PositionComponent, playerStep *components.BlockStep) {

	targetBlockPos := playerPosition.CurrentBlockPosition.Add(playerStep.Increment)

	if target := s.stage.GetEntity(targetBlockPos); target != nil {
		targetPosition := ecs.GetComponent[components.PositionComponent](target)
		targetPushable := ecs.GetComponent[components.PushableComponent](target)
		targetStep := ecs.GetComponent[components.BlockStep](target)

		if targetPosition == nil || targetPushable == nil {
			return
		}

		if !playerStep.Increment.IsZero() && playerStep.Increment.IsHorizontal() {
			pushPos := targetPosition.CurrentBlockPosition.Add(playerStep.Increment)

			if s.stage.CheckBlockAtPosition(game.Void, pushPos) && s.stage.GetEntity(pushPos) == nil {
				targetStep.Move(playerStep.Increment, moveSpeed)
			}
		}
	}
}
