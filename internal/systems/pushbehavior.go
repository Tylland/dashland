package systems

import (
	"github.com/tylland/dashland/internal/components"
	"github.com/tylland/dashland/internal/ecs"
	"github.com/tylland/dashland/internal/game"
	"github.com/tylland/dashland/internal/game/types"
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
		position := ecs.GetComponent[components.Position](player)
		step := ecs.GetComponent[components.BlockStep](player)

		if position != nil && step != nil {
			p.handlePlayerPush(world, position, step)
		}
	}
}

func (s *PushBehavior) handlePlayerPush(world *ecs.World, playerPosition *components.Position, playerStep *components.BlockStep) {

	targetBlockPos := playerPosition.CurrentBlockPosition.Add(playerStep.Increment)

	if target := s.stage.GetEntity(targetBlockPos); target != nil {
		targetPosition := ecs.GetComponent[components.Position](target)
		targetPushable := ecs.GetComponent[components.Pushable](target)
		targetStep := ecs.GetComponent[components.BlockStep](target)

		if targetPosition == nil || targetPushable == nil {
			return
		}

		if !playerStep.Increment.IsZero() && playerStep.Increment.IsHorizontal() {
			pushPos := targetPosition.CurrentBlockPosition.Add(playerStep.Increment)

			if s.stage.CheckBlockAtPosition(types.EntityVoid, pushPos) && s.stage.GetEntity(pushPos) == nil {
				targetStep.Move(playerStep.Increment, moveSpeed)
			}
		}
	}
}
