package systems

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/tylland/dashland/internal/common"
	"github.com/tylland/dashland/internal/components"
	"github.com/tylland/dashland/internal/ecs"
	"github.com/tylland/dashland/internal/game"
)

const moveSpeed = 128 // pixels per second

type BlockMovementSystem struct {
	stage *game.Stage
	BlockResolver
	//	PlayerResolver
}

func NewBlockMovementSystem(stage *game.Stage, blocks BlockResolver) *BlockMovementSystem {
	return &BlockMovementSystem{
		stage:         stage,
		BlockResolver: blocks,
	}
}

func (s *BlockMovementSystem) Update(world *ecs.World, deltaTime float32) {

	// if !s.Player().IsDead {
	// 	s.handlePlayerPush(world)
	// 	s.updatePlayerMovement(deltaTime)
	// }

	// Handle entities movement
	for _, entity := range world.Entities() {
		comps := world.GetComponents(entity)

		input := ecs.GetComponent[components.InputComponent](comps)
		position := ecs.GetComponent[components.PositionComponent](comps)

		if input != nil && position != nil {
			s.UpdatePosition(world, input, position)
		}

		velocity := ecs.GetComponent[components.VelocityComponent](comps)

		if position != nil && velocity != nil {
			s.updateBlockMovement(entity, position, velocity, deltaTime)
		}
	}

}

func (s *BlockMovementSystem) UpdatePosition(world *ecs.World, input *components.InputComponent, position *components.PositionComponent) {

	if s.stage.GetPosition(position.CurrentBlockPosition) == position.Vector2 {
		//	if !position.HasTarget() {
		if input.RightKeyPressed {
			fmt.Print("Move player right!")
			position.SetTarget(common.NewBlockVector(1, 0))
		} else if input.LeftKeyPressed {
			fmt.Print("Move player left!")
			position.SetTarget(common.NewBlockVector(-1, 0))
		} else if input.DownKeyPressed {
			fmt.Print("Move player down!")
			position.SetTarget(common.NewBlockVector(0, 1))
		} else if input.UpKeyPressed {
			fmt.Print("Move player up!")
			position.SetTarget(common.NewBlockVector(0, -1))
		}

		if s.stage.IsObstacleForPlayer(world, position.TargetBlockPosition) {
			position.CancelTarget()
		}
	}
}

// func (s *BlockMovementSystem) updatePlayerMovement(deltaTime float32) {
// 	playerEntity, _ := world.GetEntity("player")
// 	playerPosition := ecs.GetComponent[components.PositionComponent](world.GetComponents(playerEntity))

// 	// if !p.Position.CurrentBlockPosition.IsSame(p.Position.TargetBlockPosition) {
// 	// 	p.Position.UseTarget())
// 	// }

// 	// s.handleBlockMovement(p.Position, p.Velocity, deltaTime)

// 	if !playerPosition.CurrentBlockPosition.IsSame(playerPosition.TargetBlockPosition) {

// 		targetPosition := s.GetPosition(p.Position.TargetBlockPosition)
// 		p.Movement.Start(s.GetPosition(p.Position.CurrentBlockPosition), targetPosition, moveSpeed, nil)
// 		s.VisitBlock(p.Position.TargetBlockPosition)
// 		//		p.game.world.VisitObject(p, p.targetBlockPosition)

// 		p.Position.UseTarget()

// 		p.Position.Vector2 = targetPosition
// 	}

// 	if p.Movement.Moving {
// 		p.Movement.Update(deltaTime)
// 		p.Position.Vector2 = p.Movement.Position
// 	}
// }

// func (s *BlockMovementSystem) handlePlayerPush(world *ecs.World) {
// 	// Get entity at player's target position

// 	playerEntity, _ := world.GetEntity("player")
// 	playerPosition := ecs.GetComponent[components.PositionComponent](world.GetComponents(playerEntity))

// 	if targetEntity := s.stage.GetEntity(playerPosition.TargetBlockPosition); targetEntity != nil {
// 		comps := world.GetComponents(targetEntity)
// 		targetCharacteristic := ecs.GetComponent[components.CharacteristicComponent](comps)
// 		targetPosition := ecs.GetComponent[components.PositionComponent](comps)

// 		if !targetCharacteristic.Has(characteristics.Pushable) {
// 			return
// 		}

// 		// Calculate push direction based on player's movement
// 		if playerPosition.HasTarget() {
// 			pushPos := targetPosition.TargetBlockPosition

// 			if playerPosition.CurrentBlockPosition.X > targetPosition.TargetBlockPosition.X {
// 				pushPos = pushPos.Offset(-1, 0)
// 			} else if playerPosition.CurrentBlockPosition.X < targetPosition.TargetBlockPosition.X {
// 				pushPos = pushPos.Offset(1, 0)
// 			}

// 			// Check if push position is free
// 			if s.CheckBlockAtPosition(game.Void, pushPos) && s.stage.GetEntity(pushPos) == nil {
// 				// Move the pushed entity
// 				//s.world.MoveEntity(targetEntity, pushPos)
// 				targetPosition.TargetBlockPosition = pushPos
// 			}
// 		}
// 	}
// }

func (s *BlockMovementSystem) updateBlockMovement(entity *ecs.Entity, position *components.PositionComponent, velocity *components.VelocityComponent, deltaTime float32) {

	if position.HasTarget() {
		if entity.ID == "player" {
			fmt.Print("Player is here!!")
		}
		s.stage.MoveEntity(entity, position.CurrentBlockPosition, position.TargetBlockPosition)
		position.UseTarget()
	}

	s.handleBlockMovement(position, velocity, deltaTime)
}

// Handle grid-based movement
func (s *BlockMovementSystem) handleBlockMovement(position *components.PositionComponent, velocity *components.VelocityComponent, deltaTime float32) {
	targetPos := s.stage.GetPosition(position.CurrentBlockPosition)
	currentPos := position.Vector2

	if currentPos != targetPos {
		diff := rl.Vector2Subtract(targetPos, currentPos)
		length := rl.Vector2Length(diff)

		if length > 0 {
			moveAmount := moveSpeed * deltaTime

			if moveAmount >= length {
				position.Vector2 = targetPos
			} else {
				moveVec := rl.Vector2Scale(rl.Vector2Normalize(diff), moveAmount)
				position.Vector2 = rl.Vector2Add(currentPos, moveVec)
			}

		}
	} else if !velocity.BlockVector.IsZero() {
		position.SetTarget(velocity.BlockVector)
	}

}
