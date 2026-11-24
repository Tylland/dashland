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

	// rl.DrawRectangle(
	// 	0,
	// 	0,
	// 	int32(s.stage.Width*int(s.stage.BlockWidth)),
	// 	int32(s.stage.Height*int(s.stage.BlockHeight)),
	// 	rl.Gray,
	// )

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

	s.DrawGrid()
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

// DrawGrid draws an adaptive background grid (call between Begin and End).
func (s *EditorRenderer) DrawGrid() {

	zoom := s.camera.GetCamera().Zoom
	thick := 1 / zoom // constant 1px on screen
	line := rl.Fade(rl.White, 0.07)
	axis := rl.Fade(rl.NewColor(90, 150, 255, 255), 0.55)

	maxX := float32(s.stage.Width) * s.stage.BlockWidth
	maxY := float32(s.stage.Height) * s.stage.BlockHeight

	for x := float32(0); x <= maxX; x = x + s.stage.BlockWidth {
		c := line

		rl.DrawLineEx(rl.NewVector2(x, 0), rl.NewVector2(x, maxY), thick, c)
	}

	for y := float32(0); y <= maxY; y = y + s.stage.BlockWidth {
		c := line

		rl.DrawLineEx(rl.NewVector2(0, y), rl.NewVector2(maxX, y), thick, c)
	}

	area := rl.NewRectangle(0, 0, maxX, maxY)
	rl.DrawRectangleLinesEx(area, 2/zoom, axis)
}
