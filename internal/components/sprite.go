package components

import (
	"github.com/tylland/dashland/internal/common"
	"github.com/tylland/dashland/internal/ecs"
)

type SpriteComponent struct {
	*common.Sprite
}

func NewSpriteComponent(sprite *common.Sprite) *SpriteComponent {
	return &SpriteComponent{Sprite: sprite}
}

func (c *SpriteComponent) WriteTo(w ecs.Writer) error {
	w.StartSection(ecs.ComponentName(c))
	w.WriteFloat("Width", c.Width, true)
	w.WriteFloat("Height", c.Height, true)
	//w.WriteInt(int(c.Texture.ID), true)

	return nil
}

func (c *SpriteComponent) ReadFrom(r ecs.Reader) error {
	r.ConsumeLabel(ecs.ComponentName(c))
	r.ConsumeLine()
	r.ConsumeLine()
	//r.ConsumeLine()

	return nil
}
