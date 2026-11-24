package gui

import (
	"sync"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	GripSize   = 18 // thickness of the drag handle
	TbPad      = 8
	TbGap      = 6
	TbBtnW     = 96
	TbBtnH     = 30
	TbSepExtra = 10 // handy value for Button.GapAfter (group separator)
	SnapDist   = 16 // magnetic distance to screen edges
	SnapMargin = 8  // resting distance from the edge once snapped
	StatusH    = 24 // height of the status bar at the bottom
)

// Toolbar is a free-floating, draggable panel that owns and lays out a list of Buttons.
type Toolbar struct {
	Pos        rl.Vector2
	Vertical   bool
	Buttons    []*Button
	dragging   bool
	dragOffset rl.Vector2
	lastClick  float64
	lastHover  *Button
	hoverTimer float32
	ActiveTool Tool
	workspace  any
	mutex      sync.Mutex
}

func NewToolbar(pos rl.Vector2, workspace any) *Toolbar {
	return &Toolbar{Pos: pos, workspace: workspace, mutex: sync.Mutex{}}
}

func (t *Toolbar) ActivateTool(tool Tool) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	for _, b := range t.Buttons {
		b.Active = false
	}

	t.ActiveTool = tool
	tool.Activate(t)
}

func (t *Toolbar) Workspace() any {
	return t.workspace
}

// Add appends a button and returns it so callers can tweak it further.
func (t *Toolbar) Add(b *Button) *Button {
	t.Buttons = append(t.Buttons, b)
	return b
}

// Flip switches between horizontal and vertical layout.
func (t *Toolbar) Flip() {
	t.Vertical = !t.Vertical
	t.ClampToScreen()
}

func (t *Toolbar) size() (w, h float32) {
	n := float32(len(t.Buttons))
	var extra float32
	for _, b := range t.Buttons {
		extra += b.GapAfter
	}
	gaps := n - 1
	if gaps < 0 {
		gaps = 0
	}
	if t.Vertical {
		span := n*TbBtnH + gaps*TbGap + extra
		return TbBtnW + 2*TbPad, GripSize + TbPad + span + TbPad
	}
	span := n*TbBtnW + gaps*TbGap + extra
	return GripSize + TbPad + span + TbPad, TbBtnH + 2*TbPad
}

func (t *Toolbar) Bounds() rl.Rectangle {
	w, h := t.size()
	return rl.NewRectangle(t.Pos.X, t.Pos.Y, w, h)
}

func (t *Toolbar) gripRect() rl.Rectangle {
	w, h := t.size()
	if t.Vertical {
		return rl.NewRectangle(t.Pos.X, t.Pos.Y, w, GripSize)
	}
	return rl.NewRectangle(t.Pos.X, t.Pos.Y, GripSize, h)
}

// clampToScreen keeps the whole toolbar inside the window (above the status bar).
func (t *Toolbar) ClampToScreen() {
	w, h := t.size()
	t.Pos.X = Clamp(t.Pos.X, 0, ScreenW-w)
	t.Pos.Y = Clamp(t.Pos.Y, 0, ScreenH-StatusH-h)
}

// snap pulls the toolbar to a screen edge when it is close to one.
func (t *Toolbar) snap() {
	w, h := t.size()
	if t.Pos.X < SnapDist {
		t.Pos.X = SnapMargin
	}
	if t.Pos.X+w > ScreenW-SnapDist {
		t.Pos.X = ScreenW - w - SnapMargin
	}
	if t.Pos.Y < SnapDist {
		t.Pos.Y = SnapMargin
	}
	if t.Pos.Y+h > ScreenH-StatusH-SnapDist {
		t.Pos.Y = ScreenH - StatusH - h - SnapMargin
	}
}

