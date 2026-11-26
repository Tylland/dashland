package types

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

var (
	ComponentTypeAnimation       = ecs.ComponentType{Name: "animation", Version: 1}
	ComponentTypeBlock           = ecs.ComponentType{Name: "corners", Version: 2}
	ComponentTypeCharacteristics = ecs.ComponentType{Name: "characteristics", Version: 1}
	ComponentTypeCollectable     = ecs.ComponentType{Name: "collectable", Version: 1}
	ComponentTypeCollider        = ecs.ComponentType{Name: "collider", Version: 1}
	ComponentTypeDamage          = ecs.ComponentType{Name: "damage", Version: 1}
	ComponentTypeDoor            = ecs.ComponentType{Name: "door", Version: 1}
	ComponentTypeFlash           = ecs.ComponentType{Name: "flash", Version: 1}
	ComponentTypeHealth          = ecs.ComponentType{Name: "health", Version: 1}
	ComponentTypeInput           = ecs.ComponentType{Name: "input", Version: 1}
	ComponentTypePosition        = ecs.ComponentType{Name: "position", Version: 1}
	ComponentTypePushable        = ecs.ComponentType{Name: "pushable", Version: 1}
	ComponentTypeSprite          = ecs.ComponentType{Name: "sprite", Version: 1}
	ComponentTypeBlockStep       = ecs.ComponentType{Name: "blockstep", Version: 1}
)
