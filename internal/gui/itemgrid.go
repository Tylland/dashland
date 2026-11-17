package gui

import rl "github.com/gen2brain/raylib-go/raylib"

const (
	gridCellW = 96
	gridCellH = 92
	gridGap   = 8
	gridIcon  = 56
)

// ItemGrid is a scrollable grid of entities; the user selects one cell.
type ItemGrid struct {
	Bounds     rl.Rectangle
	Items      []Entity
	Selected   int // -1 = nothing selected
	Scroll     int // first visible row
	Tint       rl.Color
	OnSelect   func(index int) // single click / keyboard change
	OnActivate func(index int) // double-click

	lastClickIdx  int
	lastClickTime float64
}

func NewItemGrid(items []Entity) *ItemGrid {
	return &ItemGrid{Items: items, Selected: -1, Tint: rl.White, lastClickIdx: -1}
}

func (g *ItemGrid) metrics() (cols, visRows int) {
	cols = int((g.Bounds.Width - gridGap) / (gridCellW + gridGap))
	visRows = int((g.Bounds.Height - gridGap) / (gridCellH + gridGap))
	if cols < 1 {
		cols = 1
	}
	if visRows < 1 {
		visRows = 1
	}
	return
}

func (g *ItemGrid) cellRect(row, col int, ox float32) rl.Rectangle {
	return rl.NewRectangle(
		g.Bounds.X+ox+float32(col)*(gridCellW+gridGap),
		g.Bounds.Y+gridGap+float32(row)*(gridCellH+gridGap),
		gridCellW, gridCellH)
}

// hit returns the index of the visible cell under p, or -1.
func (g *ItemGrid) hit(p rl.Vector2, cols, visRows int, ox float32) int {
	for row := 0; row < visRows; row++ {
		for col := 0; col < cols; col++ {
			i := (g.Scroll+row)*cols + col
			if i >= len(g.Items) {
				return -1
			}
			if rl.CheckCollisionPointRec(p, g.cellRect(row, col, ox)) {
				return i
			}
		}
	}
	return -1
}

func (g *ItemGrid) selectIndex(i int) {
	g.Selected = i
	if g.OnSelect != nil {
		g.OnSelect(i)
	}
}

func (g *ItemGrid) move(delta int) {
	switch {
	case g.Selected < 0:
		g.selectIndex(0)
	case g.Selected+delta >= 0 && g.Selected+delta < len(g.Items):
		g.selectIndex(g.Selected + delta)
	}
}

// UpdateDraw handles mouse/keyboard input (when enabled) and draws the grid inside g.Bounds.
func (g *ItemGrid) UpdateDraw(enabled bool) {
	m := rl.GetMousePosition()
	cols, visRows := g.metrics()
	totalRows := (len(g.Items) + cols - 1) / cols
	maxScroll := totalRows - visRows
	if maxScroll < 0 {
		maxScroll = 0
	}
	ox := (g.Bounds.Width - (float32(cols)*gridCellW + float32(cols-1)*gridGap)) / 2

	inside := rl.CheckCollisionPointRec(m, g.Bounds)
	hovered := -1
	if enabled && inside {
		hovered = g.hit(m, cols, visRows, ox)
	}

	// --- input
	if enabled {
		if inside {
			g.Scroll -= int(rl.GetMouseWheelMove())
		}
		if rl.IsKeyPressed(rl.KeyRight) {
			g.move(1)
		}
		if rl.IsKeyPressed(rl.KeyLeft) {
			g.move(-1)
		}
		if rl.IsKeyPressed(rl.KeyDown) {
			g.move(cols)
		}
		if rl.IsKeyPressed(rl.KeyUp) {
			g.move(-cols)
		}
		if hovered >= 0 && rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
			g.selectIndex(hovered)
			now := rl.GetTime()
			if hovered == g.lastClickIdx && now-g.lastClickTime < 0.35 { // double-click
				g.lastClickIdx = -1
				if g.OnActivate != nil {
					g.OnActivate(hovered)
				}
			} else {
				g.lastClickIdx, g.lastClickTime = hovered, now
			}
		}
		if g.Selected >= 0 { // keep the selection visible
			row := g.Selected / cols
			if row < g.Scroll {
				g.Scroll = row
			}
			if row >= g.Scroll+visRows {
				g.Scroll = row - visRows + 1
			}
		}
	}
	if g.Scroll > maxScroll {
		g.Scroll = maxScroll
	}
	if g.Scroll < 0 {
		g.Scroll = 0
	}

	// --- drawing (clipped to the grid)
	rl.DrawRectangleRounded(g.Bounds, 0.02, 4, ColTitle)
	rl.BeginScissorMode(int32(g.Bounds.X), int32(g.Bounds.Y), int32(g.Bounds.Width), int32(g.Bounds.Height))
	for row := 0; row < visRows; row++ {
		for col := 0; col < cols; col++ {
			i := (g.Scroll+row)*cols + col
			if i >= len(g.Items) {
				break
			}
			r := g.cellRect(row, col, ox)

			if i == g.Selected {
				rl.DrawRectangleRounded(rl.NewRectangle(r.X-2, r.Y-2, r.Width+4, r.Height+4), 0.12, 6, ColAccent)
			}
			bg := ColBtn
			if i == hovered {
				bg = ColBtnHover
			}
			rl.DrawRectangleRounded(r, 0.12, 6, bg)
			if i == g.Selected {
				rl.DrawRectangleRounded(r, 0.12, 6, rl.Fade(ColAccent, 0.35))
			}

			icon := rl.NewRectangle(r.X+(r.Width-gridIcon)/2, r.Y+8, gridIcon, gridIcon)
			g.Items[i].Draw(icon, g.Tint)

			name := g.Items[i].Name
			tw := rl.MeasureText(name, 14)
			rl.DrawText(name, int32(r.X)+(int32(r.Width)-tw)/2, int32(r.Y)+gridIcon+16, 14, ColText)
		}
	}
	rl.EndScissorMode()

	if maxScroll > 0 { // scrollbar
		thumbH := g.Bounds.Height * float32(visRows) / float32(totalRows)
		thumbY := g.Bounds.Y + (g.Bounds.Height-thumbH)*float32(g.Scroll)/float32(maxScroll)
		rl.DrawRectangleRounded(rl.NewRectangle(g.Bounds.X+g.Bounds.Width-7, thumbY, 4, thumbH), 1, 4, rl.Fade(rl.White, 0.3))
	}
}
