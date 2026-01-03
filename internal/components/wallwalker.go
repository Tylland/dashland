package components

type WallWalkerComponent struct {
	HasWall   bool
	Clockwise bool
}

func NewWallWalkerComponent(clockwise bool) *WallWalkerComponent {
	return &WallWalkerComponent{Clockwise: clockwise}
}
