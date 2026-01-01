package game

import "github.com/tylland/dashland/internal/ecs"

// type EntityCategory uint16

// type EntityWorld interface {
// 	Entities() []*Entity
// }

const (
	EntityCategoryObject ecs.EntityCategory = 0
	EntityCategoryEnemy  ecs.EntityCategory = 100
)

// type EntityType uint16

const (
	EntityPlayer  ecs.EntityType = iota
	EntityDiamond                = 4
	EntityBoulder                = 5
	EntityDoor                   = 11
	EntityEnemy                  = 101
)

// type Entity struct {
// 	Id   common.EntityId
// 	Type EntityType
// 	//	Characteristics characteristics.Characteristics
// 	Characteristic *components.CharacteristicComponent
// 	Position       *components.PositionComponent
// 	Velocity       *components.VelocityComponent
// 	Sprite         *components.SpriteComponent
// 	Collision      *components.CollisionComponent
// 	Collectable    *components.CollectableComponent
// 	WallWalker     *components.WallWalkerComponent
// }
