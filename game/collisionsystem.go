package game

import (
	"github.com/tylland/dashland/game/components"
	"github.com/tylland/dashland/game/utils"
)

type CollisionSystem struct {
	world *World
}

func NewCollisionSystem(world *World) *CollisionSystem {
	return &CollisionSystem{
		world: world,
	}
}

func (s *CollisionSystem) Update() {
	s.checkBlockCollisions()
	s.checkPlayerCollisions()

	for _, entity := range s.world.entities {
		if entity != nil && entity.Position != nil && entity.Collision != nil {
			s.CheckEntityCollision(entity, entity.Position, entity.Collision)
		}
	}

}

func (s *CollisionSystem) CheckEntityCollision(entity1 *Entity, position1 *components.PositionComponent, collision1 *components.CollisionComponent) {

	for _, entity2 := range s.world.entities {
		if entity2 != nil && entity2.Position != nil && entity2.Collision != nil {

			if entity1.Id == entity2.Id {
				continue
			}

			position2 := entity2.Position
			collision2 := entity2.Collision
			// Check if these layers can collide
			if !collision1.CanCollideWith(entity2.Collision.Layer) {
				return
			}

			if collision1.BoxBody.Overlaps(position1.Vector2, position2.Vector2, &collision2.BoxBody) {
				collision1.IsColliding = true
				collision2.IsColliding = true

				s.resolveEntityCollision(entity1, entity2)
			}

		}
	}
}

func (cs *CollisionSystem) checkPlayerCollisions() {
	for _, entity := range cs.world.entities {
		if entity == nil || entity.Position == nil || entity.Collision == nil {
			continue
		}

		if entity.Collision.BoxBody.Overlaps(entity.Position.Vector2, cs.world.player.Position.Vector2, &cs.world.player.BoxBody) {
			// Dispatch player collision event
			event := PlayerCollisionEvent{
				Player: cs.world.player,
				Entity: entity,
			}

			cs.world.OnEvent(event)
		}
	}
}

func (cs *CollisionSystem) checkBlockCollisions() {
	for _, entity := range cs.world.entities {
		if entity == nil || entity.Position == nil || entity.Collision == nil {
			continue
		}

		// Get the blocks that this entity might be colliding with
		nearbyBlocks := cs.world.GetNearbyBlocks(&entity.Position.CurrentBlockPosition)

		for _, block := range nearbyBlocks {
			if block.blockType == Soil {
				entityRect := entity.Collision.BoxBody.Rectangle(entity.Position.Vector2)
				blockRect := block.Rectangle()

				// Check overlap of the two rectangles

				if utils.RectangleOverlaps(&entityRect, &blockRect) {
					cs.resolveBlockCollision(entity, block)
				}
			}
		}
	}
}

func (cs *CollisionSystem) resolveEntityCollision(entity1, entity2 *Entity) {
	// Get the rectangles for overlap calculation
	rect1 := entity1.Collision.BoxBody.Rectangle(entity1.Position.Vector2)
	rect2 := entity2.Collision.BoxBody.Rectangle(entity2.Position.Vector2)

	// Calculate overlap on both axes
	//entity1Right := rect1.X + rect1.Width
	entity1Bottom := rect1.Y + rect1.Height
	//entity2Right := rect2.X + rect2.Width
	entity2Bottom := rect2.Y + rect2.Height

	overlapY := utils.Min(entity1Bottom, entity2Bottom) - utils.Max(rect1.Y, rect2.Y)

	// if entity1.Type == Boulder && entity2.Type == Boulder {
	// 	return
	// }

	// Vertical resolution - adjust the upper entity
	if rect1.Y < rect2.Y {
		// entity1 is above entity2
		entity1.Position.Y = rect2.Y - rect1.Height
		if entity1.Velocity.Vector.Y > 0 {
			entity1.Velocity.Vector.Y = 0
		}

	} else {
		// entity2 is above entity1
		entity2.Position.Y += overlapY
		if entity2.Velocity.Vector.Y > 0 {
			entity2.Velocity.Vector.Y = 0
		}
	}

	collisionEvent := EntityCollisionEvent{
		Entity1: entity1,
		Entity2: entity2,
		Layer1:  entity1.Collision.Layer,
		Layer2:  entity2.Collision.Layer,
	}

	cs.world.OnEvent(collisionEvent)
}

type EntityCollisionEvent struct {
	Entity1 *Entity
	Entity2 *Entity
	Layer1  components.CollisionLayer
	Layer2  components.CollisionLayer
}

func (ce EntityCollisionEvent) IsEvent() {}

type BlockCollisionEvent struct {
	Block  *Block
	Entity *Entity
}

func (ce BlockCollisionEvent) IsEvent() {}

type PlayerCollisionEvent struct {
	Player        *Player
	Entity        *Entity
	EntityFalling bool
}

func (ce PlayerCollisionEvent) IsEvent() {}

func (cs *CollisionSystem) resolveBlockCollision(entity *Entity, block *Block) {
	// Get the rectangles for overlap calculation
	entityRect := entity.Collision.BoxBody.Rectangle(entity.Position.Vector2)
	blockRect := block.Rectangle()

	// Calculate overlap on both axes
	// Assuming rectangles have X, Y, Width, Height properties
	entityRight := entityRect.X + entityRect.Width
	entityBottom := entityRect.Y + entityRect.Height
	blockRight := blockRect.X + blockRect.Width
	blockBottom := blockRect.Y + blockRect.Height

	overlapX := utils.Min(entityRight, blockRight) - utils.Max(entityRect.X, blockRect.X)
	overlapY := utils.Min(entityBottom, blockBottom) - utils.Max(entityRect.Y, blockRect.Y)

	// Resolve collision by moving entity out of the block
	// Move in the direction of smallest overlap
	if overlapX < overlapY {
		// Horizontal resolution
		entityCenterX := entityRect.X + entityRect.Width/2
		blockCenterX := blockRect.X + blockRect.Width/2
		if entityCenterX < blockCenterX {
			entity.Position.X = blockRect.X - entityRect.Width
		} else {
			entity.Position.X = blockRect.X + blockRect.Width
		}

		// Stop horizontal momentum
		entity.Velocity.Vector.X = 0
	} else {
		// Vertical resolution
		entityCenterY := entityRect.Y + entityRect.Height/2
		blockCenterY := blockRect.Y + blockRect.Height/2
		if entityCenterY < blockCenterY {
			entity.Position.Y = blockRect.Y - entityRect.Height
			// If falling onto the block
			if entity.Velocity.Vector.Y > 0 {
				entity.Velocity.Vector.Y = 0
			}
		} else {
			entity.Position.Y = blockRect.Y + blockRect.Height
			// If hitting head on block
			if entity.Velocity.Vector.Y < 0 {
				entity.Velocity.Vector.Y = 0
			}
		}
	}

	// Dispatch collision event
	collisionEvent := BlockCollisionEvent{
		Entity: entity,
		Block:  block,
	}

	cs.world.OnEvent(collisionEvent)
}
