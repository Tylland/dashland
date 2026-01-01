package systems

import (
	"github.com/tylland/dashland/internal/common"
	"github.com/tylland/dashland/internal/components"
	"github.com/tylland/dashland/internal/ecs"
	"github.com/tylland/dashland/internal/game"
)

type StageLoader interface {
	LoadStage(name string, pos common.BlockPosition) error
}

type Cleanup struct {
	stage  *game.Stage
	loader StageLoader
}

func NewCleanup(stage *game.Stage, loader StageLoader) *Cleanup {
	return &Cleanup{stage: stage, loader: loader}
}

func (s *Cleanup) Update(world *ecs.World, deltaTime float32) {
	for _, doomed := range world.RemovalQueue() {
		if position := ecs.GetComponent[components.PositionComponent](doomed); position != nil {
			s.stage.RemoveEntity(doomed, position.CurrentBlockPosition)
		}

		world.RemoveEntity(doomed)
	}

	world.ResetRemovalQueue()

	for _, event := range world.Events() {
		if event != nil && event.Name == "stagechange" {
			defer s.handleStageChange(event.Data.(*game.StageChangeEvent))
		}
	}

	world.ClearEvents()
}

func (s *Cleanup) handleStageChange(statechance *game.StageChangeEvent) {
	s.loader.LoadStage(statechance.Stage, statechance.Position)
}
