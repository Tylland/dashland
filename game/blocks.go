package game

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/tylland/dashland/game/core"
)

type BlockType uint16

const (
	Unknown BlockType = iota
	Bedrock
	Void
	Soil
	Diamond
	Boulder
	All = BlockType(0xFFFF)
)

type BlockBehavior uint16

const (
	NoBehavior               = BlockBehavior(0)
	Obstacle   BlockBehavior = 1 << iota
	CanFall
	Collectable
	Pushable
)

type Block struct {
	world     *world
	blockType BlockType
	position  core.BlockPosition
	behavior  BlockBehavior
	corners   [4]uint8
}

func NewBlock(world *world, blockType BlockType, x int, y int) *Block {
	switch blockType {
	case Diamond:
		return &Block{world: world, blockType: blockType, position: core.BlockPosition{X: x, Y: y}, behavior: CanFall | Collectable}
	case Boulder:
		return &Block{world: world, blockType: blockType, position: core.BlockPosition{X: x, Y: y}, behavior: CanFall | Obstacle | Pushable}
	default:
		return &Block{world: world, blockType: blockType, position: core.BlockPosition{X: x, Y: y}, behavior: NoBehavior}
	}

}

const (
	BlockLeft   = -1
	BlockCenter = 0
	BlockRight  = 1
	BlockOver   = -1
	BlockMiddle = 0
	BlockUnder  = 1
)

const (
	OverLeft = iota
	OverCenter
	OverRight
	MiddleLeft
	MiddleCenter
	MiddleRight
	UnderLeft
	UnderCenter
	UnderRight
)

type CornerPosition uint16

const (
	CornerOverLeft CornerPosition = iota
	CornerOverRight
	CornerUnderLeft
	CornerUnderRight
)

func (b *Block) cornerIndex(neighbors [9]bool, corner CornerPosition, strict bool) uint8 {

	if corner == CornerOverLeft && neighbors[MiddleLeft] && (neighbors[OverLeft] || !strict) && neighbors[OverCenter] {
		return 1
	}

	if corner == CornerOverRight && neighbors[OverCenter] && (neighbors[OverRight] || !strict) && neighbors[MiddleRight] {
		return 2
	}

	if corner == CornerUnderRight && neighbors[MiddleRight] && (neighbors[UnderRight] || !strict) && neighbors[UnderCenter] {
		return 3
	}

	if corner == CornerUnderLeft && neighbors[UnderCenter] && (neighbors[UnderLeft] || !strict) && neighbors[MiddleLeft] {
		return 4
	}

	return 0
}

func (b *Block) update(deltaTime float32) {

	if b.blockType == Soil {
		voids := b.world.CheckNeighbourTypes(Void, b.position)
		strict := true

		b.corners[CornerOverLeft] = b.cornerIndex(voids, CornerOverLeft, strict)
		b.corners[CornerOverRight] = b.cornerIndex(voids, CornerOverRight, strict)
		b.corners[CornerUnderRight] = b.cornerIndex(voids, CornerUnderRight, strict)
		b.corners[CornerUnderLeft] = b.cornerIndex(voids, CornerUnderLeft, strict)
	}

	if b.blockType == Void {
		soils := b.world.CheckNeighbourTypes(Soil, b.position)
		strict := false

		b.corners[CornerOverLeft] = b.cornerIndex(soils, CornerOverLeft, strict)
		b.corners[CornerOverRight] = b.cornerIndex(soils, CornerOverRight, strict)
		b.corners[CornerUnderRight] = b.cornerIndex(soils, CornerUnderRight, strict)
		b.corners[CornerUnderLeft] = b.cornerIndex(soils, CornerUnderLeft, strict)
	}
}

func (b *Block) renderCorners(row float32) {
	cornerWidth := b.world.blockWidth / 2
	cornerHeight := b.world.blockHeight / 2

	rl.DrawTextureRec(b.world.groundCorners, rl.NewRectangle(float32(b.corners[CornerOverLeft])*cornerWidth, row*cornerHeight, cornerWidth, cornerHeight), rl.NewVector2(float32(b.position.X)*b.world.blockWidth, float32(b.position.Y)*b.world.blockHeight), rl.White)
	rl.DrawTextureRec(b.world.groundCorners, rl.NewRectangle(float32(b.corners[CornerOverRight])*cornerWidth, row*cornerHeight, cornerWidth, cornerHeight), rl.NewVector2(float32(b.position.X)*b.world.blockWidth+cornerWidth, float32(b.position.Y)*b.world.blockHeight), rl.White)
	rl.DrawTextureRec(b.world.groundCorners, rl.NewRectangle(float32(b.corners[CornerUnderRight])*cornerWidth, row*cornerHeight, cornerWidth, cornerHeight), rl.NewVector2(float32(b.position.X)*b.world.blockWidth+cornerWidth, float32(b.position.Y)*b.world.blockHeight+cornerHeight), rl.White)
	rl.DrawTextureRec(b.world.groundCorners, rl.NewRectangle(float32(b.corners[CornerUnderLeft])*cornerWidth, row*cornerHeight, cornerWidth, cornerHeight), rl.NewVector2(float32(b.position.X)*b.world.blockWidth, float32(b.position.Y)*b.world.blockHeight+cornerHeight), rl.White)
}

func (b *Block) render() {
	if b.blockType == Soil {
		b.renderCorners(1)
	} else if b.blockType == Void {
		b.renderCorners(0)
	} else {
		rl.DrawTextureRec(b.world.blockTextures, rl.NewRectangle(float32(b.blockType)*b.world.blockWidth, 0, b.world.blockWidth, b.world.blockHeight), rl.NewVector2(float32(b.position.X)*b.world.blockWidth, float32(b.position.Y)*b.world.blockHeight), rl.White)
	}
}

func (b *Block) IsObstacleForPlayer(player *Player) bool {
	return b.blockType == Bedrock || b.blockType == Boulder
}

func (b *Block) Rectangle() rl.Rectangle {
	return rl.Rectangle{X: float32(b.position.X) * b.world.blockWidth, Y: float32(b.position.Y) * b.world.blockHeight, Width: b.world.blockWidth, Height: b.world.blockHeight}
}

type IBlock interface {
	IsObstacleForPlayer(player *Player) bool
}

type UnknownBlock struct {
}

func (ub UnknownBlock) IsObstacleForPlayer(player *Player) bool {
	return false
}

type BoulderBlock struct {
}

func (b BoulderBlock) IsObstacleForPlayer(player *Player) bool {
	return !player.pickaxe
}
