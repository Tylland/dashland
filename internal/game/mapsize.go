package game

import "github.com/tylland/dashland/internal/ecs"

type MapSize struct {
	Width       int
	Height      int
	BlockWidth  float32
	BlockHeight float32
}

func NewMapSize(width, height, blockWidth, blockHeight int) MapSize {
	return MapSize{Width: width, Height: height, BlockWidth: float32(blockWidth), BlockHeight: float32(blockHeight)}
}

func (s *MapSize) Count() int {
	return s.Height * s.Width
}

func (s *MapSize) WriteTo(w ecs.Writer) error {
	w.WriteLabel("size")
	w.WriteInt(s.Width, false)
	w.WriteInt(s.Height, false)
	w.WriteFloat(s.BlockWidth, false)
	w.WriteFloat(s.BlockHeight, false)

	return nil
}

func (s *MapSize) ReadFrom(r ecs.Reader) error {
	r.ConsumeLabel("size")

	ecs.WithChangedValue(r.ReadInt, func(width int) {
		s.Width = width
	})

	ecs.WithChangedValue(r.ReadInt, func(height int) {
		s.Height = height
	})

	ecs.WithChangedValue(r.ReadFloat, func(blockWidth float32) {
		s.BlockWidth = blockWidth
	})

	ecs.WithChangedValue(r.ReadFloat, func(blockHeight float32) {
		s.BlockHeight = blockHeight
	})

	return nil
}
