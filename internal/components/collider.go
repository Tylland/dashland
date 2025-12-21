package components

import "github.com/tylland/dashland/internal/common"

type ColliderComponent struct {
	Layer common.CollisionLayer
	Mask  common.CollisionLayer
}

func NewColliderComponent(layer common.CollisionLayer, mask common.CollisionLayer) *ColliderComponent {
	return &ColliderComponent{Layer: layer, Mask: mask}
}

func (c *ColliderComponent) CollidesWith(other *ColliderComponent) bool {
	return c.Mask&other.Layer != 0
}

func (c *ColliderComponent) Collides(layer common.CollisionLayer) bool {
	return c.Mask&layer != 0
}
