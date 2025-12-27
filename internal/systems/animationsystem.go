package systems

import (
    "github.com/tylland/dashland/internal/components"
    "github.com/tylland/dashland/internal/ecs"
)

// AnimationSystem applies AnimationComponent selections to Sprite components.
type AnimationSystem struct{}

func NewAnimationSystem() *AnimationSystem { return &AnimationSystem{} }

func (s *AnimationSystem) Update(world *ecs.World, deltaTime float32) {
    for _, entity := range world.Entities() {
        animation := ecs.GetComponent[components.AnimationComponent](entity)
        sprite := ecs.GetComponent[components.SpriteComponent](entity)

        if animation == nil || sprite == nil {
            continue
        }

        animation.ApplyAnimation(sprite.Sprite)
        animation.Advance(deltaTime, sprite.Sprite)
    }
}