// UpdateDraw handles dragging, lays out and updates every button, and draws the tooltip.
func (t *Toolbar) UpdateDraw(enabled bool) {
	m := rl.GetMousePosition()
	if !enabled {
		t.dragging = false
	}

	// --- grip: drag to move, double-click to flip orientation
	overGrip := enabled && rl.CheckCollisionPointRec(m, t.gripRect())
	if overGrip && rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
		now := rl.GetTime()
		if now-t.lastClick < 0.3 {
			t.Flip()
			t.dragging = false
			t.lastClick = 0
		} else {
			t.dragging = true
			t.dragOffset = rl.NewVector2(m.X-t.Pos.X, m.Y-t.Pos.Y)
			t.lastClick = now
		}
	}
	if t.dragging {
		if rl.IsMouseButtonDown(rl.MouseButtonLeft) {
			t.Pos = rl.NewVector2(m.X-t.dragOffset.X, m.Y-t.dragOffset.Y)
			t.ClampToScreen()
			t.snap()
		} else {
			t.dragging = false
		}
	}
	overToolbar := enabled && rl.CheckCollisionPointRec(m, t.Bounds())
	if t.dragging || overGrip {
		rl.SetMouseCursor(rl.MouseCursorResizeAll)
	} else if overToolbar {
		rl.SetMouseCursor(rl.MouseCursorDefault)
	}

	// --- panel
	b := t.Bounds()
	rl.DrawRectangleRounded(rl.NewRectangle(b.X+3, b.Y+5, b.Width, b.Height), 0.25, 8, rl.Fade(rl.Black, 0.40))
	rl.DrawRectangleRounded(rl.NewRectangle(b.X-1, b.Y-1, b.Width+2, b.Height+2), 0.25, 8, rl.Fade(rl.White, 0.12))
	rl.DrawRectangleRounded(b, 0.25, 8, ColToolbar)

	// grip dots
	dotCol := ColTextDim
	if t.dragging || overGrip {
		dotCol = ColAccent
	}
	for i := -1; i <= 1; i++ {
		if t.Vertical {
			rl.DrawCircle(int32(b.X+b.Width/2)+int32(i*8), int32(b.Y+GripSize/2), 2, dotCol)
		} else {
			rl.DrawCircle(int32(b.X+GripSize/2), int32(b.Y+b.Height/2)+int32(i*8), 2, dotCol)
		}
	}

	// --- buttons: the toolbar only positions them, each Button handles itself
	vertical := t.Vertical // a button callback may flip the layout mid-loop
	x, y := b.X+TbPad, b.Y+TbPad
	if vertical {
		y += GripSize
	} else {
		x += GripSize
	}

	var hovered *Button
	for _, btn := range t.Buttons {
		btn.Bounds = rl.NewRectangle(x, y, TbBtnW, TbBtnH)
		btn.Enabled = enabled && !t.dragging
		btn.UpdateDraw()
		if btn.Hovered() {
			hovered = btn
		}
		if vertical {
			y += TbBtnH + TbGap + btn.GapAfter
		} else {
			x += TbBtnW + TbGap + btn.GapAfter
		}
	}

	// --- tooltip after a short delay
	if hovered != t.lastHover {
		t.lastHover = hovered
		t.hoverTimer = 0
	} else if hovered != nil {
		t.hoverTimer += rl.GetFrameTime()
	}
	if hovered != nil && hovered.Tooltip != "" && t.hoverTimer > 0.5 {
		tw := float32(rl.MeasureText(hovered.Tooltip, 16))
		tr := rl.NewRectangle(m.X+12, m.Y+22, tw+16, 26)
		tr.X = Clamp(tr.X, 4, ScreenW-tr.Width-4)
		tr.Y = Clamp(tr.Y, 4, ScreenH-tr.Height-4)
		rl.DrawRectangleRounded(tr, 0.3, 6, rl.NewColor(20, 20, 26, 240))
		rl.DrawText(hovered.Tooltip, int32(tr.X)+8, int32(tr.Y)+5, 16, ColText)
	}
}
