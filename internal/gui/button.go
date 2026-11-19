package gui

import rl "github.com/gen2brain/raylib-go/raylib"

// Button is a reusable clickable widget. The owner positions it (Bounds / Enabled)
// and calls UpdateDraw once per frame; the button runs OnClick when it is clicked.
type Button struct {
	Label    string
	Tooltip  string
	Bounds   rl.Rectangle
	Enabled  bool
	Active   bool            // toggled-on look (e.g. the selected tool)
	GapAfter float32         // extra layout space after this button (used by Toolbar)
	OnClick  func()          // called when the button is clicked
	Decorate func(b *Button) // optional hook to draw extra content (icons, swatches...)

	hovered bool
	pressed bool
}

func NewButton(label, tooltip string, onClick func()) *Button {
	return &Button{Label: label, Tooltip: tooltip, Enabled: true, OnClick: onClick}
}

func (b *Button) Hovered() bool { return b.hovered }
func (b *Button) Pressed() bool { return b.pressed }

// Update refreshes the hover/pressed state and fires OnClick. It reports whether the button was clicked.
func (b *Button) Update() bool {
	b.hovered = b.Enabled && rl.CheckCollisionPointRec(rl.GetMousePosition(), b.Bounds)
	b.pressed = b.hovered && rl.IsMouseButtonDown(rl.MouseButtonLeft)
	if b.hovered && rl.IsMouseButtonReleased(rl.MouseButtonLeft) {
		if b.OnClick != nil {
			b.OnClick()
		}
		return true
	}
	return false
}

func (b *Button) Draw() {
	bg := ColBtn
	switch {
	case !b.Enabled:
		bg = ColBtnOff
	case b.pressed:
		bg = ColBtnDown
	case b.hovered:
		bg = ColBtnHover
	}
	rl.DrawRectangleRounded(b.Bounds, 0.25, 6, bg)

	if b.Active {
		rl.DrawRectangleRounded(b.Bounds, 0.25, 6, rl.Fade(ColAccent, 0.55))
	}

	fg := ColText
	if !b.Enabled {
		fg = ColTextDim
	}
	tw := rl.MeasureText(b.Label, 18)
	rl.DrawText(b.Label,
		int32(b.Bounds.X)+(int32(b.Bounds.Width)-tw)/2,
		int32(b.Bounds.Y)+(int32(b.Bounds.Height)-18)/2, 18, fg)

	if b.Decorate != nil {
		b.Decorate(b)
	}
}

func (b *Button) UpdateDraw() bool {
	clicked := b.Update()
	b.Draw()
	return clicked
}
