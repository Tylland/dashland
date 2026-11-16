package game

import "github.com/tylland/dashland/internal/ecs"

const (
	EntityCategoryObject ecs.EntityCategory = 0
	EntityCategoryEnemy  ecs.EntityCategory = 100
	EntityCategoryBlock  ecs.EntityCategory = 200
)

const (
	EntityUnknown   ecs.EntityType = 0
	EntityPlayer    ecs.EntityType = 1
	EntityDiamond   ecs.EntityType = 4
	EntityBoulder   ecs.EntityType = 5
	EntityExplosion ecs.EntityType = 6
	EntityFlash     ecs.EntityType = 7
	EntityMagicWall ecs.EntityType = 8
	//EntityDoor      ecs.EntityType = 11
	EntityExitDoor ecs.EntityType = 12
	// Enemies
	EntityUnknownEnemy ecs.EntityType = 100
	EntityFirefly      ecs.EntityType = 101
	EntityButterfly    ecs.EntityType = 102
	// Blocks
	EntityUnknownBlock ecs.EntityType = 200
	EntityBedrock      ecs.EntityType = 201
	EntityVoid         ecs.EntityType = 202
	EntitySoil         ecs.EntityType = 203
	EntityWall         ecs.EntityType = 204
)

func EntityCategory(entityType ecs.EntityType) ecs.EntityCategory {
	switch {
	case entityType >= EntityUnknownBlock:
		return EntityCategoryBlock
	case entityType >= EntityUnknownEnemy:
		return EntityCategoryEnemy
	default:
		return EntityCategoryObject
	}
}

func EntityCategoryName(category ecs.EntityCategory) string {
	switch category {
	case EntityCategoryObject:
		return "object"
	case EntityCategoryEnemy:
		return "enemy"
	case EntityCategoryBlock:
		return "block"
	default:
		return "unknown"
	}
}

func EntityTypeName(entityType ecs.EntityType) string {
	switch entityType {
	case EntityUnknown:
		return "unknown"
	case EntityPlayer:
		return "player"
	case EntityDiamond:
		return "diamond"
	case EntityBoulder:
		return "boulder"
	case EntityExplosion:
		return "explosion"
	case EntityFlash:
		return "flash"
	case EntityMagicWall:
		return "magicwall"
	case EntityExitDoor:
		return "exitdoor"
	case EntityFirefly:
		return "firefly"
	case EntityButterfly:
		return "butterfly"
	case EntityUnknownBlock:
		return "unknown-block"
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
