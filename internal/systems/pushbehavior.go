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
		if position, step, ok := player.ComponentTuple[components.Position, components.BlockStep](); ok {
			p.handlePlayerPush(position, step)
		}
	}
}

func (s *PushBehavior) handlePlayerPush(playerPosition *components.Position, playerStep *components.BlockStep) {
	targetBlockPos := playerPosition.CurrentBlockPosition.Add(playerStep.Increment)

	if target := s.stage.GetEntity(targetBlockPos); target != nil {
		if targetPosition, _, targetStep, ok := target.ComponentTriple[components.Position, components.Pushable, components.BlockStep](); ok {
			if !playerStep.Increment.IsZero() && playerStep.Increment.IsHorizontal() {
				pushPos := targetPosition.CurrentBlockPosition.Add(playerStep.Increment)

				if s.stage.CheckBlockAtPosition(types.EntityVoid, pushPos) && s.stage.GetEntity(pushPos) == nil {
					targetStep.Move(playerStep.Increment, moveSpeed)
				}
			}
		}
	}
}
