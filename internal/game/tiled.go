package game

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/lafriks/go-tiled"
	"github.com/tylland/dashland/internal/common"
	"github.com/tylland/dashland/internal/ecs"
)

func InitBlocks(world *ecs.World, s *Stage, tiles []*tiled.LayerTile) {
	s.blocks = make([]*ecs.Entity, s.MapSize.Count())

	for index, tile := range tiles {
		blockPosition := common.NewBlockPositionFromIndex(index, s.BlockMap.Width)
		header := NewEntityHeader(ecs.EntityType(uint32(EntityCategoryBlock) + tile.ID))
		block := NewBlock(header, s.BlockMap, s.EntityMap, blockPosition)
		s.BlockMap.SetBlock(block, blockPosition)
		world.AddEntity(block)
	}

	//	s.BlockMap.PrintBlockMap()
}

func InitEntities(world *ecs.World, stage *Stage, category ecs.EntityCategory, tiles []*tiled.LayerTile) {
	for index, tile := range tiles {
		entityType := ecs.EntityType(uint32(category) + tile.ID)

		if entityType == EntityUnknown || entityType == EntityUnknownEnemy {
			continue
		}

		blockPosition := common.NewBlockPositionFromIndex(index, stage.EntityMap.Width)

		entity, err := NewEntityAndHeader(entityType, world, stage, blockPosition)

		if err != nil {
			fmt.Printf("%s at %s %d:%d: %s\n", EntityCategoryName(category), stage.Name, blockPosition.X, blockPosition.Y, err.Error())
			continue
		}

		stage.EntityMap.SetEntity(entity, blockPosition)
	}
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

func InitObjectsEntities(world *ecs.World, s *Stage, category ecs.EntityCategory, objectLayer *tiled.ObjectGroup) {
	for _, obj := range objectLayer.Objects {
		if obj.Type == "EntityExitDoor" {

			blockPos := s.GetBlockPosition(rl.Vector2{X: float32(obj.X), Y: float32(obj.Y)})
			targetStage := getString(obj, "Stage")
			targetPosition, err := common.ParseBlockPosition(getString(obj, "Position"))

			if err != nil {
				panic(err.Error())
			}

			header := NewEntityHeader(EntityExitDoor)
			door := NewExitDoor(header, world, s, blockPos, s.GetPosition(blockPos), targetStage, targetPosition)
			s.EntityMap.SetEntity(door, blockPos)
			s.ExitPosition = blockPos
		}
	}
}
