package components

import "github.com/tylland/dashland/game/core"

type WallWalkerComponent struct {
	HasWall          bool
	DefaultDirection core.BlockVector
}
