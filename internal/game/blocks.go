package game

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/tylland/dashland/internal/characteristics"
	"github.com/tylland/dashland/internal/components"
	"github.com/tylland/dashland/internal/ecs"
	"github.com/tylland/dashland/internal/game/types"
)

const (
	BlockWidth  float32 = 32
	BlockHeight float32 = 32
)

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

func BlockHasCharacteristic(entity *ecs.Entity, characteristic characteristics.Characteristics) bool {
	component := ecs.GetComponent[components.CharacteristicComponent](entity)
	return component != nil && component.Has(characteristic)
}

// cornerIndex returns the corner texture index for a block corner.
// A strict soil corner requires all three neighboring cells to match,
// while void corners allow the diagonal neighbor to differ.
func cornerIndex(neighbors [9]bool, corner CornerPosition, strict bool) uint8 {
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

func UpdateBlock(entity *ecs.Entity, blockMap *BlockMap) {
	position := ecs.GetComponent[components.PositionComponent](entity)
	block := ecs.GetComponent[components.BlockComponent](entity)
	if position == nil || block == nil {
		return
	}

	switch entity.Type {
	case types.EntitySoil:
		neighbors := blockMap.CheckNeighbourTypes(types.EntityVoid, position.CurrentBlockPosition)
		block.Corners[CornerOverLeft] = cornerIndex(neighbors, CornerOverLeft, true)
		block.Corners[CornerOverRight] = cornerIndex(neighbors, CornerOverRight, true)
		block.Corners[CornerUnderRight] = cornerIndex(neighbors, CornerUnderRight, true)
		block.Corners[CornerUnderLeft] = cornerIndex(neighbors, CornerUnderLeft, true)
	case types.EntityVoid:
		neighbors := blockMap.CheckNeighbourTypes(types.EntitySoil, position.CurrentBlockPosition)
		block.Corners[CornerOverLeft] = cornerIndex(neighbors, CornerOverLeft, false)
		block.Corners[CornerOverRight] = cornerIndex(neighbors, CornerOverRight, false)
		block.Corners[CornerUnderRight] = cornerIndex(neighbors, CornerUnderRight, false)
		block.Corners[CornerUnderLeft] = cornerIndex(neighbors, CornerUnderLeft, false)
	}
}

func WriteBlockTo(entity *ecs.Entity, w ecs.Writer) error {
	if err := entity.EntityHeader.WriteTo(w); err != nil {
		return err
	}
	position := ecs.GetComponent[components.PositionComponent](entity)
	character := ecs.GetComponent[components.CharacteristicComponent](entity)
	collider := ecs.GetComponent[components.ColliderComponent](entity)
	block := ecs.GetComponent[components.BlockComponent](entity)
	if position == nil || character == nil || collider == nil || block == nil {
		return fmt.Errorf("block entity %q is missing required components", entity.ID)
	}

	w.StartSection("position", 1)
	w.WriteInt("X", position.CurrentBlockPosition.X, false)
	w.WriteInt("Y", position.CurrentBlockPosition.Y, false)
	w.StartSection("characteristics", 1)
	w.WriteEnum("Characteristics", character.Value(), true, "")
	if err := collider.WriteTo(w); err != nil {
		return err
	}
	w.StartSection("corners", 1)
	w.WriteInt("Upper left", int(block.Corners[0]), block.Corners[0] == 0)
	w.WriteInt("Upper right", int(block.Corners[1]), block.Corners[1] == 0)
	w.WriteInt("Lower left", int(block.Corners[2]), block.Corners[2] == 0)
	w.WriteInt("Lower right", int(block.Corners[3]), block.Corners[3] == 0)
	return nil
}

func ReadBlockFrom(entity *ecs.Entity, r ecs.Reader, blockMap *BlockMap) error {
	if err := entity.EntityHeader.ReadFrom(r); err != nil {
		return err
	}
	position := ecs.GetComponent[components.PositionComponent](entity)
	character := ecs.GetComponent[components.CharacteristicComponent](entity)
	collider := ecs.GetComponent[components.ColliderComponent](entity)
	block := ecs.GetComponent[components.BlockComponent](entity)
	if position == nil || character == nil || collider == nil || block == nil {
		return fmt.Errorf("block entity %q is missing required components", entity.ID)
	}

	r.ConsumeSection("position")
	if err := ecs.WithChangedValue(r.ReadInt, func(x int) { position.CurrentBlockPosition.X = x }); err != nil {
		return err
	}
	if err := ecs.WithChangedValue(r.ReadInt, func(y int) { position.CurrentBlockPosition.Y = y }); err != nil {
		return err
	}
	position.PreviousBlockPosition = position.CurrentBlockPosition
	position.Vector2 = rl.NewVector2(float32(position.CurrentBlockPosition.X)*blockMap.BlockWidth, float32(position.CurrentBlockPosition.Y)*blockMap.BlockHeight)

	r.ConsumeSection("characteristics")
	if err := ecs.WithChangedValue(r.ReadEnum, func(value uint32) { character.Set(characteristics.Characteristics(value)) }); err != nil {
		return err
	}
	if err := collider.ReadFrom(r); err != nil {
		return err
	}
	r.ConsumeSection("corners")
	for i := range block.Corners {
		if err := ecs.WithChangedValue(r.ReadInt, func(value int) { block.Corners[i] = uint8(value) }); err != nil {
			return err
		}
	}
	return nil
}
