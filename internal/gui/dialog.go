package gui

import rl "github.com/gen2brain/raylib-go/raylib"

type DialogKind int

const (
	DlgNone DialogKind = iota
	DlgMessage
	DlgConfirm
	DlgInput
)

type DialogAction int

const (
	ActNone DialogAction = iota
	ActAbout
	ActClear
	ActRename
)

type DialogResult int

const (
	ResNone DialogResult = iota
	ResOK
	ResCancel
)

type Dialog struct {
	Kind        DialogKind
	Action      DialogAction
	Bounds      rl.Rectangle
	Title       string
	Message     string
	Input       []rune
	dragging    bool
	dragOffset  rl.Vector2
	openedFrame uint64
	backHold    float32 // backspace auto-repeat timer

	okBtn     *Button
	cancelBtn *Button // nil for message dialogs
	pending   DialogResult
}

var Frame uint64

func (d *Dialog) Open(kind DialogKind, action DialogAction, title, message, initial string) {
	*d = Dialog{
		Kind:        kind,
		Action:      action,
		Title:       title,
		Message:     message,
		Input:       []rune(initial),
		openedFrame: Frame,
	}
	w, h := float32(420), float32(170)
	if kind == DlgInput {
		h = 200
	}
	d.Bounds = rl.NewRectangle((ScreenW-w)/2, (ScreenH-h)/2, w, h)

	switch kind {
	case DlgConfirm:
		d.okBtn = NewButton("Yes", "", func() { d.pending = ResOK })
		d.cancelBtn = NewButton("No", "", func() { d.pending = ResCancel })
	case DlgInput:
		d.okBtn = NewButton("OK", "", func() { d.pending = ResOK })
		d.cancelBtn = NewButton("Cancel", "", func() { d.pending = ResCancel })
	default:
		d.okBtn = NewButton("OK", "", func() { d.pending = ResOK })
	}
	for rl.GetCharPressed() > 0 { // flush pending characters
	}
}

// UpdateDraw handles input and draws the dialog. It returns a result when the dialog should close.
func (d *Dialog) UpdateDraw() DialogResult {
	if d.Kind == DlgNone {
		return ResNone
	}

	canInteract := d.openedFrame != Frame // ignore the click that opened it
	m := rl.GetMousePosition()

	// --- dragging via title bar
	titleBar := rl.NewRectangle(d.Bounds.X, d.Bounds.Y, d.Bounds.Width, 30)
	if canInteract && rl.IsMouseButtonPressed(rl.MouseButtonLeft) && rl.CheckCollisionPointRec(m, titleBar) {
		d.dragging = true
		d.dragOffset = rl.NewVector2(m.X-d.Bounds.X, m.Y-d.Bounds.Y)
	}
	if d.dragging {
		if rl.IsMouseButtonDown(rl.MouseButtonLeft) {
			d.Bounds.X = Clamp(m.X-d.dragOffset.X, 0, ScreenW-d.Bounds.Width)
			d.Bounds.Y = Clamp(m.Y-d.dragOffset.Y, 0, ScreenH-d.Bounds.Height)
		} else {
			d.dragging = false
		}
	}

	// --- backdrop + panel
	rl.DrawRectangle(0, 0, int32(ScreenW), int32(ScreenH), rl.Fade(rl.Black, 0.55))
	shadow := rl.NewRectangle(d.Bounds.X+4, d.Bounds.Y+6, d.Bounds.Width, d.Bounds.Height)
	rl.DrawRectangleRounded(shadow, 0.04, 8, rl.Fade(rl.Black, 0.35))
	rl.DrawRectangleRounded(d.Bounds, 0.04, 8, ColPanel)
	rl.DrawRectangleRounded(titleBar, 0.4, 8, ColTitle)
	rl.DrawRectangle(int32(d.Bounds.X), int32(d.Bounds.Y)+18, int32(d.Bounds.Width), 12, ColTitle)
	rl.DrawText(d.Title, int32(d.Bounds.X)+12, int32(d.Bounds.Y)+6, 18, ColText)

	// --- message
	rl.DrawText(d.Message, int32(d.Bounds.X)+20, int32(d.Bounds.Y)+46, 18, ColText)

	// --- text input field
	if d.Kind == DlgInput {
		box := rl.NewRectangle(d.Bounds.X+20, d.Bounds.Y+82, d.Bounds.Width-40, 34)
		border := rl.NewRectangle(box.X-1, box.Y-1, box.Width+2, box.Height+2)
		rl.DrawRectangleRounded(border, 0.2, 6, ColAccent)
		rl.DrawRectangleRounded(box, 0.2, 6, ColTitle)
		rl.DrawText(string(d.Input), int32(box.X)+10, int32(box.Y)+8, 18, ColText)

		if canInteract {
			for {
				ch := rl.GetCharPressed()
				if ch <= 0 {
					break
				}
				if ch >= 32 && ch <= 125 && len(d.Input) < MaxText {
					d.Input = append(d.Input, rune(ch))
				}
			}
			// backspace with auto-repeat
			if rl.IsKeyPressed(rl.KeyBackspace) {
				d.backHold = 0
				if len(d.Input) > 0 {
					d.Input = d.Input[:len(d.Input)-1]
				}
			} else if rl.IsKeyDown(rl.KeyBackspace) {
				d.backHold += rl.GetFrameTime()
				if d.backHold > 0.4 && len(d.Input) > 0 {
					d.Input = d.Input[:len(d.Input)-1]
					d.backHold = 0.36
				}
			}
		}
		if int(rl.GetTime()*2)%2 == 0 { // blinking caret
			cx := int32(box.X) + 10 + rl.MeasureText(string(d.Input), 18) + 2
			rl.DrawRectangle(cx, int32(box.Y)+7, 2, 20, ColText)
		}
	}

	// --- buttons
	const bw, bh = 96, 34
	by := d.Bounds.Y + d.Bounds.Height - bh - 16
	rightX := d.Bounds.X + d.Bounds.Width - 20

	d.pending = ResNone
	d.okBtn.Bounds = rl.NewRectangle(rightX-bw, by, bw, bh)
	d.okBtn.Enabled = canInteract
	d.okBtn.UpdateDraw()
	if d.cancelBtn != nil {
		d.cancelBtn.Bounds = rl.NewRectangle(rightX-2*bw-10, by, bw, bh)
		d.cancelBtn.Enabled = canInteract
		d.cancelBtn.UpdateDraw()
	}

	if canInteract {
		if rl.IsKeyPressed(rl.KeyEnter) {
			d.pending = ResOK
		}
		if rl.IsKeyPressed(rl.KeyEscape) {
			if d.cancelBtn != nil {
				d.pending = ResCancel
			} else {
				d.pending = ResOK
			}
		}
	}
	return d.pending
}
