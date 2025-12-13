package game

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/lafriks/go-tiled"
	"github.com/tylland/dashland/game/core"
)

type GroundMap struct {
	MapSize
	entityTextures        *rl.Texture2D
	enemyTextures         *rl.Texture2D
	groundCorners         *rl.Texture2D
	entities              []*Entity
	initialPlayerPosition core.BlockPosition
}

func NewGroundMap(mapSize MapSize, entityTextures *rl.Texture2D, enemyTextures *rl.Texture2D, groundCorners *rl.Texture2D) *GroundMap {
	return &GroundMap{MapSize: mapSize, entityTextures: entityTextures, enemyTextures: enemyTextures, groundCorners: groundCorners, entities: []*Entity{}}
}

func (gm *GroundMap) InitPlayerPosition(tiles []*tiled.LayerTile) {
	gm.entities = make([]*Entity, len(tiles))

	for index, tile := range tiles {
		if tile.ID == 1 {
			gm.initialPlayerPosition = core.BlockPosition{X: index % gm.width, Y: index / gm.width}
			return
		}
	}
}

func (gm *GroundMap) InitEntities(world *World, category EntityCategory, tiles []*tiled.LayerTile) {
	for index, tile := range tiles {
		entity, err := NewGameEntity(world, EntityType(uint32(category)+tile.ID), index%gm.width, index/gm.width)

		if err == nil {
			gm.entities[index] = entity
		}
	}
}

func (gm *GroundMap) AddEntities(world *World, category EntityCategory, tiles []*tiled.LayerTile) {

	for index, tile := range tiles {
		entity, err := NewGameEntity(world, EntityType(uint32(category)+tile.ID), index%world.width, index/world.width)

		if err == nil {
			gm.entities[index] = entity
		}
	}
}

func (gm *GroundMap) SetEntity(entity *Entity, pos core.BlockPosition) {

	if pos.X < 0 || pos.X >= gm.width || pos.Y < 0 || pos.Y >= gm.height {
		return
	}

	gm.entities[pos.Y*gm.width+pos.X] = entity
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

func (gm *GroundMap) CheckEntityAtPosition(entityType EntityType, position core.BlockPosition) bool {
	return gm.entities[position.Y*gm.width+position.X].Type == entityType
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
