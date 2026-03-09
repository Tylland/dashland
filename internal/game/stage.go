package game

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/lafriks/go-tiled"
	"github.com/tylland/dashland/internal/characteristics"
	"github.com/tylland/dashland/internal/common"
	"github.com/tylland/dashland/internal/components"
	"github.com/tylland/dashland/internal/ecs"
)

type SoundPlayer interface {
	PlayFx(name string)
}

type Stage struct {
	MapSize
	*BlockMap
	*EntityMap
	Name               string
	EnterPosition      common.BlockPosition
	DiamondsRequired   int
	DiamondPoints      int
	DiamondBonusPoints int
	ExitCondition      bool
	ExitPosition       common.BlockPosition
	MagicWallActive    bool
	MagicWallTimer     float32
	MagicWallDuration  float32
}

func NewStage(name string, size MapSize, blockTexture, entityTextures, groundCorners *rl.Texture2D) *Stage {
	return &Stage{Name: name, MapSize: size, BlockMap: NewBlockMap(size, blockTexture), EntityMap: NewEntityMap(size, entityTextures, groundCorners), MagicWallDuration: 20.0}
}

func (s *Stage) GetPosition(position common.BlockPosition) rl.Vector2 {
	return rl.NewVector2(float32(position.X)*s.BlockMap.BlockWidth, float32(position.Y)*s.BlockMap.BlockHeight)
}

func (s *Stage) GetBlockPosition(position rl.Vector2) common.BlockPosition {
	bx := int(position.X / float32(s.BlockWidth))
	by := int(position.Y / float32(s.BlockHeight))

	return common.BlockPosition{X: bx, Y: by}
}

func (w *Stage) Render(deltaTime float32) {
	for _, block := range w.blocks {
		block.Render()
	}

}

func (s *Stage) VisitBlock(position common.BlockPosition) {
	fmt.Printf("Block at position %d,%d changed from %s", position.X, position.Y, s.blocks[position.Y*s.Width+position.X].BlockType.String())
	s.SetBlock(NewBlock(s.BlockMap, s.EntityMap, Void, position), position)
	fmt.Printf(" to %s \n", s.blocks[position.Y*s.Width+position.X].BlockType.String())
}

func (s *Stage) ActivateMagicWall() {
	if !s.MagicWallActive {
		s.MagicWallActive = true
		s.MagicWallTimer = s.MagicWallDuration
		fmt.Println("Magic wall activated!")
	}
}

func (s *Stage) UpdateMagicWall(deltaTime float32) {
	if s.MagicWallActive {
		s.MagicWallTimer -= deltaTime
		if s.MagicWallTimer <= 0 {
			s.MagicWallActive = false
			s.MagicWallTimer = 0
			fmt.Println("Magic wall expired!")
		}
	}
}

func (s *Stage) IsMagicWallActive() bool {
	return s.MagicWallActive
}

func (s *Stage) CheckCharacteristics(position common.BlockPosition, character characteristics.Characteristics) bool {

	if block, ok := s.GetBlock(position.X, position.Y); ok {
		if block.HasCharacteristic(character) {
			return true
		}
	}

	if entity, ok := s.GetEntityAtPosition(position); ok {
		characteristics := ecs.GetComponent[components.CharacteristicComponent](entity)

		return characteristics != nil && characteristics.Has(character)
	}

	return false
}

func (s *Stage) CheckBlocked(world *ecs.World, position common.BlockPosition, collider *components.ColliderComponent) bool {

	if block, ok := s.GetBlock(position.X, position.Y); ok {
		blocked, _ := collider.Result(block.Collider)
		if blocked {
			return true
		}
	}

	if entity, ok := s.GetEntityAtPosition(position); ok {
		entCharacter := ecs.GetComponent[components.CharacteristicComponent](entity)
		entCollider := ecs.GetComponent[components.ColliderComponent](entity)

		if entCharacter != nil && entCollider != nil {
			// Falling entities are not considered obstacles
			if entCharacter.Has(characteristics.Falling) {
				return false
			}
			blocked, _ := collider.Result(entCollider)
			if blocked {
				return true
			}
		}
	}

	return false
}

func (s *Stage) CheckPositionOccupied(position common.BlockPosition) bool {
	if !s.CheckBlockAtPosition(Void, position) {
		return true
	}

	return s.GetEntity(position) != nil
}

func (s *Stage) InitEntities(world *ecs.World, category ecs.EntityCategory, tiles []*tiled.LayerTile) {
	for index, tile := range tiles {
		blockPosition := common.NewBlockPositionFromIndex(index, s.EntityMap.Width)

		entity, err := NewGameEntity(world, s, ecs.EntityType(uint32(category)+tile.ID), blockPosition)

		if err == nil {
			s.EntityMap.SetEntity(entity, blockPosition)
		}
	}
}

func (s *Stage) InitBlocks(world *Stage, tiles []*tiled.LayerTile) {
	s.blocks = make([]*Block, len(tiles))

	for index, tile := range tiles {
		blockPosition := common.NewBlockPositionFromIndex(index, s.BlockMap.Width)
		s.BlockMap.SetBlock(NewBlock(s.BlockMap, s.EntityMap, BlockType(tile.ID), blockPosition), blockPosition)
	}

	s.BlockMap.PrintBlockMap()
}

// InitBlocksFromGrid populates the block map from a flat BlockType slice (used by procedural generation).
func (s *Stage) InitBlocksFromGrid(blockTypes []BlockType) {
	s.blocks = make([]*Block, len(blockTypes))

	for i, bt := range blockTypes {
		pos := common.NewBlockPositionFromIndex(i, s.BlockMap.Width)
		s.BlockMap.SetBlock(NewBlock(s.BlockMap, s.EntityMap, bt, pos), pos)
	}

	s.BlockMap.PrintBlockMap()
}

// InitEmptyEntityMap allocates the entity grid without requiring Tiled layer data.
func (s *Stage) InitEmptyEntityMap() {
	s.EntityMap.entities = make([]*ecs.Entity, s.Width*s.Height)
}

func getString(obj *tiled.Object, name string) string {
	value := ""
	for _, p := range obj.Properties {
		if p.Name == name {
			var val any = p.Value
			switch v := val.(type) {
			case string:
				value = v
			case float64:
				value = fmt.Sprintf("%v", v)
			default:
				value = fmt.Sprintf("%v", v)
			}
			break
		}
	}

	return value
}

func (s *Stage) InitObjectsEntities(world *ecs.World, category ecs.EntityCategory, objectLayer *tiled.ObjectGroup) {
	for _, obj := range objectLayer.Objects {
		if obj.Type == "EntityExitDoor" {

			blockPos := s.GetBlockPosition(rl.Vector2{X: float32(obj.X), Y: float32(obj.Y)})
			targetStage := getString(obj, "Stage")
			targetPosition, err := common.ParseBlockPosition(getString(obj, "Position"))

			if err != nil {
				panic(err.Error())
			}

			door := NewExitDoor(world, s, blockPos, s.GetPosition(blockPos), targetStage, targetPosition)
			s.EntityMap.SetEntity(door, blockPos)
			s.ExitPosition = blockPos
		}
	}
}

// func (s *Stage) SpawnExitDoor(world *ecs.World, blockPos common.BlockPosition, targetStage string, targetPos common.BlockPosition) *ecs.Entity {
// 	door := NewExitDoor(world, s, blockPos, s.GetPosition(blockPos), targetStage, targetPos)
// 	s.EntityMap.SetEntity(door, blockPos)
// 	return door
// }
