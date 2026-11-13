package components

import (
	"github.com/tylland/dashland/internal/collider"
	"github.com/tylland/dashland/internal/ecs"
)

type ColliderComponent struct {
	Layer        collider.CollisionLayer
	ObstructMask collider.CollisionLayer
	InteractMask collider.CollisionLayer
}

func NewColliderComponent(layer collider.CollisionLayer, block collider.CollisionLayer, overlap collider.CollisionLayer) *ColliderComponent {
	return &ColliderComponent{Layer: layer, ObstructMask: block, InteractMask: overlap}
}

func (c *ColliderComponent) Collide(other *ColliderComponent) (bool, bool) {
	obstructs := c.ObstructMask&other.Layer != 0 || other.ObstructMask&c.Layer != 0
	interacts := c.InteractMask&other.Layer != 0 || other.InteractMask&c.Layer != 0

	return obstructs, interacts
}

func (c *ColliderComponent) WriteTo(w ecs.Writer) error {
	w.WriteLabel(ecs.ComponentName(c))
	w.WriteLabel("layer")
	w.WriteEnumValue(c.Layer, true)

	w.WriteLabel("obstructs")
	w.WriteEnumValue(c.ObstructMask, true)

	w.WriteLabel("interacts with")
	w.WriteEnumValue(c.InteractMask, true)

	return nil
}
