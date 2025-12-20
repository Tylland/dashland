package game

import (
	"github.com/tylland/dashland/internal/components"
	"github.com/tylland/dashland/internal/ecs"
)

type GameEvent interface {
	World() *ecs.World
	IsEvent()
}

type GameEventListner interface {
	OnEvent(event GameEvent)
}

type GameEventDispatcher struct {
	listeners []GameEventListner
}

func (d *GameEventDispatcher) AddListener(listener GameEventListner) {
	d.listeners = append(d.listeners, listener)
}

type EntityCollisionEvent struct {
	world   *ecs.World
	Entity1 *ecs.Entity
	Entity2 *ecs.Entity
}

func NewEntityCollisionEvent(world *ecs.World, e1 *ecs.Entity, e2 *ecs.Entity) *EntityCollisionEvent {
	return &EntityCollisionEvent{
		world:   world,
		Entity1: e1,
		Entity2: e2,
	}
}

func (ce EntityCollisionEvent) World() *ecs.World {
	return ce.world
}

func (ce EntityCollisionEvent) IsEvent() {}

type BlockCollisionEvent struct {
	world  *ecs.World
	Block  *Block
	Entity *ecs.Entity
}

func NewBlockCollisionEvent(world *ecs.World, block *Block, entity *ecs.Entity) *BlockCollisionEvent {
	return &BlockCollisionEvent{
		world:  world,
		Block:  block,
		Entity: entity,
	}
}

func (ce BlockCollisionEvent) World() *ecs.World {
	return ce.world
}

func (ce BlockCollisionEvent) IsEvent() {}

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

func (ce PlayerCollisionEvent) World() *ecs.World {
	return ce.world
}

func (ce PlayerCollisionEvent) IsEvent() {}
