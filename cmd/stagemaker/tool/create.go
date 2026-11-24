package tool

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/tylland/dashland/internal/common"
	"github.com/tylland/dashland/internal/ecs"
	"github.com/tylland/dashland/internal/game"
	"github.com/tylland/dashland/internal/gui"
)

type ToolbarWorkspace interface {
	World() *ecs.World
	View() *gui.View
	Stage() *game.Stage
	StagePath() string
	Entities() []*gui.EntityItem
	SelectedEntity() *gui.EntityItem
	Inspection() *gui.EntityInspection
	Inspect(inspect *gui.EntityInspection)
	HoveredTarget() *Target
	SetHoveredTarget(target *Target)
	SelectedTarget() *Target
	SetSelectedTarget(target *Target)
}

// ------------------------------------------------------------------ Rectangle grid tool

const cellSize = 32 // world units per grid cell; entities are drawn cellSize x cellSize

// FillRectTool fills a dragged rectangle of grid cells with the selected entity.
// Cells are aligned to the world grid, so repeated fills line up exactly.
type FillRectTool struct {
	active   bool // tool is selected
	dragging bool
	start    common.BlockPosition // cell coordinates where the drag began
	end      common.BlockPosition // cell coordinates where the drag began
	stage    *game.Stage
	ws       ToolbarWorkspace
}

func NewFillRectTool(stage *game.Stage) *FillRectTool {
	return &FillRectTool{stage: stage}
}

// Lifetime
func (t *FillRectTool) IsActive() bool {
	return t.active
}

func (t *FillRectTool) Activate(tb gui.ToolBar) {
	t.active = true

	if ws, ok := tb.Workspace().(ToolbarWorkspace); ok {
		t.ws = ws
		t.ws.SetHoveredTarget(nil)
	}
}

func (t *FillRectTool) Deactivate() {
	t.active = false
}

func (t *FillRectTool) Cancel() {
	t.dragging = false
}

// Mouse Tracking
func (t *FillRectTool) WorldMouseDown(world rl.Vector2, button rl.MouseButton) {
	t.dragging = true

	t.start, _ = t.stage.ClampedBlockPosition(world)

	t.end = t.start
}

func (t *FillRectTool) WorldMouseMove(world rl.Vector2) {
	if t.dragging {
		//		t.endX, t.endY = cellOf(world)

		t.end, _ = t.stage.ClampedBlockPosition(world)

		t.DrawPreview(t.ws.SelectedEntity(), rl.White, t.ws.View().Camera.Zoom)
	}
}

func (t *FillRectTool) WorldMouseUp(world rl.Vector2, button rl.MouseButton) {
	if t.dragging {
		//		t.endX, t.endY = cellOf(world)
		t.end, _ = t.stage.ClampedBlockPosition(world)
		t.CreateEntities()
	}

	t.Cancel()
}

func (t *FillRectTool) ScreenMouseMove(coord rl.Vector2) {
	if t.dragging {
		t.DrawLabel(coord)
	}
}

// cellOf returns the grid cell containing a world position.
func cellOf(world rl.Vector2) (int, int) {
	return int(math.Floor(float64(world.X / cellSize))), int(math.Floor(float64(world.Y / cellSize)))
}

// Cells returns the dragged range as inclusive, normalized cell coordinates.
func (t *FillRectTool) Cells() (x0, y0, x1, y1 int) {
	x0, x1 = t.start.X, t.end.X
	y0, y1 = t.start.Y, t.end.Y
	if x0 > x1 {
		x0, x1 = x1, x0
	}
	if y0 > y1 {
		y0, y1 = y1, y0
	}
	return
}

// DrawPreview draws the selection and ghost entities (call in world space).
func (t *FillRectTool) DrawPreview(e *gui.EntityItem, tint rl.Color, zoom float32) {
	if !t.dragging {
		return
	}
	x0, y0, x1, y1 := t.Cells()
	area := rl.NewRectangle(float32(x0)*cellSize, float32(y0)*cellSize,
		float32(x1-x0+1)*cellSize, float32(y1-y0+1)*cellSize)
	rl.DrawRectangleRec(area, rl.Fade(gui.ColAccent, 0.12))

	ghost := rl.Fade(tint, 0.55)
	for cy := y0; cy <= y1; cy++ {
		for cx := x0; cx <= x1; cx++ {
			e.Draw(rl.NewRectangle(float32(cx)*cellSize, float32(cy)*cellSize, cellSize, cellSize), ghost)
		}
	}
	rl.DrawRectangleLinesEx(area, 2/zoom, gui.ColAccent)
}

// DrawLabel shows "W x H = N" next to the mouse (call in screen space).
func (t *FillRectTool) DrawLabel(mouse rl.Vector2) {
	if !t.dragging {
		return
	}
	x0, y0, x1, y1 := t.Cells()
	w, h := x1-x0+1, y1-y0+1
	txt := fmt.Sprintf("%d x %d = %d", w, h, w*h)
	col := gui.ColText

	r := rl.NewRectangle(mouse.X+16, mouse.Y+16, float32(rl.MeasureText(txt, 16))+16, 26)
	r.X = gui.Clamp(r.X, 4, gui.ScreenW-r.Width-4)
	r.Y = gui.Clamp(r.Y, 4, gui.ScreenH-r.Height-4)
	rl.DrawRectangleRounded(r, 0.3, 6, rl.NewColor(20, 20, 26, 230))
	rl.DrawText(txt, int32(r.X)+8, int32(r.Y)+5, 16, col)
}

func (t *FillRectTool) CreateEntities() {

	x0, y0, x1, y1 := t.Cells()

	for cy := y0; cy <= y1; cy++ {
		for cx := x0; cx <= x1; cx++ {
			pos := common.NewBlockPosition(cx, cy)

			entity := game.Must(game.NewEntityAndHeader(t.ws.SelectedEntity().Type, t.ws.World(), t.ws.Stage(), pos))

			if existing, ok := t.ws.Stage().GetEntityAtPosition(pos); ok {
				t.ws.World().RemoveEntity(existing)
			}

			t.ws.Stage().SetEntity(entity, pos)
		}
	}

}
