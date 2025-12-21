package game

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/tylland/dashland/internal/characteristics"
	"github.com/tylland/dashland/internal/common"
	"github.com/tylland/dashland/internal/components"
)

type BlockType uint16

const (
	Unknown BlockType = iota
	Bedrock
	Void
	Soil
	Wall
	All = BlockType(0xFFFF)
)

type Block struct {
	bm *BlockMap
	gm *EntityMap

	BlockType BlockType
	Position  common.BlockPosition
	Character characteristics.Characteristics
	Collider  *components.ColliderComponent
	Corners   [4]uint8
}

func NewBlockWithCharacteristics(bm *BlockMap, gm *EntityMap, blockType BlockType, position common.BlockPosition, character characteristics.Characteristics, layer common.CollisionLayer, mask common.CollisionLayer) *Block {
	return &Block{bm: bm, gm: gm, BlockType: blockType, Position: position, Collider: components.NewColliderComponent(layer, mask)}
}

func NewBlock(bm *BlockMap, gm *EntityMap, blockType BlockType, position common.BlockPosition) *Block {
	switch blockType {
	case Void:
		return NewBlockWithCharacteristics(bm, gm, blockType, position, characteristics.Void, LayerGround, LayerNone)
	case Soil:
		return NewBlockWithCharacteristics(bm, gm, blockType, position, characteristics.RollOff|characteristics.EnemyObstacle, LayerGround, LayerEnemy)
	case Wall:
		return NewBlockWithCharacteristics(bm, gm, blockType, position, characteristics.EnemyObstacle|characteristics.PlayerObstacle, LayerWall, LayerAll)
	case Bedrock:
		return NewBlockWithCharacteristics(bm, gm, blockType, position, characteristics.EnemyObstacle|characteristics.PlayerObstacle, LayerWall, LayerAll)
	default:
		return NewBlockWithCharacteristics(bm, gm, blockType, position, characteristics.None, LayerGround, LayerNone)
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

func (b *Block) HasCharacteristic(characteristics characteristics.Characteristics) bool {
	return b.Character&characteristics == characteristics
}

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

func (b *Block) Update(_ float32) {

	if b.BlockType == Soil {
		voids := b.bm.CheckNeighbourTypes(Void, b.Position)
		strict := true

		b.Corners[CornerOverLeft] = b.cornerIndex(voids, CornerOverLeft, strict)
		b.Corners[CornerOverRight] = b.cornerIndex(voids, CornerOverRight, strict)
		b.Corners[CornerUnderRight] = b.cornerIndex(voids, CornerUnderRight, strict)
		b.Corners[CornerUnderLeft] = b.cornerIndex(voids, CornerUnderLeft, strict)
	}

	if b.BlockType == Void {
		soils := b.bm.CheckNeighbourTypes(Soil, b.Position)
		strict := false

		b.Corners[CornerOverLeft] = b.cornerIndex(soils, CornerOverLeft, strict)
		b.Corners[CornerOverRight] = b.cornerIndex(soils, CornerOverRight, strict)
		b.Corners[CornerUnderRight] = b.cornerIndex(soils, CornerUnderRight, strict)
		b.Corners[CornerUnderLeft] = b.cornerIndex(soils, CornerUnderLeft, strict)
	}
}

func (b *Block) renderCorners(row float32) {
	cornerWidth := b.gm.BlockWidth / 2
	cornerHeight := b.gm.BlockHeight / 2

	rl.DrawTextureRec(*b.gm.groundCorners, rl.NewRectangle(float32(b.Corners[CornerOverLeft])*cornerWidth, row*cornerHeight, cornerWidth, cornerHeight), rl.NewVector2(float32(b.Position.X)*b.gm.BlockWidth, float32(b.Position.Y)*b.gm.BlockHeight), rl.White)
	rl.DrawTextureRec(*b.gm.groundCorners, rl.NewRectangle(float32(b.Corners[CornerOverRight])*cornerWidth, row*cornerHeight, cornerWidth, cornerHeight), rl.NewVector2(float32(b.Position.X)*b.gm.BlockWidth+cornerWidth, float32(b.Position.Y)*b.gm.BlockHeight), rl.White)
	rl.DrawTextureRec(*b.gm.groundCorners, rl.NewRectangle(float32(b.Corners[CornerUnderRight])*cornerWidth, row*cornerHeight, cornerWidth, cornerHeight), rl.NewVector2(float32(b.Position.X)*b.gm.BlockWidth+cornerWidth, float32(b.Position.Y)*b.gm.BlockHeight+cornerHeight), rl.White)
	rl.DrawTextureRec(*b.gm.groundCorners, rl.NewRectangle(float32(b.Corners[CornerUnderLeft])*cornerWidth, row*cornerHeight, cornerWidth, cornerHeight), rl.NewVector2(float32(b.Position.X)*b.gm.BlockWidth, float32(b.Position.Y)*b.gm.BlockHeight+cornerHeight), rl.White)
}

func (b *Block) Render() {
	switch b.BlockType {
	case Soil:
		b.renderCorners(1)
	case Void:
		b.renderCorners(0)
	default:
		rl.DrawTextureRec(*b.bm.blockTextures, rl.NewRectangle(float32(b.BlockType)*b.bm.BlockWidth, 0, b.bm.BlockWidth, b.bm.BlockHeight), rl.NewVector2(float32(b.Position.X)*b.bm.BlockWidth, float32(b.Position.Y)*b.bm.BlockHeight), rl.White)
	}
}

func (b *Block) Rectangle() rl.Rectangle {
	return rl.Rectangle{X: float32(b.Position.X) * b.bm.BlockWidth, Y: float32(b.Position.Y) * b.bm.BlockHeight, Width: b.bm.BlockWidth, Height: b.bm.BlockHeight}
}
