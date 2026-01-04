package systems

import (
	"fmt"

	"github.com/tylland/dashland/internal/components"
	"github.com/tylland/dashland/internal/ecs"
	"github.com/tylland/dashland/internal/game"
)

type Collector struct {
	stage *game.Stage
	sound game.SoundPlayer
}

func NewCollector(stage *game.Stage, sound game.SoundPlayer) *Collector {
	return &Collector{stage: stage, sound: sound}
}

func (s *Collector) Update(world *ecs.World, deltaTime float32) {
	for _, event := range world.Events() {
		if event != nil && event.Name == "collect" {
			s.handleCollect(world, event.Data.(*game.CollectEvent))
		}
	}
}

func (s *Collector) handleCollect(world *ecs.World, collect *game.CollectEvent) {
	fmt.Printf("Player collected %s!!\n", collect.Collectable.ID)

	collactable := ecs.GetComponent[components.CollectableComponent](collect.Collectable)

	if collactable != nil && collactable.Name == "diamond" {
		inventory := ecs.GetComponent[components.Inventory](collect.Collector)

		if inventory != nil {
			inventory.Diamonds += collactable.Amount

			if inventory.Diamonds >= s.stage.DiamondsRequired && !s.stage.ExitCondition {
				game.NewFlash(world)
				s.stage.ExitCondition = true
			}

		}

		s.sound.PlayFx("diamond_collected")
	}

	s.stage.TryRemoveEntity(collect.Collectable)
	world.EnqueueRemoval(collect.Collectable)
}
