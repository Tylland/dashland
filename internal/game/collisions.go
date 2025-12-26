package game

import "github.com/tylland/dashland/internal/common"

const (
	LayerNone common.CollisionLayer = 0
	LayerWall common.CollisionLayer = 1 << iota
	LayerPlayer
	LayerEnemy
	LayerBedrock
	LayerGround
	LayerProjectile
	LayerItem
	LayerCollectable
	// LayerWater
	// LayerTrap
	//	LayerTrigger
	LayerAll common.CollisionLayer = 0xFFFF
)
