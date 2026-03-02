package systems

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/tylland/dashland/internal/components"
	"github.com/tylland/dashland/internal/ecs"
	"github.com/tylland/dashland/internal/game"
)

type Renderer struct {
	stage  *game.Stage
	camera game.Camera
}

func NewRenderSystem(stage *game.Stage, camera game.Camera) *Renderer {
	return &Renderer{stage: stage, camera: camera}
}

func (s *Renderer) Update(world *ecs.World, deltaTime float32) {
	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)
	rl.BeginMode2D(s.camera.GetCamera())

	s.stage.Update(deltaTime)
	s.stage.Render(deltaTime)

	for _, entity := range world.Entities() {
		position := ecs.GetComponent[components.PositionComponent](entity)
		sprite := ecs.GetComponent[components.SpriteComponent](entity)

		if position != nil && sprite != nil {
			rl.DrawTextureRec(*sprite.Texture, sprite.Source, position.Vector2, rl.White)
		}

	}

	for _, entity := range world.Entities() {
		flash := ecs.GetComponent[components.FlashComponent](entity)

		if flash != nil {
			rl.DrawRectangle(
				0,
				0,
				int32(rl.GetScreenWidth()),
				int32(rl.GetScreenHeight()),
				rl.White,
			)
		}
	}

	rl.EndMode2D()

	// Draw inventory at the top of the screen
	s.drawInventory(world)

	rl.EndDrawing()
}

func (s *Renderer) drawInventory(world *ecs.World) {
	// Find the player entity
	player := world.GetEntity("player")
	if player == nil {
		return
	}

	inventory := ecs.GetComponent[components.Inventory](player)
	if inventory == nil {
		return
	}

	// Only show inventory if diamonds are required
	if s.stage.DiamondsRequired <= 0 {
		return
	}

	// Draw inventory UI at top-left
	padding := int32(10)
	fontSize := int32(32)
	spriteSize := float32(32) // Size of the diamond sprite in the UI

	// Draw required diamonds text first
	requiredText := fmt.Sprint(s.stage.DiamondsRequired)
	requiredColor := rl.Gray
	if inventory.Diamonds >= s.stage.DiamondsRequired {
		requiredColor = rl.White
	}
	rl.DrawText(requiredText, padding, padding, fontSize, requiredColor)
	textWidth := rl.MeasureText(requiredText, fontSize)

	// Draw diamond sprite after required count
	spriteX := padding + textWidth + 5
	if s.stage.EntityMap != nil {
		entityTextures := s.stage.EntityMap.GetEntityTextures()
		if entityTextures != nil {
			diamondSrcX := float32(4) * s.stage.BlockWidth // Diamond is entity type 4
			diamondSrcRect := rl.NewRectangle(diamondSrcX, 0, s.stage.BlockWidth, s.stage.BlockHeight)
			diamondDstRect := rl.NewRectangle(float32(spriteX), float32(padding), spriteSize, spriteSize)

			rl.DrawTexturePro(*entityTextures, diamondSrcRect, diamondDstRect, rl.NewVector2(0, 0), 0, rl.White)
		}
	}

	// Draw current diamonds counter after sprite
	currentText := fmt.Sprint(inventory.Diamonds)
	textX := spriteX + int32(spriteSize) + 5

	rl.DrawText(currentText, textX, padding, fontSize, rl.White)

	// Draw points after diamonds
	pointsText := fmt.Sprint(inventory.Points)
	pointsX := textX + rl.MeasureText(currentText, fontSize) + 20
	rl.DrawText(pointsText, pointsX, padding, fontSize, rl.White)
}
