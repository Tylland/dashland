package gui

import (
	"path/filepath"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/tylland/dashland/internal/ecs"
)

// ------------------------------------------------------------------ ValuesDialog

const (
	vdRowH    = 32
	vdHeaderH = 40
)

type valueRow struct {
	header         string // non-empty for a "# label" row
	val            *Value
	y, h           float32
	field1, field2 *TextField // text / decimal field, hex field
}

type fieldRef struct {
	field *TextField
	row   *valueRow
}

// ValuesDialog is a modal editor for a Document loaded from a file.
//   - click a field to edit it; Tab / Shift+Tab / Enter move between fields
//   - numbers have a decimal and a hex field that stay in sync
//   - editing a value marks it with '*' (not default)
//   - Save writes the file, Esc / Cancel discards the changes
type ValuesDialog struct {
	Window
	Active      bool
	Path        string
	entity      *ecs.Entity
	Doc         *Document
	rows        []*valueRow
	fields      []fieldRef
	contentH    float32
	listH       float32
	scroll      float32
	focus       *TextField
	saveBtn     *Button
	cancelBtn   *Button
	result      DialogResult
	err         string
	openedFrame uint64
}

func (d *ValuesDialog) Open(entity *ecs.Entity, doc *Document, path string) error {
	title := "Edit values: " + FitText(filepath.Base(path), 360, 18)
	*d = ValuesDialog{
		Window:      NewWindow(title, 580, 540),
		Active:      true,
		Path:        path,
		entity:      entity,
		Doc:         doc,
		openedFrame: Frame,
	}
	d.buildRows()
	d.saveBtn = NewButton("Save", "", func() { d.save() })
	d.cancelBtn = NewButton("Cancel", "", func() { d.result = ResCancel })
	return nil
}

func (d *ValuesDialog) buildRows() {
	decimal := func(r rune) bool { return (r >= '0' && r <= '9') || strings.ContainsRune(".-+eE", r) }
	hexDigit := func(r rune) bool {
		return (r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F')
	}
	printable := func(r rune) bool { return r >= 32 && r != 127 }

	y := float32(8)
	for _, sec := range d.Doc.Sections {
		if sec.Label != "" {
			d.rows = append(d.rows, &valueRow{header: sec.Label, y: y, h: vdHeaderH})
			y += vdHeaderH
		}
		for _, v := range sec.Values {
			v := v // own copy for the closures below
			row := &valueRow{val: v, y: y, h: vdRowH}
			switch v.Kind {
			case KindText:
				row.field1 = &TextField{MaxLen: 60, Allow: printable, OnChange: func(t string) bool {
					if nt := strings.TrimSpace(t); nt != v.Text {
						v.Text, v.Edited = nt, true
					}
					return true
				}}
			case KindFloat, KindInt:
				row.field1 = &TextField{MaxLen: 24, Allow: decimal, OnChange: func(t string) bool {
					ok, changed := v.SetDecimal(t)
					if changed {
						v.Edited = true
					}
					return ok
				}}
				row.field2 = &TextField{MaxLen: 8, Allow: hexDigit, OnChange: func(t string) bool {
					ok, changed := v.SetHex(t)
					if changed {
						v.Edited = true
					}
					return ok
				}}
			}
			for _, f := range []*TextField{row.field1, row.field2} {
				if f != nil {
					d.fields = append(d.fields, fieldRef{f, row})
				}
			}
			d.rows = append(d.rows, row)
			y += vdRowH
		}
	}
	d.contentH = y + 8
}

func (d *ValuesDialog) save() {
	d.result = ResOK
}

// setFocus focuses a field and scrolls its row into view.
func (d *ValuesDialog) setFocus(fr fieldRef) {
	d.focus = fr.field
	if fr.row.y < d.scroll {
		d.scroll = fr.row.y - 8
	}
	if fr.row.y+fr.row.h > d.scroll+d.listH {
		d.scroll = fr.row.y + fr.row.h - d.listH + 8
	}
}

func (d *ValuesDialog) focusIndex() int {
	for i, fr := range d.fields {
		if fr.field == d.focus {
			return i
		}
	}
	return -1
}

// focusMove moves the focus to the next (+1) or previous (-1) field, wrapping around.
func (d *ValuesDialog) focusMove(delta int) {
	n := len(d.fields)
	if n == 0 {
		return
	}
	idx := d.focusIndex()
	switch {
	case idx < 0 && delta > 0:
		idx = 0
	case idx < 0:
		idx = n - 1
	default:
		idx = (idx + delta + n) % n
	}
	d.setFocus(d.fields[idx])
}

// drawField syncs an unfocused field with its model value, then draws it.
func (d *ValuesDialog) drawField(f *TextField, model string, bounds, list rl.Rectangle, can bool) bool {
	f.Bounds = bounds
	if d.focus != f {
		f.Text, f.Invalid = model, false // also reverts half-typed invalid text
	}
	if f.UpdateDraw(can, d.focus == f, list) {
		d.focus = f
		return true
	}
	return false
}

// drawRow draws one row at screen y and reports whether it consumed a click.
func (d *ValuesDialog) drawRow(r *valueRow, list rl.Rectangle, y float32, can bool) (clicked bool) {
	x := list.X
	if r.header != "" {
		rl.DrawText(r.header, int32(x)+12, int32(y)+14, 18, ColAccent)
		rl.DrawRectangle(int32(x)+12, int32(y+r.h)-5, int32(list.Width)-24, 1, rl.Fade(rl.White, 0.12))
		return false
	}

	v := r.val
	if v.Edited { // '*' = edited (not default)
		rl.DrawText("*", int32(x)+14, int32(y)+8, 22, ColAccent)
	}
	name := FitText(v.Name, 104, 16)
	rl.DrawText(name, int32(x)+36, int32(y)+8, 16, ColTextDim)
	fx := x + 150

	switch v.Kind {
	case KindBool:
		m := rl.GetMousePosition()
		area := rl.NewRectangle(fx, y+2, 130, 28)
		hover := can && rl.CheckCollisionPointRec(m, area) && rl.CheckCollisionPointRec(m, list)
		box := rl.NewRectangle(fx, y+4, 24, 24)
		bg := ColBtn
		if hover {
			bg = ColBtnHover
		}
		rl.DrawRectangleRounded(box, 0.25, 4, bg)
		if v.Text == "true" {
			rl.DrawRectangleRounded(rl.NewRectangle(box.X+5, box.Y+5, 14, 14), 0.25, 4, ColAccent)
		}
		rl.DrawText(v.Text, int32(fx)+34, int32(y)+7, 18, ColText)
		if hover && rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
			if v.Text == "true" {
				v.Text = "false"
			} else {
				v.Text = "true"
			}
			v.Edited = true
			clicked = true
		}
	case KindText:
		clicked = d.drawField(r.field1, v.Text, rl.NewRectangle(fx, y+3, 340, 26), list, can)
	case KindEnum:
		display := v.Display
		if display == "" {
			display = v.Text
		}
		rl.DrawText(FitText(display, 340, 16), int32(fx), int32(y)+7, 16, ColText)
	case KindInt, KindFloat:
		c1 := d.drawField(r.field1, v.Text, rl.NewRectangle(fx, y+3, 190, 26), list, can)
		rl.DrawText(":", int32(fx)+198, int32(y)+6, 18, ColTextDim)
		c2 := d.drawField(r.field2, v.Hex, rl.NewRectangle(fx+214, y+3, 118, 26), list, can)
		clicked = c1 || c2
	}
	return clicked
}

