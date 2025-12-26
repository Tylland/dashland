package game

import (
	"github.com/tylland/dashland/internal/components"
	"github.com/tylland/dashland/internal/ecs"
)

// type GameEvent interface {
// 	World() *ecs.World
// 	IsEvent()
// }

// type GameEventListner interface {
// 	OnEvent(event GameEvent)
// }

// type GameEventDispatcher struct {
// 	listeners []GameEventListner
// }

// func (d *GameEventDispatcher) AddListener(listener GameEventListner) {
// 	d.listeners = append(d.listeners, listener)
// }

type EntityCollisionEvent struct {
	Entity1 *ecs.Entity
	Entity2 *ecs.Entity
}

func NewEntityCollisionEvent(e1 *ecs.Entity, e2 *ecs.Entity) *EntityCollisionEvent {
	return &EntityCollisionEvent{
		Entity1: e1,
		Entity2: e2,
	}
}

type EntityEvent struct {
	Actor  *ecs.Entity
	Target *ecs.Entity
}

func NewEntityEvent(actor *ecs.Entity, target *ecs.Entity) *EntityEvent {
	return &EntityEvent{
		Actor:  actor,
		Target: target,
	}
}

type BlockEvent struct {
	Block  *Block
	Entity *ecs.Entity
}

func NewBlockEvent(block *Block, entity *ecs.Entity) *BlockEvent {
	return &BlockEvent{
		Block:  block,
		Entity: entity,
	}
}

type BlockCollisionEvent struct {
	Block  *Block
	Entity *ecs.Entity
}

func NewBlockCollisionEvent(block *Block, entity *ecs.Entity) *BlockCollisionEvent {
	return &BlockCollisionEvent{
		Block:  block,
		Entity: entity,
	}
}

type PlayerCollisionEvent struct {
	world          *ecs.World
	Player         *ecs.Entity
	Entity         *ecs.Entity
	EntityPosition *components.PositionComponent
	EntityFalling  bool
}

func NewPlayerCollisionEvent(world *ecs.World, player *ecs.Entity, entity *ecs.Entity, entityPosition *components.PositionComponent, entityFalling bool) *PlayerCollisionEvent {
	return &PlayerCollisionEvent{
		world:          world,
		Player:         player,
		Entity:         entity,
		EntityPosition: entityPosition,
		EntityFalling:  entityFalling,
	}
}

type CollectEvent struct {
	Collector   *ecs.Entity
	Collectable *ecs.Entity
}

func NewCollectEvent(collector *ecs.Entity, collectable *ecs.Entity) *CollectEvent {
	return &CollectEvent{
		Collector:   collector,
		Collectable: collectable,
	}
}

type DamageEvent struct {
	Source *ecs.Entity
	Target *ecs.Entity
	Damage *components.Damage
	Health *components.Health
}

func NewDamageEvent(source *ecs.Entity, target *ecs.Entity, damage *components.Damage, health *components.Health) *DamageEvent {
	return &DamageEvent{
		Source: source,
		Target: target,
		Damage: damage,
		Health: health,
	}
}
