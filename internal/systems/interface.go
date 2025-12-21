package systems

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/tylland/dashland/internal/common"
	"github.com/tylland/dashland/internal/ecs"
	"github.com/tylland/dashland/internal/game"
)

type EntityLister interface {
	Entities() []*ecs.Entity
}

type PositionResolver interface {
	GetPosition(common.BlockPosition) rl.Vector2
}

type EntityResolver interface {
	GetEntity(common.BlockPosition) *ecs.Entity
}

type BlockResolver interface {
	GetPosition(common.BlockPosition) rl.Vector2
	VisitBlock(common.BlockPosition)
	GetBlockAtPosition(common.BlockPosition) (*game.Block, bool)
	CheckBlockAtPosition(blockType game.BlockType, position common.BlockPosition) bool
}
