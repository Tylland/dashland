package systems

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
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
	// Handle entities movement
	for _, entity := range world.Entities() {
		comps := world.GetComponents(entity)

		position := ecs.GetComponent[components.PositionComponent](comps)
		step := ecs.GetComponent[components.BlockStep](comps)
		collider := ecs.GetComponent[components.ColliderComponent](comps)

		if position != nil && step != nil && collider != nil {
			s.moveEntity(world, entity, position, step, collider, deltaTime)
		}
	}

}

// func (s *BlockMovementSystem) nextStep(world *ecs.World, input *components.InputComponent, position *components.PositionComponent, step *components.BlockStep) {

// 	if position.Vector2 == step.Target {
// 		//	if !position.HasTarget() {

// 		var direction common.BlockVector = common.BlockVector{}

// 		if input.RightKeyPressed {
// 			fmt.Print("Move player right!")
// 			direction = common.DirectionRight
// 		} else if input.LeftKeyPressed {
// 			fmt.Print("Move player left!")
// 			direction = common.DirectionLeft
// 		} else if input.DownKeyPressed {
// 			fmt.Print("Move player down!")
// 			direction = common.DirectionDown
// 		} else if input.UpKeyPressed {
// 			fmt.Print("Move player up!")
// 			direction = common.DirectionUp
// 		}

// 		targetBlockPos := position.CurrentBlockPosition.Add(direction)

// 		step.Move(direction, s.stage.GetPosition(targetBlockPos), moveSpeed)

// 		// if !s.stage.IsObstacleForPlayer(world, targetBlockPos) {
// 		// 	step.Move(direction, s.stage.GetPosition(targetBlockPos), moveSpeed)
// 		// }
// 	}
// }

func (s *BlockMovementSystem) moveEntity(world *ecs.World, entity *ecs.Entity, position *components.PositionComponent, step *components.BlockStep, collider *components.ColliderComponent, deltaTime float32) {
	if entity.ID == "player" {
		fmt.Print("Player is here!!")
	}

	if !step.Increment.IsZero() {

		target := position.CurrentBlockPosition.Add(step.Increment)

		if block, ok := s.stage.GetBlock(target.X, target.Y); ok {
			if collider.CollidesWith(block.Collider) {

			}
		}

		position.Update(target)
		s.stage.MoveEntity(entity, position.PreviousBlockPosition, position.CurrentBlockPosition)
		s.stage.VisitBlock(position.CurrentBlockPosition)
		step.Increment.Clear()
	}

	s.moveToTarget(position, step, deltaTime)
}

// Handle grid-based movement
func (s *BlockMovementSystem) moveToTarget(position *components.PositionComponent, step *components.BlockStep, deltaTime float32) {
	targetPos := step.Target
	currentPos := position.Vector2

	if currentPos != targetPos {
		diff := rl.Vector2Subtract(targetPos, currentPos)
		length := rl.Vector2Length(diff)

		if length > 0 {
			moveAmount := step.Speed * deltaTime

			if moveAmount >= length {
				position.Vector2 = targetPos
			} else {
				moveVec := rl.Vector2Scale(rl.Vector2Normalize(diff), moveAmount)
				position.Vector2 = rl.Vector2Add(currentPos, moveVec)
			}
		}
	}
	// else if step.IsMoving() {
	// 	step.Halt()
	// }

}
