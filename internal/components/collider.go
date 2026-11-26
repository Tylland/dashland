package components

import (
	"github.com/tylland/dashland/internal/collider"
	"github.com/tylland/dashland/internal/ecs"
	"github.com/tylland/dashland/internal/game/types"
)

type Collider struct {
	Layer    collider.CollisionLayer
	Obstruct collider.CollisionLayer
	Interact collider.CollisionLayer
}

func NewCollider(layer collider.CollisionLayer, obstructs collider.CollisionLayer, interact collider.CollisionLayer) *Collider {
	return &Collider{Layer: layer, Obstruct: obstructs, Interact: interact}
}

func (c *Collider) Collide(other *Collider) (bool, bool) {
	obstructs := c.Obstruct&other.Layer != collider.LayerNone || other.Obstruct&c.Layer != collider.LayerNone
	interacts := c.Interact&other.Layer != collider.LayerNone || other.Interact&c.Layer != collider.LayerNone

	return obstructs, interacts
}

func (c *Collider) WriteTo(w ecs.Writer) error {
	w.StartSection(Section(types.ComponentTypeCollider))
	w.WriteEnum("Layer", c.Layer, true, "is")
	w.WriteEnum("Obstruct", c.Obstruct, true, "obstructs")
	w.WriteEnum("Interact", c.Interact, true, "interacts with")

	return nil
}

func (c *Collider) ReadFrom(r ecs.Reader) error {
	r.ConsumeSection(types.ComponentTypeCollider.Name)
	r.ConsumeLine()
	r.ConsumeLine()
	r.ConsumeLine()

	return nil
}
