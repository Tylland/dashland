package systems

import (
	"github.com/tylland/dashland/internal/components"
	"github.com/tylland/dashland/internal/ecs"
	"github.com/tylland/dashland/internal/game"
)

type Cleanup struct {
	stage *game.Stage
}

func NewCleanup(stage *game.Stage) *Cleanup {
	return &Cleanup{stage: stage}
}

func (s *Cleanup) Update(world *ecs.World, deltaTime float32) {

	for _, doomed := range world.RemovalQueue() {
		if position := ecs.GetComponent[components.PositionComponent](doomed.Components); position != nil {
			s.stage.RemoveEntity(doomed, position.CurrentBlockPosition)
		}

		world.RemoveEntity(doomed)
	}

	world.ResetRemovalQueue()
	world.ClearEvents()
}
