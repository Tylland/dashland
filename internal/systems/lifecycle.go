package systems

import (
	"github.com/tylland/dashland/internal/components"
	"github.com/tylland/dashland/internal/ecs"
)

// AnimationSystem applies AnimationComponent selections to Sprite components.
type LifecycleSystem struct{}

func NewLifecycleSystem() *LifecycleSystem { return &LifecycleSystem{} }

func (s *LifecycleSystem) Update(world *ecs.World, deltaTime float32) {
	for _, entity := range world.Entities() {

		if lifetime := ecs.GetComponent[components.Lifetime](entity); lifetime != nil {
			lifetime.Remaining -= deltaTime

			if lifetime.Remaining < 0 {
				world.EnqueueRemoval(entity)
			}
		}
	}
}
