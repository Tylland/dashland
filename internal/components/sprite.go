package components

import "github.com/tylland/dashland/internal/common"

type SpriteComponent struct {
	*common.Sprite
}

func NewSpriteComponent(sprite *common.Sprite) *SpriteComponent {
	return &SpriteComponent{Sprite: sprite}
}
