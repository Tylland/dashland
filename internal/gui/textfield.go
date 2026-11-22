package gui

import rl "github.com/gen2brain/raylib-go/raylib"

// ------------------------------------------------------------------ TextField

// TextField is a single-line text box. The owner decides which field has focus.
type TextField struct {
	Text     string
	Bounds   rl.Rectangle
	MaxLen   int
	Allow    func(r rune) bool      // which characters may be typed (nil = all)
	OnChange func(text string) bool // called after every edit; false = text is not valid
	Invalid  bool

	backHold float32 // backspace auto-repeat timer
}

// UpdateDraw handles typing (when focused) and draws the field. It reports a click on it.
// clip is the visible area the field lives in (clicks outside it are ignored).
func (f *TextField) UpdateDraw(enabled, focused bool, clip rl.Rectangle) (clicked bool) {
	m := rl.GetMousePosition()
	hover := enabled && rl.CheckCollisionPointRec(m, f.Bounds) && rl.CheckCollisionPointRec(m, clip)
	clicked = hover && rl.IsMouseButtonPressed(rl.MouseButtonLeft)

	if focused && enabled {
		changed := false
		for {
			ch := rl.GetCharPressed()
			if ch <= 0 {
				break
			}
			r := rune(ch)
			if (f.Allow == nil || f.Allow(r)) && len([]rune(f.Text)) < f.MaxLen {
				f.Text += string(r)
				changed = true
			}
		}
		del := false
		if rl.IsKeyPressed(rl.KeyBackspace) {
			f.backHold, del = 0, true
		} else if rl.IsKeyDown(rl.KeyBackspace) {
			f.backHold += rl.GetFrameTime()
			if f.backHold > 0.4 {
				f.backHold, del = 0.36, true
			}
		}
		if del {
			if r := []rune(f.Text); len(r) > 0 {
				f.Text = string(r[:len(r)-1])
				changed = true
			}
		}
		if changed && f.OnChange != nil {
			f.Invalid = !f.OnChange(f.Text)
		}
	}

	border := ColBtn
	switch {
	case f.Invalid:
		border = rl.NewColor(230, 90, 90, 255)
	case focused:
		border = ColAccent
	case hover:
		border = ColBtnHover
	}
	b := f.Bounds
	rl.DrawRectangleRounded(rl.NewRectangle(b.X-1, b.Y-1, b.Width+2, b.Height+2), 0.2, 4, border)
	rl.DrawRectangleRounded(b, 0.2, 4, rl.NewColor(30, 32, 40, 255))
	txt := FitText(f.Text, int32(b.Width)-16, 16)
	rl.DrawText(txt, int32(b.X)+8, int32(b.Y)+5, 16, ColText)
	if focused && int(rl.GetTime()*2)%2 == 0 { // blinking caret
		rl.DrawRectangle(int32(b.X)+8+rl.MeasureText(txt, 16)+1, int32(b.Y)+4, 2, 18, ColText)
	}
	return clicked
}
