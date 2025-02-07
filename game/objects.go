package game

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/lafriks/go-tiled"
	"github.com/tylland/dashland/game/components"
	"github.com/tylland/dashland/game/core"
)

var entityId core.EntityId = 0

func NewEntityId() core.EntityId {
	entityId++

	return entityId
}

func NewGameEntity(world *world, blockType BlockType, x int, y int) (Entity, error) {
	blockPosition := core.BlockPosition{X: x, Y: y}
	position := world.GetPosition(blockPosition)

	switch blockType {
	// case Soil:
	// 	return NewSoil(world, NewEntityId(), blockPosition, position), nil
	case Diamond:
		return NewDiamond(world, NewEntityId(), blockPosition, position), nil
	case Boulder:
		return NewStone(world, NewEntityId(), blockPosition, position), nil
	default:
		return Entity{}, fmt.Errorf("%v (%d) is unknown blocktype", blockType, int(blockType))
	}
}

type GroundMap struct {
	MapSize
	objectTextures rl.Texture2D
	groundCorners  rl.Texture2D
	entities       []*Entity
}

func (gm *GroundMap) InitEntities(world *world, tiles []*tiled.LayerTile) {
	gm.entities = make([]*Entity, len(tiles))

	for index, tile := range tiles {
		entity, err := NewGameEntity(world, BlockType(tile.ID), index%world.width, index/world.width)

		if err == nil {
			gm.entities[index] = &entity
		}
	}
}

func (gm *GroundMap) GetEntity(position core.BlockPosition) *Entity {
	if position.X < 0 || position.X >= gm.width || position.Y < 0 || position.Y >= gm.height {
		return nil
	}

	return gm.entities[position.Y*gm.width+position.X]
}
func (gm *GroundMap) GetEntityAtPosition(position core.BlockPosition) *Entity {
	if position.X < 0 || position.X >= gm.width || position.Y < 0 || position.Y >= gm.height {
		return nil
	}

	return gm.entities[position.Y*gm.width+position.X]
}

func (gm *GroundMap) CheckEntityAtPosition(blockType BlockType, position core.BlockPosition) bool {
	return gm.entities[position.Y*gm.width+position.X].Type == blockType
}

func (gm *GroundMap) MoveEntity(source *Entity, targetPos core.BlockPosition) {
	sourcePosition := source.Position.CurrentBlockPosition

	gm.entities[targetPos.Y*gm.width+targetPos.X] = source
	gm.entities[sourcePosition.Y*gm.width+sourcePosition.X] = nil

	source.Position.Update(targetPos)
}

func (gm *GroundMap) RemoveEntity(doomed *Entity) {

	sourcePosition := doomed.Position.CurrentBlockPosition

	gm.entities[sourcePosition.Y*gm.width+sourcePosition.X] = nil
}

// const (
// 	NoComponent                          = core.ComponentType(0)
// 	PositionComponent core.ComponentType = 1 << iota
// 	VelocityComponent
// 	SpriteComponent
// )

func NewSoil(world *world, id core.EntityId, blockPosition core.BlockPosition, position rl.Vector2) Entity {
	return Entity{
		Id:       id,
		Type:     Soil,
		Position: &components.PositionComponent{CurrentBlockPosition: blockPosition, Vector2: position},
		Sprite: &components.SpriteComponent{Sprite: core.Sprite{
			Texture: &world.objectTextures,
			Source:  rl.NewRectangle(float32(Soil)*world.blockWidth, 0, world.blockWidth, world.blockHeight)}},
		Collision: components.NewCollisionComponent(world.blockWidth, world.blockHeight, components.LayerAll),
	}
}

func NewStone(world *world, id core.EntityId, blockPosition core.BlockPosition, position rl.Vector2) Entity {
	return Entity{
		Id:       id,
		Type:     Boulder,
		Behavior: CanFall | Obstacle | Pushable,
		Position: components.NewPositionComponent(blockPosition, position),
		Velocity: &components.VelocityComponent{Vector: core.Vector{X: 0, Y: 0}},
		Sprite: &components.SpriteComponent{Sprite: core.Sprite{
			Texture: &world.objectTextures,
			Source:  rl.NewRectangle(float32(Boulder)*world.blockWidth, 0, world.blockWidth, world.blockHeight)}},
		Collision: components.NewCollisionComponent(world.blockWidth, world.blockHeight, components.LayerAll),
	}
}

func NewDiamond(world *world, id core.EntityId, blockPosition core.BlockPosition, position rl.Vector2) Entity {
	return Entity{
		Id:       id,
		Type:     Diamond,
		Behavior: CanFall | Collectable,
		Position: components.NewPositionComponent(blockPosition, position),
		Velocity: &components.VelocityComponent{Vector: core.Vector{X: 0, Y: 0}},
		Sprite: &components.SpriteComponent{Sprite: core.Sprite{
			Texture: &world.objectTextures,
			Source:  rl.NewRectangle(float32(Diamond)*world.blockWidth, 0, world.blockWidth, world.blockHeight)}},
		Collision:   components.NewCollisionComponent(world.blockWidth, world.blockHeight, components.LayerAll),
		Collectable: components.NewCollectableComponent(components.CollectableDiamond, 1),
	}
}
