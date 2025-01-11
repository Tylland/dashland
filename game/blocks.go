package game

// import (
// 	"fmt"
// 	"path"

// 	rl "github.com/gen2brain/raylib-go/raylib"
// 	"github.com/lafriks/go-tiled"
// )

type BlockType uint32

const (
	Unknown BlockType = iota
	Bedrock
	Void
	Soil
	Diamond
	Boulder
)

type BlockFlags uint8

const (
	Obstacle BlockFlags = 1 << iota
	CanFall
	Collectable
)

type Block struct {
	blockMap  *BlockMap
	blockType BlockType
	position  BlockPosition
	flags     BlockFlags
}

func NewBlock(blockMap *BlockMap, blockType BlockType, position BlockPosition) Block {
	return Block{blockMap: blockMap, blockType: blockType, position: position}
}

func (b *Block) MoveTo(pos BlockPosition) {
	b.position = pos
	b.blockMap.blocks[pos.Y*b.blockMap.width+pos.X] = *b
}

func (b *Block) Update(deltaTime float32) {
	if b.flags&CanFall == CanFall {
		if !b.blockMap.CheckTypeAtPosition(Soil, b.position.Offset(0, 1)) {

			if b.blockMap.CheckTypeAtPosition(Void, b.position.Offset(0, 1)) {
				b.blockMap.SwapBlock(b, b.position.Offset(0, 1))
			}

			if b.blockMap.CheckTypeAtPosition(Void, b.position.Offset(1, 0)) &&
				b.blockMap.CheckTypeAtPosition(Void, b.position.Offset(1, 1)) {
				b.blockMap.SwapBlock(b, b.position.Offset(1, 1))
			}

			if b.blockMap.CheckTypeAtPosition(Void, b.position.Offset(-1, 0)) &&
				b.blockMap.CheckTypeAtPosition(Void, b.position.Offset(-1, 1)) {
				b.blockMap.SwapBlock(b, b.position.Offset(-1, 1))
			}
		}

	}
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
	return !player.Pickaxe
}
