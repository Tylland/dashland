package game

import rl "github.com/gen2brain/raylib-go/raylib"

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

type BlockPosition struct {
	X int
	Y int
}

func (bp BlockPosition) Offset(deltaX int, deltaY int) BlockPosition {
	return BlockPosition{X: bp.X + deltaX, Y: bp.Y + deltaY}
}

func (bp BlockPosition) Add(position BlockPosition) BlockPosition {
	return BlockPosition{X: bp.X + position.X, Y: bp.Y + position.Y}
}

func (bp BlockPosition) Subtract(position BlockPosition) BlockPosition {
	return BlockPosition{X: bp.X - position.X, Y: bp.Y - position.Y}
}

func (bp BlockPosition) IsSame(other BlockPosition) bool {
	return bp.X == other.X && bp.Y == other.Y
}

type BlockRectangle struct {
	center rl.Vector2
	size   rl.Vector2
}

func (br *BlockRectangle) Rectangle() rl.Rectangle {
	return rl.Rectangle{X: br.center.X - br.size.X/2, Y: br.center.Y - br.size.Y/2, Width: br.size.X, Height: br.size.Y}
}

type Block struct {
	world     *world
	blockType BlockType
	position  BlockPosition
	behavior  BlockBehavior
	corners   [4]uint8
}

func NewBlock(world *world, blockType BlockType, x int, y int) *Block {
	switch blockType {
	case Diamond:
		return &Block{world: world, blockType: blockType, position: BlockPosition{X: x, Y: y}, behavior: CanFall | Collectable}
	case Boulder:
		return &Block{world: world, blockType: blockType, position: BlockPosition{X: x, Y: y}, behavior: CanFall | Obstacle | Pushable}
	default:
		return &Block{world: world, blockType: blockType, position: BlockPosition{X: x, Y: y}, behavior: NoBehavior}
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

func (b *Block) cornerIndex(voids [9]bool, corner CornerPosition) uint8 {

	if corner == CornerOverLeft && voids[MiddleLeft] && voids[OverLeft] && voids[OverCenter] {
		return 1
	}

	if corner == CornerOverRight && voids[OverCenter] && voids[OverRight] && voids[MiddleRight] {
		return 2
	}

	if corner == CornerUnderRight && voids[MiddleRight] && voids[UnderRight] && voids[UnderCenter] {
		return 3
	}

	if corner == CornerUnderLeft && voids[UnderCenter] && voids[UnderLeft] && voids[MiddleLeft] {
		return 4
	}

	return 0
}
func (b *Block) update(deltaTime float32) {

	if b.blockType == Soil {
		voids := b.world.CheckNeighbourTypes(Void, b.position)

		b.corners[CornerOverLeft] = b.cornerIndex(voids, CornerOverLeft)
		b.corners[CornerOverRight] = b.cornerIndex(voids, CornerOverRight)
		b.corners[CornerUnderRight] = b.cornerIndex(voids, CornerUnderRight)
		b.corners[CornerUnderLeft] = b.cornerIndex(voids, CornerUnderLeft)
	}

	if b.blockType == Void {
		soils := b.world.CheckNeighbourTypes(Soil, b.position)

		b.corners[CornerOverLeft] = b.cornerIndex(soils, CornerOverLeft)
		b.corners[CornerOverRight] = b.cornerIndex(soils, CornerOverRight)
		b.corners[CornerUnderRight] = b.cornerIndex(soils, CornerUnderRight)
		b.corners[CornerUnderLeft] = b.cornerIndex(soils, CornerUnderLeft)
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
