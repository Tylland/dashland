package systems

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/tylland/dashland/internal/components"
	"github.com/tylland/dashland/internal/ecs"
	"github.com/tylland/dashland/internal/game"
)

type EditorRenderer struct {
	stage       *game.Stage
	camera      game.Camera
	elapsedTime float64
	mist        rl.Texture2D
}

func NewEditorRenderer(stage *game.Stage, camera game.Camera) *EditorRenderer {
	return &EditorRenderer{stage: stage, camera: camera}
}

func (s *EditorRenderer) Update(world *ecs.World, deltaTime float32) {
	s.elapsedTime += float64(deltaTime)

	s.stage.Update(deltaTime)
	s.stage.Render(deltaTime)

	rl.DrawRectangle(
		0,
		0,
		int32(s.stage.Width*int(s.stage.BlockWidth)),
		int32(s.stage.Height*int(s.stage.BlockHeight)),
		rl.Gray,
	)

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

		particle := ecs.GetComponent[components.ParticleComponent](entity)

		if particle != nil {
			s.drawParticle(particle)
		}
	}
}

func (s *EditorRenderer) RunWhenPaused() bool {
	return true
}

func (s *EditorRenderer) drawParticle(particle *components.ParticleComponent) {
	src := rl.NewRectangle(0, 0, float32(s.mist.Width), float32(s.mist.Height))
	for _, p := range particle.Parts {
		a := p.Alpha * (0.75 + 0.25*float32(math.Sin(s.elapsedTime*0.5+float64(p.Phase))))
		c := particle.Tint
		c.A = uint8(a * 255)
		dst := rl.NewRectangle(p.X, p.Y, p.Size, p.Size*0.6)
		rl.DrawTexturePro(s.mist, src, dst, rl.NewVector2(p.Size/2, p.Size*0.3), 0, c)
	}
}
