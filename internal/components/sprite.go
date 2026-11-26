package components

import (
	"github.com/tylland/dashland/internal/common"
	"github.com/tylland/dashland/internal/ecs"
	"github.com/tylland/dashland/internal/game/types"
)

type Sprite struct {
	*common.Sprite
}

func NewSprite(sprite *common.Sprite) *Sprite {
	return &Sprite{Sprite: sprite}
}

func (c *Sprite) WriteTo(w ecs.Writer) error {
	w.StartSection(Section(types.ComponentTypeSprite))

	w.WriteFloat("Width", c.Width, true)
	w.WriteFloat("Height", c.Height, true)
	//w.WriteInt(int(c.Texture.ID), true)

	return nil
}

func (c *Sprite) ReadFrom(r ecs.Reader) error {
	r.ConsumeSection(types.ComponentTypeSprite.Name)
	r.ConsumeLine()
	r.ConsumeLine()
	//r.ConsumeLine()

	return nil
}
