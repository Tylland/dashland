package components

import "github.com/tylland/dashland/internal/common"

type ColliderComponent struct {
	Layer       common.CollisionLayer
	BlockMask   common.CollisionLayer
	CollideMask common.CollisionLayer
}

func NewColliderComponent(layer common.CollisionLayer, block common.CollisionLayer, overlap common.CollisionLayer) *ColliderComponent {
	return &ColliderComponent{Layer: layer, BlockMask: block, CollideMask: overlap}
}

func (c *ColliderComponent) Result(other *ColliderComponent) (bool, bool) {
	block := c.BlockMask&other.Layer != 0 || other.BlockMask&c.Layer != 0
	collide := c.CollideMask&other.Layer != 0 || other.CollideMask&c.Layer != 0

	return block, collide
}
