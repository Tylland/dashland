package systems

import (
	"github.com/tylland/dashland/internal/components"
	"github.com/tylland/dashland/internal/ecs"
	"github.com/tylland/dashland/internal/game"
)

type PushSystem struct {
	stage *game.Stage
}

func NewPushSystem(stage *game.Stage) *PushSystem {
	return &PushSystem{stage: stage}
}

func (p *PushSystem) Update(world *ecs.World, deltaTime float32) {
	player := world.GetEntity("player")
	playerComps := world.GetComponents(player)
	playerPosition := ecs.GetComponent[components.PositionComponent](playerComps)

	p.handlePlayerPush(world, playerPosition)
}

func (s *PushSystem) handlePlayerPush(world *ecs.World, playerPosition *components.PositionComponent) {
	if targetEntity := s.stage.GetEntity(playerPosition.TargetBlockPosition); targetEntity != nil {
		comps := world.GetComponents(targetEntity)

		targetPosition := ecs.GetComponent[components.PositionComponent](comps)
		targetPushable := ecs.GetComponent[components.PushableComponent](comps)

		if targetPosition == nil || targetPushable == nil {
			return
		}

		// Calculate push direction based on player's movement
		if playerPosition.HasTarget() {
			pushPos := targetPosition.TargetBlockPosition

			if playerPosition.CurrentBlockPosition.X > targetPosition.TargetBlockPosition.X {
				pushPos = pushPos.Offset(-1, 0)
			} else if playerPosition.CurrentBlockPosition.X < targetPosition.TargetBlockPosition.X {
				pushPos = pushPos.Offset(1, 0)
			}

			// Check if push position is free
			if s.stage.CheckBlockAtPosition(game.Void, pushPos) && s.stage.GetEntity(pushPos) == nil {
				// Move the pushed entity
				s.stage.MoveEntity(targetEntity, targetPosition.CurrentBlockPosition, pushPos)
				targetPosition.Update(pushPos)

				targetPosition.TargetBlockPosition = pushPos
			}
		}
	}
}
