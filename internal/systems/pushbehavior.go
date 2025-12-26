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
		position := ecs.GetComponent[components.PositionComponent](player.Components)
		step := ecs.GetComponent[components.BlockStep](player.Components)

		if position != nil && step != nil {
			p.handlePlayerPush(world, position, step)
		}
	}
}

func (s *PushBehavior) handlePlayerPush(world *ecs.World, playerPosition *components.PositionComponent, playerStep *components.BlockStep) {

	targetBlockPos := playerPosition.CurrentBlockPosition.Add(playerStep.Increment)

	if target := s.stage.GetEntity(targetBlockPos); target != nil {
		targetPosition := ecs.GetComponent[components.PositionComponent](target.Components)
		targetPushable := ecs.GetComponent[components.PushableComponent](target.Components)
		targetStep := ecs.GetComponent[components.BlockStep](target.Components)

		if targetPosition == nil || targetPushable == nil {
			return
		}

		// Calculate push direction based on player's movement
		if !playerStep.Increment.IsZero() {
			pushPos := targetPosition.CurrentBlockPosition.Add(playerStep.Increment)

			// if playerPosition.CurrentBlockPosition.X > targetPosition.CurrentBlockPosition.X {
			// 	pushPos = pushPos.Offset(-1, 0)
			// } else if playerPosition.CurrentBlockPosition.X < targetPosition.CurrentBlockPosition.X {
			// 	pushPos = pushPos.Offset(1, 0)
			// }

			if s.stage.CheckBlockAtPosition(game.Void, pushPos) && s.stage.GetEntity(pushPos) == nil {
				targetStep.Move(playerStep.Increment, moveSpeed)
			}
		}
	}
}
