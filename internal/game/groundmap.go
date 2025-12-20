package game

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/lafriks/go-tiled"
	"github.com/tylland/dashland/internal/common"
	"github.com/tylland/dashland/internal/ecs"
)

type EntityMap struct {
	MapSize
	entityTextures        *rl.Texture2D
	enemyTextures         *rl.Texture2D
	groundCorners         *rl.Texture2D
	entities              []*ecs.Entity
	InitialPlayerPosition common.BlockPosition
}

func NewGroundMap(mapSize MapSize, entityTextures *rl.Texture2D, enemyTextures *rl.Texture2D, groundCorners *rl.Texture2D) *EntityMap {
	return &EntityMap{MapSize: mapSize, entityTextures: entityTextures, enemyTextures: enemyTextures, groundCorners: groundCorners, entities: []*ecs.Entity{}}
}

func (m *EntityMap) InitPlayerPosition(tiles []*tiled.LayerTile) {
	m.entities = make([]*ecs.Entity, len(tiles))

	for index, tile := range tiles {
		if tile.ID == 1 {
			m.InitialPlayerPosition = common.BlockPosition{X: index % m.Width, Y: index / m.Width}
			return
		}
	}
}

func (m *EntityMap) SetEntity(entity *ecs.Entity, pos common.BlockPosition) {

	if pos.X < 0 || pos.X >= m.Width || pos.Y < 0 || pos.Y >= m.Height {
		return
	}

	m.entities[pos.Y*m.Width+pos.X] = entity
}

func (m *EntityMap) GetEntity(position common.BlockPosition) *ecs.Entity {
	if position.X < 0 || position.X >= m.Width || position.Y < 0 || position.Y >= m.Height {
		return nil
	}

	return m.entities[position.Y*m.Width+position.X]
}
func (m *EntityMap) GetEntityAtPosition(position common.BlockPosition) *ecs.Entity {
	if position.X < 0 || position.X >= m.Width || position.Y < 0 || position.Y >= m.Height {
		return nil
	}

	return m.entities[position.Y*m.Width+position.X]
}

func (m *EntityMap) CheckEntityAtPosition(entityType ecs.EntityType, position common.BlockPosition) bool {
	return m.entities[position.Y*m.Width+position.X].Type == entityType
}

func (m *EntityMap) MoveEntity(source *ecs.Entity, sourcePos common.BlockPosition, targetPos common.BlockPosition) {
	m.entities[targetPos.Y*m.Width+targetPos.X] = source
	m.entities[sourcePos.Y*m.Width+sourcePos.X] = nil
}

func (m *EntityMap) RemoveEntity(doomed *ecs.Entity, position common.BlockPosition) {

	//	sourcePosition := position.CurrentBlockPosition

	m.entities[position.Y*m.Width+position.X] = nil
}
