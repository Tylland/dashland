package systems

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/tylland/dashland/internal/components"
	"github.com/tylland/dashland/internal/ecs"
	"github.com/tylland/dashland/internal/game"
	"github.com/tylland/dashland/internal/game/types"
)

type Renderer struct {
	stage       *game.Stage
	camera      game.Camera
	elapsedTime float64
	mist        rl.Texture2D
}

func NewRender(stage *game.Stage, camera game.Camera) *Renderer {
	// Soft radial puff texture used by every mist particle
	img := rl.GenImageGradientRadial(256, 256, 0.0, rl.White, rl.Blank)
	mist := rl.LoadTextureFromImage(img)
	rl.UnloadImage(img)
	rl.SetTextureFilter(mist, rl.FilterBilinear)

	return &Renderer{stage: stage, camera: camera, mist: mist}
}

func (s *Renderer) Update(world *ecs.World, deltaTime float32) {
	s.elapsedTime += float64(deltaTime)

	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)
	rl.BeginMode2D(s.camera.GetCamera())

	for _, entity := range world.Entities() {
		if entity.Category == types.EntityCategoryBlock {
			position := ecs.GetComponent[components.Position](entity)
			sprite := ecs.GetComponent[components.Sprite](entity)

			if position != nil && sprite != nil {
				s.drawBlock(entity, position, sprite)
			}
		}
	}

	for _, entity := range world.Entities() {
		if entity.Category == types.EntityCategoryBlock {
			continue
		}

		position := ecs.GetComponent[components.Position](entity)
		sprite := ecs.GetComponent[components.Sprite](entity)

		if position != nil && sprite != nil {
			rl.DrawTextureRec(*sprite.Texture, sprite.Source, position.Vector2, rl.White)
		}
	}

	for _, entity := range world.Entities() {
		flash := ecs.GetComponent[components.Flash](entity)

		if flash != nil {
			rl.DrawRectangle(
				0,
				0,
				int32(rl.GetScreenWidth()),
				int32(rl.GetScreenHeight()),
				rl.White,
			)
		}

		particle := ecs.GetComponent[components.Particle](entity)

		if particle != nil {
			s.drawParticle(particle)
		}
	}

	rl.EndMode2D()

	rl.DrawFPS(10, int32(rl.GetScreenHeight()-30))
	s.drawTimer()
	s.drawInventory(world)

	if s.stage.TimeExpired {
		s.drawTimeExpired()
	}

	rl.EndDrawing()
}

func (s *Renderer) RunWhenPaused() bool {
	return true
}

func (s *Renderer) drawBlock(entity *ecs.Entity, position *components.Position, sprite *components.Sprite) {
	block := ecs.GetComponent[components.BlockCorners](entity)

	if !block.HasCorners {
		rl.DrawTextureRec(
			*sprite.Texture,
			rl.NewRectangle(float32(entity.Type)*s.stage.BlockWidth, 0, s.stage.BlockWidth, s.stage.BlockHeight),
			position.Vector2,
			rl.White,
		)
	} else {
		s.drawBlockCorners(entity, position, sprite, block)
	}
}

func (s *Renderer) drawBlockCorners(entity *ecs.Entity,
	position *components.Position,
	sprite *components.Sprite,
	block *components.BlockCorners) {

	row := float32(0.0)
	cornerWidth := s.stage.BlockWidth / 2
	cornerHeight := s.stage.BlockHeight / 2
	x := float32(position.CurrentBlockPosition.X) * s.stage.BlockWidth
	y := float32(position.CurrentBlockPosition.Y) * s.stage.BlockHeight

	rl.DrawTextureRec(*sprite.Texture, rl.NewRectangle(float32(block.Corners[game.CornerOverLeft])*cornerWidth, row*cornerHeight, cornerWidth, cornerHeight), rl.NewVector2(x, y), rl.White)
	rl.DrawTextureRec(*sprite.Texture, rl.NewRectangle(float32(block.Corners[game.CornerOverRight])*cornerWidth, row*cornerHeight, cornerWidth, cornerHeight), rl.NewVector2(x+cornerWidth, y), rl.White)
	rl.DrawTextureRec(*sprite.Texture, rl.NewRectangle(float32(block.Corners[game.CornerUnderRight])*cornerWidth, row*cornerHeight, cornerWidth, cornerHeight), rl.NewVector2(x+cornerWidth, y+cornerHeight), rl.White)
	rl.DrawTextureRec(*sprite.Texture, rl.NewRectangle(float32(block.Corners[game.CornerUnderLeft])*cornerWidth, row*cornerHeight, cornerWidth, cornerHeight), rl.NewVector2(x, y+cornerHeight), rl.White)
}

func (s *Renderer) drawParticle(particle *components.Particle) {
	src := rl.NewRectangle(0, 0, float32(s.mist.Width), float32(s.mist.Height))
	for _, p := range particle.Parts {
		a := p.Alpha * (0.75 + 0.25*float32(math.Sin(s.elapsedTime*0.5+float64(p.Phase))))
		c := particle.Tint
		c.A = uint8(a * 255)
		dst := rl.NewRectangle(p.X, p.Y, p.Size, p.Size*0.6)
		rl.DrawTexturePro(s.mist, src, dst, rl.NewVector2(p.Size/2, p.Size*0.3), 0, c)
	}
}

func (s *Renderer) drawTimer() {
	if s.stage.TimeLimit <= 0 {
		return
	}

	seconds := int(s.stage.TimeRemaining + 0.99)
	timerText := fmt.Sprintf("%02d:%02d", seconds/60, seconds%60)
	fontSize := int32(32)
	padding := int32(10)
	x := int32(rl.GetScreenWidth()) - rl.MeasureText(timerText, fontSize) - padding
	color := rl.White
	if seconds <= 10 {
		color = rl.Red
	}
	rl.DrawText(timerText, x, padding, fontSize, color)
}

func (s *Renderer) drawTimeExpired() {
	message := "TIME OUT"
	fontSize := int32(56)
	x := (int32(rl.GetScreenWidth()) - rl.MeasureText(message, fontSize)) / 2
	y := int32(rl.GetScreenHeight()) / 2
	rl.DrawText(message, x, y, fontSize, rl.Red)
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
