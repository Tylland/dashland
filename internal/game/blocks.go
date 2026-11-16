package game

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/tylland/dashland/internal/characteristics"
	"github.com/tylland/dashland/internal/collider"
	"github.com/tylland/dashland/internal/common"
	"github.com/tylland/dashland/internal/components"
	"github.com/tylland/dashland/internal/ecs"
)

type BlockType uint32

const (
	Unknown BlockType = iota
	Bedrock
	Void
	Soil
	Wall
	All = BlockType(0xFFFF)
)

func (bt BlockType) String() string {
	switch bt {
	case Bedrock:
		return "bedrock"
	case Void:
		return "void"
	case Soil:
		return "soil"
	case Wall:
		return "wall"
	default:
		return "unknown"
	}
}

func ToBlockName(id ecs.EntityType) string {
	switch id {
	case EntityBedrock:
		return "bedrock"
	case EntityVoid:
		return "void"
	case EntitySoil:
		return "soil"
	case EntityWall:
		return "wall"
	default:
		return "unknown"
	}
}

type Block struct {
	*ecs.EntityHeader
	bm *BlockMap
	gm *EntityMap

	Position  common.BlockPosition
	Character characteristics.Characteristics
	Collider  *components.ColliderComponent
	Corners   [4]uint8
}

func NewBlockWithCharacteristics(header *ecs.EntityHeader, bm *BlockMap, gm *EntityMap, position common.BlockPosition, character characteristics.Characteristics, layer collider.CollisionLayer, block collider.CollisionLayer, overlap collider.CollisionLayer) *Block {
	return &Block{EntityHeader: header, bm: bm, gm: gm, Position: position, Character: character, Collider: components.NewColliderComponent(layer, block, overlap)}
}

func NewBlock(header *ecs.EntityHeader, bm *BlockMap, gm *EntityMap, position common.BlockPosition) *Block {
	switch header.Type {
	case EntityVoid:
		return NewBlockWithCharacteristics(header, bm, gm, position, characteristics.Destructable, collider.LayerNone, collider.LayerNone, collider.LayerNone)
	case EntitySoil:
		return NewBlockWithCharacteristics(header, bm, gm, position, characteristics.CanHoldGravity|characteristics.Destructable, collider.LayerGround, collider.LayerEnemy, collider.LayerAll)
	case EntityWall:
		return NewBlockWithCharacteristics(header, bm, gm, position, characteristics.CanHoldGravity|characteristics.GravityRollOff|characteristics.Destructable, collider.LayerWall, collider.LayerAll, collider.LayerNone)
	case EntityBedrock:
		return NewBlockWithCharacteristics(header, bm, gm, position, characteristics.CanHoldGravity, collider.LayerBedrock, collider.LayerAll, collider.LayerNone)
	default:
		return NewBlockWithCharacteristics(header, bm, gm, position, characteristics.None, collider.LayerGround, collider.LayerNone, collider.LayerNone)
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

	if b.Type == EntitySoil {
		voids := b.bm.CheckNeighbourTypes(EntityVoid, b.Position)
		strict := true

		b.Corners[CornerOverLeft] = b.cornerIndex(voids, CornerOverLeft, strict)
		b.Corners[CornerOverRight] = b.cornerIndex(voids, CornerOverRight, strict)
		b.Corners[CornerUnderRight] = b.cornerIndex(voids, CornerUnderRight, strict)
		b.Corners[CornerUnderLeft] = b.cornerIndex(voids, CornerUnderLeft, strict)
	}

	if b.Type == EntityVoid {
		soils := b.bm.CheckNeighbourTypes(EntitySoil, b.Position)
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
	switch b.Type {
	case EntitySoil:
		b.renderCorners(1)
	case EntityVoid:
		b.renderCorners(0)
	default:
		rl.DrawTextureRec(*b.bm.blockTextures, rl.NewRectangle(float32(b.Type)*b.bm.BlockWidth, 0, b.bm.BlockWidth, b.bm.BlockHeight), rl.NewVector2(float32(b.Position.X)*b.bm.BlockWidth, float32(b.Position.Y)*b.bm.BlockHeight), rl.White)
	}
}

func (b *Block) Rectangle() rl.Rectangle {
	return rl.Rectangle{X: float32(b.Position.X) * b.bm.BlockWidth, Y: float32(b.Position.Y) * b.bm.BlockHeight, Width: b.bm.BlockWidth, Height: b.bm.BlockHeight}
}

func (b *Block) WriteTo(w ecs.Writer) error {
	b.EntityHeader.WriteTo(w)

	w.WriteLabel("position")
	w.WriteInt(b.Position.X, false)
	w.WriteInt(b.Position.Y, false)

	w.WriteLabel("characteristics")
	w.WriteEnum(b.Character, true, "")

	err := b.Collider.WriteTo(w)

	if err != nil {
		return err
	}

	w.WriteLabel("corners")
	w.WriteInt(int(b.Corners[0]), b.Corners[0] == 0)
	w.WriteInt(int(b.Corners[1]), b.Corners[1] == 0)
	w.WriteInt(int(b.Corners[2]), b.Corners[2] == 0)
	w.WriteInt(int(b.Corners[3]), b.Corners[3] == 0)

	return nil
}

func (b *Block) ReadFrom(r ecs.Reader) error {
	r.ConsumeLabel("position")

	ecs.WithChangedValue(r.ReadInt, func(x int) {
		b.Position.X = x
	})

	ecs.WithChangedValue(r.ReadInt, func(y int) {
		b.Position.Y = y
	})

	r.ConsumeLabel("characteristics")

	ecs.WithChangedValue(r.ReadEnum, func(enum uint32) {
		b.Character = characteristics.Characteristics(enum)
	})

	err := b.Collider.ReadFrom(r)

	if err != nil {
		return err
	}

	r.ConsumeLabel("corners")
	r.ConsumeLine()
	r.ConsumeLine()
	r.ConsumeLine()
	r.ConsumeLine()

	return nil
}
