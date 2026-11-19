package components

import rl "github.com/gen2brain/raylib-go/raylib"

type IconComponent struct {
	Texture *rl.Texture2D
	Source  rl.Rectangle
}

func NewIconComponent(texture *rl.Texture2D, source rl.Rectangle) *IconComponent {
	return &IconComponent{Texture: texture, Source: source}
}
