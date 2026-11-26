package components

type WallWalker struct {
	HasWall   bool
	Clockwise bool
}

func NewWallWalker(clockwise bool) *WallWalker {
	return &WallWalker{Clockwise: clockwise}
}
