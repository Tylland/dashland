package common

import rl "github.com/gen2brain/raylib-go/raylib"

type Sprite struct {
	Texture    *rl.Texture2D
	Source     rl.Rectangle
	SourceBase rl.Vector2
	Width      float32
	Height     float32
	// Sprite is a render resource; animation timing moved to AnimationComponent

}

//NewSprite creates a new sprite
func NewSprite(texture *rl.Texture2D, width float32, height float32, sourceX float32, sourceY float32, frame uint) *Sprite {
	sprite := &Sprite{Texture: texture, Width: width, Height: height, SourceBase: rl.Vector2{X: sourceX, Y: sourceY}}
	sprite.UpdateFrame(frame)
	return sprite

}

func (s *Sprite) UpdateFrame(frame uint) {
	s.Source = rl.Rectangle{
		X:      s.SourceBase.X + float32(frame)*s.Width,
		Y:      s.SourceBase.Y,
		Width:  s.Width,
		Height: s.Height,
	}

}
