package components

import "github.com/tylland/dashland/game/core"

type SpriteComponent struct {
	*core.Sprite
}

func NewSpriteComponent(sprite *core.Sprite) *SpriteComponent {
	return &SpriteComponent{Sprite: sprite}
}
