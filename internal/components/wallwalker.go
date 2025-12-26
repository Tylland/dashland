package components

import "github.com/tylland/dashland/internal/common"

type WallWalkerComponent struct {
	HasWall bool
}

func NewWallWalkerComponent(defaultDirection common.BlockVector) *WallWalkerComponent {
	return &WallWalkerComponent{}
}
