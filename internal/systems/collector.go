package systems

import (
	"fmt"

	"github.com/tylland/dashland/internal/ecs"
	"github.com/tylland/dashland/internal/game"
)

type Collector struct {
	stage *game.Stage
	sound game.SoundPlayer
}

func NewCollect(stage *game.Stage, sound game.SoundPlayer) *Collector {
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
	s.sound.PlayFx("diamond_collected")

	s.stage.TryRemoveEntity(collect.Collectable)
	world.EnqueueRemoval(collect.Collectable)
}
