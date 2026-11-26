package renderer

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/tylland/dashland/cmd/stagemaker/tool"
	"github.com/tylland/dashland/internal/components"
	"github.com/tylland/dashland/internal/ecs"
	"github.com/tylland/dashland/internal/game"
	"github.com/tylland/dashland/internal/game/types"
)

var ColHoovered = rl.NewColor(255, 255, 0, 125)
var ColSelected = rl.NewColor(0, 255, 0, 125)

type EditorRenderer struct {
	ws          tool.ToolbarWorkspace
	camera      game.Camera
	elapsedTime float64
	mist        rl.Texture2D
}

func NewEditorRenderer(ws tool.ToolbarWorkspace, camera game.Camera) *EditorRenderer {
	return &EditorRenderer{ws: ws, camera: camera}
}

func (s *EditorRenderer) Update(world *ecs.World, deltaTime float32) {
	s.elapsedTime += float64(deltaTime)

	for _, entity := range world.Entities() {
		if entity.Category == types.EntityCategoryBlock {
			position := ecs.GetComponent[components.PositionComponent](entity)
			sprite := ecs.GetComponent[components.SpriteComponent](entity)

			if position != nil && sprite != nil {
				s.drawBlock(entity, position, sprite)
			}
		}
	}

	for _, entity := range world.Entities() {
		if entity.Category == types.EntityCategoryBlock {
			continue
		}

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

	s.drawGrid()
	s.drawHoover()
	s.drawSelected()
}

func (s *EditorRenderer) RunWhenPaused() bool {
	return true
}

func (s *EditorRenderer) drawBlock(entity *ecs.Entity, position *components.PositionComponent, sprite *components.SpriteComponent) {
	block := ecs.GetComponent[components.BlockComponent](entity)

	if !block.HasCorners {
		rl.DrawTextureRec(
			*sprite.Texture,
			rl.NewRectangle(float32(entity.Type)*game.BlockWidth, 0, game.BlockWidth, game.BlockHeight),
			position.Vector2,
			rl.White,
		)
	} else {
		s.drawBlockCorners(entity, position, sprite, block)
	}

	// switch entity.Type {
	// case types.EntitySoil:
	// 	s.drawBlockCorners(entity, position, 1)
	// case types.EntityVoid:
	// 	s.drawBlockCorners(entity, position, 0)
	// default:
	// 	blockTextures := s.stage.BlockMap.GetBlockTextures()
	// 	if blockTextures != nil {
	// 	}
	// }
}

func (s *EditorRenderer) drawBlockCorners(entity *ecs.Entity,
	position *components.PositionComponent,
	sprite *components.SpriteComponent,
	block *components.BlockComponent) {

	row := float32(0.0)
	cornerWidth := game.BlockWidth / 2
	cornerHeight := game.BlockHeight / 2
	x := float32(position.CurrentBlockPosition.X) * game.BlockWidth
	y := float32(position.CurrentBlockPosition.Y) * game.BlockHeight

	rl.DrawTextureRec(*sprite.Texture, rl.NewRectangle(float32(block.Corners[game.CornerOverLeft])*cornerWidth, row*cornerHeight, cornerWidth, cornerHeight), rl.NewVector2(x, y), rl.White)
	rl.DrawTextureRec(*sprite.Texture, rl.NewRectangle(float32(block.Corners[game.CornerOverRight])*cornerWidth, row*cornerHeight, cornerWidth, cornerHeight), rl.NewVector2(x+cornerWidth, y), rl.White)
	rl.DrawTextureRec(*sprite.Texture, rl.NewRectangle(float32(block.Corners[game.CornerUnderRight])*cornerWidth, row*cornerHeight, cornerWidth, cornerHeight), rl.NewVector2(x+cornerWidth, y+cornerHeight), rl.White)
	rl.DrawTextureRec(*sprite.Texture, rl.NewRectangle(float32(block.Corners[game.CornerUnderLeft])*cornerWidth, row*cornerHeight, cornerWidth, cornerHeight), rl.NewVector2(x, y+cornerHeight), rl.White)
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

func (s *EditorRenderer) drawHoover() {
	if target := s.ws.HoveredTarget(); target != nil {
		area := s.ws.Stage().Area(target.Position)
		rl.DrawRectangleLinesEx(area, 2/s.ws.View().Camera.Zoom, ColHoovered)
	}
}

func (s *EditorRenderer) drawSelected() {
	if target := s.ws.SelectedTarget(); target != nil {
		area := s.ws.Stage().Area(target.Position)
		rl.DrawRectangleLinesEx(area, 2/s.ws.View().Camera.Zoom, ColSelected)
	}
}

// drawGrid draws an adaptive background grid (call between Begin and End).
func (s *EditorRenderer) drawGrid() {

	stage := s.ws.Stage()
	zoom := s.camera.GetCamera().Zoom
	thick := 1 / zoom // constant 1px on screen
	line := rl.Fade(rl.White, 0.07)
	axis := rl.Fade(rl.NewColor(90, 150, 255, 255), 0.55)

	maxX := float32(stage.Width) * stage.BlockWidth
	maxY := float32(stage.Height) * stage.BlockHeight

	for x := float32(0); x <= maxX; x = x + stage.BlockWidth {
		c := line

		rl.DrawLineEx(rl.NewVector2(x, 0), rl.NewVector2(x, maxY), thick, c)
	}

	for y := float32(0); y <= maxY; y = y + stage.BlockWidth {
		c := line

		rl.DrawLineEx(rl.NewVector2(0, y), rl.NewVector2(maxX, y), thick, c)
	}

	area := rl.NewRectangle(0, 0, maxX, maxY)
	rl.DrawRectangleLinesEx(area, 2/zoom, axis)
}
