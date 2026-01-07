package game

type MapSize struct {
	Width       int
	Height      int
	BlockWidth  float32
	BlockHeight float32
}

func NewMapSize(width, height, blockWidth, blockHeight int) MapSize {
	return MapSize{Width: width, Height: height, BlockWidth: float32(blockWidth), BlockHeight: float32(blockHeight)}
}