// UpdateDraw returns the saved path and ResOK once the file was written (or ResCancel).
func (d *ValuesDialog) UpdateDraw() (string, DialogResult) {
	if !d.Active {
		return "", ResNone
	}
	can := d.openedFrame != Frame // ignore the click that opened it
	d.UpdateDrag(can)
	d.DrawFrame()

	b := d.Bounds
	m := rl.GetMousePosition()
	d.result = ResNone

	list := rl.NewRectangle(b.X+16, b.Y+44, b.Width-32, b.Height-44-70)
	d.listH = list.Height
	maxScroll := d.contentH - list.Height
	if maxScroll < 0 {
		maxScroll = 0
	}
	inList := rl.CheckCollisionPointRec(m, list)

	// --- keyboard / wheel
	if can {
		if inList {
			d.scroll -= rl.GetMouseWheelMove() * 48
		}
		if rl.IsKeyPressed(rl.KeyTab) {
			if rl.IsKeyDown(rl.KeyLeftShift) || rl.IsKeyDown(rl.KeyRightShift) {
				d.focusMove(-1)
			} else {
				d.focusMove(1)
			}
		}
		if rl.IsKeyPressed(rl.KeyEnter) {
			switch idx := d.focusIndex(); {
			case d.focus == nil:
				d.save()
			case idx == len(d.fields)-1:
				d.focus = nil
			default:
				d.focusMove(1)
			}
		}
		if rl.IsKeyPressed(rl.KeyEscape) {
			if d.focus != nil {
				d.focus = nil
			} else {
				d.result = ResCancel
			}
		}
	}
	d.scroll = Clamp(d.scroll, 0, maxScroll)

	// --- rows (clipped to the list)
	rl.DrawRectangleRounded(list, 0.02, 4, ColTitle)
	rl.BeginScissorMode(int32(list.X), int32(list.Y), int32(list.Width), int32(list.Height))
	clickedWidget := false
	for _, r := range d.rows {
		y := list.Y + r.y - d.scroll
		if y+r.h < list.Y || y > list.Y+list.Height {
			continue
		}
		if d.drawRow(r, list, y, can) {
			clickedWidget = true
		}
	}
	rl.EndScissorMode()
	if can && inList && !clickedWidget && rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
		d.focus = nil // clicked empty space
	}

	if maxScroll > 0 { // scrollbar
		thumbH := list.Height * list.Height / d.contentH
		thumbY := list.Y + (list.Height-thumbH)*d.scroll/maxScroll
		rl.DrawRectangleRounded(rl.NewRectangle(list.X+list.Width-7, thumbY, 4, thumbH), 1, 4, rl.Fade(rl.White, 0.3))
	}

	// --- footer
	const bw, bh = 96, 34
	by := b.Y + b.Height - bh - 16
	if d.err != "" {
		rl.DrawText(FitText("Save failed: "+d.err, int32(b.Width)-32-2*bw-30, 14), int32(b.X)+16, int32(by)+10, 14, rl.NewColor(255, 120, 120, 255))
	} else {
		rl.DrawText("* = edited (not default)", int32(b.X)+16, int32(by)+10, 14, ColTextDim)
	}
	rightX := b.X + b.Width - 16
	d.saveBtn.Bounds = rl.NewRectangle(rightX-bw, by, bw, bh)
	d.saveBtn.Enabled = can
	d.saveBtn.UpdateDraw()
	d.cancelBtn.Bounds = rl.NewRectangle(rightX-2*bw-10, by, bw, bh)
	d.cancelBtn.Enabled = can
	d.cancelBtn.UpdateDraw()

	if d.result != ResNone {
		d.Active = false
	}
	return d.Path, d.result
}
