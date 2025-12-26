package systems

import (
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
		position := ecs.GetComponent[components.PositionComponent](entity.Components)
		sprite := ecs.GetComponent[components.SpriteComponent](entity.Components)

		if position != nil && sprite != nil {
			// if entity.Collision != nil {
			// 	rl.DrawRectangle(int32(position.Vector2.X), int32(position.Vector2.Y), int32(entity.Collision.Width), int32(entity.Collision.Height), rl.Red)
			// }

			if sprite.FrameCount > 1 {
				sprite.FrameTimer += deltaTime
				if sprite.FrameTimer >= 1.0/sprite.FrameSpeed {
					sprite.UpdateFrame((sprite.Frame + 1) % sprite.FrameCount)

					sprite.FrameTimer = 0

				}
			}

			rl.DrawTextureRec(*sprite.Texture, sprite.Source, position.Vector2, rl.White)
		}
	}

	rl.EndMode2D()
	rl.EndDrawing()
}
