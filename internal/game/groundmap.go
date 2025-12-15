package game

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/lafriks/go-tiled"
	"github.com/tylland/dashland/internal/common"
	"github.com/tylland/dashland/internal/components"
	"github.com/tylland/dashland/internal/ecs"
)

type GroundMap struct {
	MapSize
	entityTextures        *rl.Texture2D
	enemyTextures         *rl.Texture2D
	groundCorners         *rl.Texture2D
	entities              []*ecs.Entity
	InitialPlayerPosition common.BlockPosition
}

func NewGroundMap(mapSize MapSize, entityTextures *rl.Texture2D, enemyTextures *rl.Texture2D, groundCorners *rl.Texture2D) *GroundMap {
	return &GroundMap{MapSize: mapSize, entityTextures: entityTextures, enemyTextures: enemyTextures, groundCorners: groundCorners, entities: []*ecs.Entity{}}
}

func (gm *GroundMap) InitPlayerPosition(tiles []*tiled.LayerTile) {
	gm.entities = make([]*ecs.Entity, len(tiles))

	for index, tile := range tiles {
		if tile.ID == 1 {
			gm.InitialPlayerPosition = common.BlockPosition{X: index % gm.Width, Y: index / gm.Width}
			return
		}
	}
}

func (gm *GroundMap) SetEntity(entity *ecs.Entity, pos common.BlockPosition) {

	if pos.X < 0 || pos.X >= gm.Width || pos.Y < 0 || pos.Y >= gm.Height {
		return
	}

	gm.entities[pos.Y*gm.Width+pos.X] = entity
}

func (gm *GroundMap) GetEntity(position common.BlockPosition) *ecs.Entity {
	if position.X < 0 || position.X >= gm.Width || position.Y < 0 || position.Y >= gm.Height {
		return nil
	}

	return gm.entities[position.Y*gm.Width+position.X]
}
func (gm *GroundMap) GetEntityAtPosition(position common.BlockPosition) *ecs.Entity {
	if position.X < 0 || position.X >= gm.Width || position.Y < 0 || position.Y >= gm.Height {
		return nil
	}

	return gm.entities[position.Y*gm.Width+position.X]
}

func (gm *GroundMap) CheckEntityAtPosition(entityType ecs.EntityType, position common.BlockPosition) bool {
	return gm.entities[position.Y*gm.Width+position.X].Type == entityType
}

func (gm *GroundMap) MoveEntity(source *ecs.Entity, position *components.PositionComponent, targetPos common.BlockPosition) {
	sourcePosition := position.CurrentBlockPosition

	gm.entities[targetPos.Y*gm.Width+targetPos.X] = source
	gm.entities[sourcePosition.Y*gm.Width+sourcePosition.X] = nil

	position.Update(targetPos)
}

func (gm *GroundMap) RemoveEntity(doomed *ecs.Entity, position *components.PositionComponent) {

	sourcePosition := position.CurrentBlockPosition

	gm.entities[sourcePosition.Y*gm.Width+sourcePosition.X] = nil
}
