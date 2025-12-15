package common

import rl "github.com/gen2brain/raylib-go/raylib"

type Sprite struct {
	Texture    *rl.Texture2D
	Source     rl.Rectangle
	SourceBase rl.Vector2
	Width      float32
	Height     float32
	// Animation properties
	Frame      uint    // Current animation frame
	FrameCount uint    // Total frames in animation
	FrameSpeed float32 // Animation speed (frames per second)
	FrameTimer float32 // Time accumulator for animation

}

//NewSprite creates a new sprite
func NewSprite(texture *rl.Texture2D, width float32, height float32, sourceX float32, sourceY float32, frameCount uint, frame uint) *Sprite {
	sprite := &Sprite{Texture: texture, Width: width, Height: height, SourceBase: rl.Vector2{X: sourceX, Y: sourceY}, FrameCount: frameCount, FrameSpeed: 10}
	sprite.UpdateFrame(frame)
	return sprite

}

func (s *Sprite) UpdateFrame(frame uint) {
	s.Frame = frame

	s.Source = rl.Rectangle{
		X:      s.SourceBase.X + float32(s.Frame)*s.Width,
		Y:      s.SourceBase.Y,
		Width:  s.Width,
		Height: s.Height,
	}

}
