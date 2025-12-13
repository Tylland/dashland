package game

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type RenderSystem struct {
	world *World
}

func NewRenderSystem(w *World) *RenderSystem {
	return &RenderSystem{
		world: w,
	}
}

func (s *RenderSystem) Update(deltaTime float32) {

	for _, entity := range s.world.entities {
		if entity != nil && entity.Position != nil && entity.Sprite != nil {
			position := entity.Position
			sprite := entity.Sprite

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
}
