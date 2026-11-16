package components

import (
	"github.com/tylland/dashland/internal/collider"
	"github.com/tylland/dashland/internal/ecs"
)

type ColliderComponent struct {
	Layer    collider.CollisionLayer
	Obstruct collider.CollisionLayer
	Interact collider.CollisionLayer
}

func NewColliderComponent(layer collider.CollisionLayer, obstructs collider.CollisionLayer, interact collider.CollisionLayer) *ColliderComponent {
	return &ColliderComponent{Layer: layer, Obstruct: obstructs, Interact: interact}
}

func (c *ColliderComponent) Collide(other *ColliderComponent) (bool, bool) {
	obstructs := c.Obstruct&other.Layer != collider.LayerNone || other.Obstruct&c.Layer != collider.LayerNone
	interacts := c.Interact&other.Layer != collider.LayerNone || other.Interact&c.Layer != collider.LayerNone

	return obstructs, interacts
}

func (c *ColliderComponent) WriteTo(w ecs.Writer) error {
	w.WriteLabel(ecs.ComponentName(c))
	w.WriteEnum(c.Layer, true, "is")
	w.WriteEnum(c.Obstruct, true, "obstructs")
	w.WriteEnum(c.Interact, true, "interacts with")

	return nil
}

func (c *ColliderComponent) ReadFrom(r ecs.Reader) error {
	r.ConsumeLabel(ecs.ComponentName(c))
	r.ConsumeLine()
	r.ConsumeLine()
	r.ConsumeLine()

	return nil
}
