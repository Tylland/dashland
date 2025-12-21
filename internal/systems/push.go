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
	playerStep := ecs.GetComponent[components.BlockStep](playerComps)

	p.handlePlayerPush(world, playerPosition, playerStep)
}

func (s *PushSystem) handlePlayerPush(world *ecs.World, playerPosition *components.PositionComponent, playerStep *components.BlockStep) {

	targetBlockPos := playerPosition.CurrentBlockPosition.Add(playerStep.Increment)

	if targetEntity := s.stage.GetEntity(targetBlockPos); targetEntity != nil {
		comps := world.GetComponents(targetEntity)

		targetPosition := ecs.GetComponent[components.PositionComponent](comps)
		targetPushable := ecs.GetComponent[components.PushableComponent](comps)
		targetStep := ecs.GetComponent[components.BlockStep](comps)

		if targetPosition == nil || targetPushable == nil {
			return
		}

		// Calculate push direction based on player's movement
		if !playerStep.Increment.IsZero() {
			pushPos := targetPosition.CurrentBlockPosition.Add(playerStep.Increment)

			if playerPosition.CurrentBlockPosition.X > targetPosition.CurrentBlockPosition.X {
				pushPos = pushPos.Offset(-1, 0)
			} else if playerPosition.CurrentBlockPosition.X < targetPosition.CurrentBlockPosition.X {
				pushPos = pushPos.Offset(1, 0)
			}

			if s.stage.CheckBlockAtPosition(game.Void, pushPos) && s.stage.GetEntity(pushPos) == nil {
				targetStep.Move(playerStep.Increment, s.stage.GetPosition(pushPos), moveSpeed)
			}
		}
	}
}
