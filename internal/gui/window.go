package gui

import rl "github.com/gen2brain/raylib-go/raylib"

// Window is a reusable draggable modal frame: dimmed backdrop, shadow, panel and title bar.
type Window struct {
	Bounds     rl.Rectangle
	Title      string
	dragging   bool
	dragOffset rl.Vector2
}

func NewWindow(title string, w, h float32) Window {
	return Window{Title: title, Bounds: rl.NewRectangle((ScreenW-w)/2, (ScreenH-h)/2, w, h)}
}

func (w *Window) titleBar() rl.Rectangle {
	return rl.NewRectangle(w.Bounds.X, w.Bounds.Y, w.Bounds.Width, 30)
}

// UpdateDrag moves the window while its title bar is dragged.
func (w *Window) UpdateDrag(can bool) {
	m := rl.GetMousePosition()
	if can && rl.IsMouseButtonPressed(rl.MouseButtonLeft) && rl.CheckCollisionPointRec(m, w.titleBar()) {
		w.dragging = true
		w.dragOffset = rl.NewVector2(m.X-w.Bounds.X, m.Y-w.Bounds.Y)
	}
	if w.dragging {
		if rl.IsMouseButtonDown(rl.MouseButtonLeft) {
			w.Bounds.X = Clamp(m.X-w.dragOffset.X, 0, ScreenW-w.Bounds.Width)
			w.Bounds.Y = Clamp(m.Y-w.dragOffset.Y, 0, ScreenH-w.Bounds.Height)
		} else {
			w.dragging = false
		}
	}
}

func (w *Window) DrawFrame() {
	b := w.Bounds
	rl.DrawRectangle(0, 0, int32(ScreenW), int32(ScreenH), rl.Fade(rl.Black, 0.55))
	rl.DrawRectangleRounded(rl.NewRectangle(b.X+4, b.Y+6, b.Width, b.Height), 0.03, 8, rl.Fade(rl.Black, 0.35))
	rl.DrawRectangleRounded(b, 0.03, 8, ColPanel)
	rl.DrawRectangleRounded(w.titleBar(), 0.4, 8, ColTitle)
	rl.DrawRectangle(int32(b.X), int32(b.Y)+18, int32(b.Width), 12, ColTitle)
	rl.DrawText(w.Title, int32(b.X)+12, int32(b.Y)+6, 18, ColText)
}

// EntityDialog is a modal window that hosts an ItemGrid for choosing an entity.
type EntityDialog struct {
	Window
	Active      bool
	Grid        *ItemGrid
	okBtn       *Button
	cancelBtn   *Button
	result      DialogResult
	openedFrame uint64
}

func (d *EntityDialog) Open(items []Entity, selected int, tint rl.Color) {
	*d = EntityDialog{Window: NewWindow("Select entity", 500, 470), Active: true, openedFrame: Frame}
	d.Grid = NewItemGrid(items)
	d.Grid.Selected = selected
	d.Grid.Tint = tint
	d.Grid.OnActivate = func(int) { d.result = ResOK }
	d.okBtn = NewButton("Select", "", func() { d.result = ResOK })
	d.cancelBtn = NewButton("Cancel", "", func() { d.result = ResCancel })
}

// UpdateDraw returns the selected index and a result once the dialog closes.
func (d *EntityDialog) UpdateDraw() (int, DialogResult) {
	if !d.Active {
		return -1, ResNone
	}
	can := d.openedFrame != Frame // ignore the click that opened it
	d.UpdateDrag(can)
	d.DrawFrame()

	b := d.Bounds
	d.result = ResNone
	d.Grid.Bounds = rl.NewRectangle(b.X+16, b.Y+44, b.Width-32, b.Height-44-76)
	d.Grid.UpdateDraw(can)

	// footer: selected name + buttons
	const bw, bh = 96, 34
	by := b.Y + b.Height - bh - 16
	label := "Nothing selected"
	if sel := d.Grid.Selected; sel >= 0 && sel < len(d.Grid.Items) {
		label = d.Grid.Items[sel].Name
	}
	rl.DrawText(FitText(label, int32(b.Width)-32-2*bw-30, 16), int32(b.X)+16, int32(by)+9, 16, ColTextDim)

	rightX := b.X + b.Width - 16
	d.okBtn.Bounds = rl.NewRectangle(rightX-bw, by, bw, bh)
	d.okBtn.Enabled = can && d.Grid.Selected >= 0
	d.okBtn.UpdateDraw()
	d.cancelBtn.Bounds = rl.NewRectangle(rightX-2*bw-10, by, bw, bh)
	d.cancelBtn.Enabled = can
	d.cancelBtn.UpdateDraw()

	if can {
		if rl.IsKeyPressed(rl.KeyEnter) && d.Grid.Selected >= 0 {
			d.result = ResOK
		}
		if rl.IsKeyPressed(rl.KeyEscape) {
			d.result = ResCancel
		}
	}
	if d.result != ResNone {
		d.Active = false
	}
	return d.Grid.Selected, d.result
}
