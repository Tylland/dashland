package systems

import (
	"fmt"

	"github.com/tylland/dashland/internal/components"
	"github.com/tylland/dashland/internal/ecs"
	"github.com/tylland/dashland/internal/game"
)

type CollectorSystem struct {
	stage *game.Stage
	sound game.SoundPlayer
}

func NewCollectorSystem(stage *game.Stage, sound game.SoundPlayer) *CollectorSystem {
	return &CollectorSystem{stage: stage, sound: sound}
}

func (s *CollectorSystem) Update(world *ecs.World, deltaTime float32) {
	for _, event := range world.Events() {
		if event != nil && event.Name == "collect" {
			s.handleCollect(world, event.Data.(*game.CollectEvent))
		}
	}
}

func (s *CollectorSystem) handleCollect(world *ecs.World, collect *game.CollectEvent) {
	fmt.Printf("Player collected %s!!\n", collect.Collectable.ID)

	collactable := ecs.GetComponent[components.CollectableComponent](collect.Collectable)

	if collactable != nil && collactable.Name == "diamond" {
		inventory := ecs.GetComponent[components.Inventory](collect.Collector)

		if inventory != nil {
			inventory.Diamonds += collactable.Amount

			if inventory.Diamonds >= s.stage.DiamondsRequired && !s.stage.ExitCondition {
				s.stage.ExitCondition = true
				world.AddEvent("exitopen", game.NewExitOpenEvent())
			}

		}

		s.sound.PlayFx("diamond_collected")
	}

	s.stage.TryRemoveEntity(collect.Collectable)
	world.EnqueueRemoval(collect.Collectable)
}
