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
	w.WriteLabel(ecs.ComponentName(c))
	w.WriteFloat(c.Width, true)
	w.WriteFloat(c.Height, true)
	w.WriteInt(int(c.Texture.ID), true)

	return nil
}
