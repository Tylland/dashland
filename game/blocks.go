package game

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
	position  BlockPosition
	behavior  BlockBehavior
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

func (b *Block) update(deltaTime float32) {
	// if b.behavior&CanFall == CanFall {
	// 	under := b.position.Offset(0, 1)

	// 	if !b.world.CheckBlockAtPosition(Soil, under) && !b.world.checkPlayerAtPosition(under) {

	// 		if !b.world.checkPositionOccupied(under) {
	// 			b.world.SwapBlock(b, under)
	// 		}

	// 		right := b.position.Offset(1, 0)
	// 		rightUnder := b.position.Offset(1, 1)

	// 		if !b.world.checkPositionOccupied(right) && !b.world.checkPositionOccupied(rightUnder) {
	// 			b.world.SwapBlock(b, rightUnder)
	// 		}

	// 		left := b.position.Offset(-1, 0)
	// 		leftUnder := b.position.Offset(-1, 1)

	// 		if !b.world.checkPositionOccupied(left) && !b.world.checkPositionOccupied(leftUnder) {
	// 			b.world.SwapBlock(b, leftUnder)
	// 		}
	// 	}

	// }
}

func (b *Block) ObstacleForPlayer(player *Player) bool {
	return b.blockType == Bedrock || b.blockType == Boulder
}

type IBlock interface {
	ObstacleForPlayer(player *Player) bool
}

type UnknownBlock struct {
}

func (ub UnknownBlock) ObstacleForPlayer(player *Player) bool {
	return false
}

type BoulderBlock struct {
}

func (b BoulderBlock) ObstacleForPlayer(player *Player) bool {
	return !player.pickaxe
}
