// raylib-go FLOATING toolbar + dialogs demo (pure raylib, no raygui)
//
// Object-oriented layout: Button is a reusable widget (label, tooltip, OnClick,
// optional Decorate hook). Toolbar, Dialog and FileDialog own Buttons and only
// position them; main() wires each toolbar button to its action with a closure.
//
// Setup:
//
//	go mod init raylibdemo
//	go get github.com/gen2brain/raylib-go/raylib
//	go run .
//
// Features:
//   - Floating toolbar: drag it by the grip, it snaps to window edges,
//     double-click the grip (or press Flip) to switch horizontal/vertical
//   - Hover/pressed button states and tooltips
//   - Resizable / maximizable window (maximize button or F11); the UI follows the window size
//   - Canvas view: Ctrl+wheel zooms at the cursor; hold Space, press the left
//     button and drag to pan; Home resets the view
//   - Rect tool (toolbar button or R): drag a rectangle to fill its grid cells
//     with the selected entity; Esc or right-click cancels
//   - Modal dialogs: message (OK), confirm (Yes/No), text input (OK/Cancel)
//   - Dialogs are draggable by their title bar, dim the background,
//     and block toolbar/canvas input while open
//   - Keyboard: Enter = confirm, Esc = cancel
package main

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	initialW  = 900 // starting (and minimum) window size
	initialH  = 600
	maxText   = 32
	maxPoints = 2048
)

// Current window size. Refreshed every frame so the whole UI follows
// resizing and maximizing.
var screenW, screenH float32 = initialW, initialH

// ---------------------------------------------------------------- Theme

var (
	colBG       = rl.NewColor(36, 38, 46, 255)
	colToolbar  = rl.NewColor(52, 55, 66, 255)
	colBtn      = rl.NewColor(70, 74, 90, 255)
	colBtnHover = rl.NewColor(92, 98, 122, 255)
	colBtnDown  = rl.NewColor(56, 60, 76, 255)
	colBtnOff   = rl.NewColor(60, 62, 70, 255)
	colAccent   = rl.NewColor(90, 150, 255, 255)
	colPanel    = rl.NewColor(58, 61, 74, 255)
	colTitle    = rl.NewColor(44, 47, 58, 255)
	colText     = rl.NewColor(235, 237, 245, 255)
	colTextDim  = rl.NewColor(140, 145, 160, 255)
)

// -------------------------------------------------------------- Widgets

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
	bg := colBtn
	switch {
	case !b.Enabled:
		bg = colBtnOff
	case b.pressed:
		bg = colBtnDown
	case b.hovered:
		bg = colBtnHover
	}
	rl.DrawRectangleRounded(b.Bounds, 0.25, 6, bg)
	if b.Active {
		rl.DrawRectangleRounded(b.Bounds, 0.25, 6, rl.Fade(colAccent, 0.55))
	}

	fg := colText
	if !b.Enabled {
		fg = colTextDim
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

// keepOnScreen moves r back inside the window if a resize left it outside.
func keepOnScreen(r *rl.Rectangle) {
	r.X = clamp(r.X, 0, screenW-r.Width)
	r.Y = clamp(r.Y, 0, screenH-r.Height)
}

func clamp(v, lo, hi float32) float32 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

// --------------------------------------------------------------- Dialogs

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

var frame uint64

func (d *Dialog) Open(kind DialogKind, action DialogAction, title, message, initial string) {
	*d = Dialog{
		Kind:        kind,
		Action:      action,
		Title:       title,
		Message:     message,
		Input:       []rune(initial),
		openedFrame: frame,
	}
	w, h := float32(420), float32(170)
	if kind == DlgInput {
		h = 200
	}
	d.Bounds = rl.NewRectangle((screenW-w)/2, (screenH-h)/2, w, h)

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

	canInteract := d.openedFrame != frame // ignore the click that opened it
	m := rl.GetMousePosition()

	// --- dragging via title bar
	titleBar := rl.NewRectangle(d.Bounds.X, d.Bounds.Y, d.Bounds.Width, 30)
	if canInteract && rl.IsMouseButtonPressed(rl.MouseButtonLeft) && rl.CheckCollisionPointRec(m, titleBar) {
		d.dragging = true
		d.dragOffset = rl.NewVector2(m.X-d.Bounds.X, m.Y-d.Bounds.Y)
	}
	if d.dragging {
		if rl.IsMouseButtonDown(rl.MouseButtonLeft) {
			d.Bounds.X = clamp(m.X-d.dragOffset.X, 0, screenW-d.Bounds.Width)
			d.Bounds.Y = clamp(m.Y-d.dragOffset.Y, 0, screenH-d.Bounds.Height)
		} else {
			d.dragging = false
		}
	}

	// --- backdrop + panel
	rl.DrawRectangle(0, 0, int32(screenW), int32(screenH), rl.Fade(rl.Black, 0.55))
	shadow := rl.NewRectangle(d.Bounds.X+4, d.Bounds.Y+6, d.Bounds.Width, d.Bounds.Height)
	rl.DrawRectangleRounded(shadow, 0.04, 8, rl.Fade(rl.Black, 0.35))
	rl.DrawRectangleRounded(d.Bounds, 0.04, 8, colPanel)
	rl.DrawRectangleRounded(titleBar, 0.4, 8, colTitle)
	rl.DrawRectangle(int32(d.Bounds.X), int32(d.Bounds.Y)+18, int32(d.Bounds.Width), 12, colTitle)
	rl.DrawText(d.Title, int32(d.Bounds.X)+12, int32(d.Bounds.Y)+6, 18, colText)

	// --- message
	rl.DrawText(d.Message, int32(d.Bounds.X)+20, int32(d.Bounds.Y)+46, 18, colText)

	// --- text input field
	if d.Kind == DlgInput {
		box := rl.NewRectangle(d.Bounds.X+20, d.Bounds.Y+82, d.Bounds.Width-40, 34)
		border := rl.NewRectangle(box.X-1, box.Y-1, box.Width+2, box.Height+2)
		rl.DrawRectangleRounded(border, 0.2, 6, colAccent)
		rl.DrawRectangleRounded(box, 0.2, 6, colTitle)
		rl.DrawText(string(d.Input), int32(box.X)+10, int32(box.Y)+8, 18, colText)

		if canInteract {
			for {
				ch := rl.GetCharPressed()
				if ch <= 0 {
					break
				}
				if ch >= 32 && ch <= 125 && len(d.Input) < maxText {
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
			rl.DrawRectangle(cx, int32(box.Y)+7, 2, 20, colText)
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

// ----------------------------------------------------------- File selector

// matchPattern reports whether a file name matches a wildcard pattern such as "*.entity.txt".
// '*' matches any run of characters and '?' a single one; matching ignores case.
// Several patterns can be given separated by ';' or ',' (e.g. "*.png;*.jpg").
// An empty pattern matches everything.
func matchPattern(pattern, name string) bool {
	pattern = strings.TrimSpace(pattern)
	if pattern == "" {
		return true
	}
	name = strings.ToLower(name)
	for _, p := range strings.FieldsFunc(pattern, func(r rune) bool { return r == ';' || r == ',' }) {
		p = strings.ToLower(strings.TrimSpace(p))
		if p == "" {
			continue
		}
		if ok, err := filepath.Match(p, name); err == nil && ok {
			return true
		}
	}
	return false
}

type FileEntry struct {
	Name  string
	IsDir bool
	Size  int64
}

// FileDialog is an in-app, modal file browser (raylib has no native file dialog).
type FileDialog struct {
	Active        bool
	Bounds        rl.Rectangle
	Dir           string
	Entries       []FileEntry
	Selected      int
	Scroll        int
	Err           string
	dragging      bool
	dragOffset    rl.Vector2
	openedFrame   uint64
	lastClickIdx  int
	lastClickTime float64

	Pattern string // wildcard filter for file names, e.g. "*.entity.txt"; empty = all files

	upBtn     *Button
	openBtn   *Button
	cancelBtn *Button
	pick      string       // set by openBtn
	result    DialogResult // set by cancelBtn
}

const fdRowH = 24

func humanSize(n int64) string {
	const unit = 1024
	if n < unit {
		return fmt.Sprintf("%d B", n)
	}
	div, exp := int64(unit), 0
	for m := n / unit; m >= unit; m /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(n)/float64(div), "KMGT"[exp])
}

// fitText shortens s from the left ("...tail") until it fits maxW pixels.
func fitText(s string, maxW, fontSize int32) string {
	if rl.MeasureText(s, fontSize) <= maxW {
		return s
	}
	r := []rune(s)
	for len(r) > 4 {
		r = r[1:]
		if t := "..." + string(r); rl.MeasureText(t, fontSize) <= maxW {
			return t
		}
	}
	return "..." + string(r)
}

// Open shows the dialog. pattern restricts the listed files with wildcards
// (see matchPattern), e.g. "*.entity.txt"; pass "" to list every file.
func (f *FileDialog) Open(startDir, pattern string) {
	if startDir == "" {
		if wd, err := os.Getwd(); err == nil {
			startDir = wd
		} else {
			startDir = "."
		}
	}
	*f = FileDialog{Active: true, openedFrame: frame, lastClickIdx: -1}
	f.Pattern = pattern
	w, h := float32(560), float32(430)
	f.Bounds = rl.NewRectangle((screenW-w)/2, (screenH-h)/2, w, h)
	f.upBtn = NewButton("Up", "", f.goUp)
	f.openBtn = NewButton("Open", "", func() { f.pick = f.activate() })
	f.cancelBtn = NewButton("Cancel", "", func() { f.result = ResCancel })
	f.loadDir(startDir)
	for rl.GetCharPressed() > 0 { // flush pending characters
	}
}

func (f *FileDialog) loadDir(dir string) {
	f.Dir, f.Entries, f.Selected, f.Scroll, f.Err = dir, nil, -1, 0, ""

	des, err := os.ReadDir(dir)
	if err != nil {
		f.Err = err.Error()
	}
	if filepath.Dir(dir) != dir { // not at the filesystem root
		f.Entries = append(f.Entries, FileEntry{Name: "..", IsDir: true})
	}

	var dirs, files []FileEntry
	for _, de := range des {
		if strings.HasPrefix(de.Name(), ".") { // skip hidden files
			continue
		}
		e := FileEntry{Name: de.Name(), IsDir: de.IsDir()}
		if info, err := de.Info(); err == nil {
			e.Size = info.Size()
		}
		if de.Type()&os.ModeSymlink != 0 { // follow symlinks
			if st, err := os.Stat(filepath.Join(dir, de.Name())); err == nil {
				e.IsDir, e.Size = st.IsDir(), st.Size()
			}
		}
		if e.IsDir {
			dirs = append(dirs, e)
		} else if matchPattern(f.Pattern, e.Name) {
			files = append(files, e)
		}
	}
	byName := func(l []FileEntry) {
		sort.Slice(l, func(i, j int) bool { return strings.ToLower(l[i].Name) < strings.ToLower(l[j].Name) })
	}
	byName(dirs)
	byName(files)
	f.Entries = append(append(f.Entries, dirs...), files...)
}

// activate opens the selected entry: directories are entered, files are returned as a path.
func (f *FileDialog) activate() string {
	if f.Selected < 0 || f.Selected >= len(f.Entries) {
		return ""
	}
	e := f.Entries[f.Selected]
	full := filepath.Join(f.Dir, e.Name)
	if e.IsDir {
		f.loadDir(full)
		return ""
	}
	return full
}

func (f *FileDialog) goUp() {
	if parent := filepath.Dir(f.Dir); parent != f.Dir {
		f.loadDir(parent)
	}
}

// UpdateDraw handles input and draws the dialog. When it closes it returns the chosen path (if any).
func (f *FileDialog) UpdateDraw() (string, DialogResult) {
	if !f.Active {
		return "", ResNone
	}
	can := f.openedFrame != frame // ignore the click that opened it
	m := rl.GetMousePosition()
	result, pick := ResNone, ""

	// --- dragging via title bar
	titleBar := rl.NewRectangle(f.Bounds.X, f.Bounds.Y, f.Bounds.Width, 30)
	if can && rl.IsMouseButtonPressed(rl.MouseButtonLeft) && rl.CheckCollisionPointRec(m, titleBar) {
		f.dragging = true
		f.dragOffset = rl.NewVector2(m.X-f.Bounds.X, m.Y-f.Bounds.Y)
	}
	if f.dragging {
		if rl.IsMouseButtonDown(rl.MouseButtonLeft) {
			f.Bounds.X = clamp(m.X-f.dragOffset.X, 0, screenW-f.Bounds.Width)
			f.Bounds.Y = clamp(m.Y-f.dragOffset.Y, 0, screenH-f.Bounds.Height)
		} else {
			f.dragging = false
		}
	}
	b := f.Bounds

	// --- backdrop + panel
	rl.DrawRectangle(0, 0, int32(screenW), int32(screenH), rl.Fade(rl.Black, 0.55))
	rl.DrawRectangleRounded(rl.NewRectangle(b.X+4, b.Y+6, b.Width, b.Height), 0.03, 8, rl.Fade(rl.Black, 0.35))
	rl.DrawRectangleRounded(b, 0.03, 8, colPanel)
	rl.DrawRectangleRounded(titleBar, 0.4, 8, colTitle)
	rl.DrawRectangle(int32(b.X), int32(b.Y)+18, int32(b.Width), 12, colTitle)
	rl.DrawText("Open file", int32(b.X)+12, int32(b.Y)+6, 18, colText)

	// --- path bar + [active filter chip] + Up button
	chipText, chipW := "", float32(0)
	if strings.TrimSpace(f.Pattern) != "" {
		chipText = fitText(f.Pattern, 140, 14)
		chipW = float32(rl.MeasureText(chipText, 14)) + 20 + 8
	}
	pathRect := rl.NewRectangle(b.X+16, b.Y+40, b.Width-32-60-chipW, 28)
	chipRect := rl.NewRectangle(pathRect.X+pathRect.Width+8, b.Y+40, chipW-8, 28)
	upRect := rl.NewRectangle(b.X+b.Width-16-52, b.Y+40, 52, 28)
	rl.DrawRectangleRounded(pathRect, 0.2, 6, colTitle)
	rl.DrawText(fitText(f.Dir, int32(pathRect.Width)-16, 16), int32(pathRect.X)+8, int32(pathRect.Y)+6, 16, colTextDim)
	if chipText != "" {
		rl.DrawRectangleRounded(chipRect, 0.3, 6, rl.Fade(colAccent, 0.35))
		rl.DrawText(chipText, int32(chipRect.X)+10, int32(chipRect.Y)+7, 14, colText)
	}
	f.upBtn.Bounds, f.upBtn.Enabled = upRect, can
	f.upBtn.UpdateDraw()

	// --- list geometry
	list := rl.NewRectangle(b.X+16, b.Y+78, b.Width-32, b.Height-78-76)
	visible := int(list.Height) / fdRowH
	maxScroll := len(f.Entries) - visible
	if maxScroll < 0 {
		maxScroll = 0
	}

	// --- list input
	if can {
		if rl.CheckCollisionPointRec(m, list) {
			f.Scroll -= int(rl.GetMouseWheelMove()) * 3
		}
		if rl.IsKeyPressed(rl.KeyDown) && f.Selected < len(f.Entries)-1 {
			f.Selected++
		}
		if rl.IsKeyPressed(rl.KeyUp) && f.Selected > 0 {
			f.Selected--
		}
		if rl.IsKeyPressed(rl.KeyBackspace) {
			f.goUp()
		}
		if rl.IsKeyPressed(rl.KeyEnter) && f.Selected >= 0 {
			pick = f.activate()
		}
		if rl.IsMouseButtonPressed(rl.MouseButtonLeft) && rl.CheckCollisionPointRec(m, list) {
			idx := f.Scroll + int((m.Y-list.Y)/fdRowH)
			if idx >= 0 && idx < len(f.Entries) {
				f.Selected = idx
				now := rl.GetTime()
				if idx == f.lastClickIdx && now-f.lastClickTime < 0.35 { // double-click
					pick = f.activate()
					f.lastClickIdx = -1
				} else {
					f.lastClickIdx, f.lastClickTime = idx, now
				}
			}
		}
		// keep the selection visible
		if f.Selected >= 0 {
			if f.Selected < f.Scroll {
				f.Scroll = f.Selected
			}
			if f.Selected >= f.Scroll+visible {
				f.Scroll = f.Selected - visible + 1
			}
		}
	}
	if f.Scroll > maxScroll {
		f.Scroll = maxScroll
	}
	if f.Scroll < 0 {
		f.Scroll = 0
	}

	// --- list drawing (clipped)
	rl.DrawRectangleRounded(list, 0.02, 4, colTitle)
	rl.BeginScissorMode(int32(list.X), int32(list.Y), int32(list.Width), int32(list.Height))
	for i := f.Scroll; i < len(f.Entries) && i < f.Scroll+visible+1; i++ {
		e := f.Entries[i]
		row := rl.NewRectangle(list.X, list.Y+float32((i-f.Scroll)*fdRowH), list.Width, fdRowH)
		switch {
		case i == f.Selected:
			rl.DrawRectangleRec(row, rl.Fade(colAccent, 0.45))
		case can && rl.CheckCollisionPointRec(m, row) && rl.CheckCollisionPointRec(m, list):
			rl.DrawRectangleRec(row, rl.Fade(rl.White, 0.08))
		}

		// icon: folder = accent tab+body, file = gray page
		ix, iy := int32(row.X)+8, int32(row.Y)+5
		if e.IsDir {
			rl.DrawRectangle(ix, iy, 8, 4, colAccent)
			rl.DrawRectangle(ix, iy+3, 16, 11, colAccent)
		} else {
			rl.DrawRectangle(ix+2, iy, 12, 14, colTextDim)
		}

		rl.DrawText(e.Name, int32(row.X)+32, int32(row.Y)+4, 16, colText)
		if !e.IsDir {
			sz := humanSize(e.Size)
			rl.DrawText(sz, int32(row.X+row.Width)-rl.MeasureText(sz, 14)-16, int32(row.Y)+5, 14, colTextDim)
		}
	}
	if f.Err != "" {
		rl.DrawText(fitText("Cannot read folder: "+f.Err, int32(list.Width)-24, 14), int32(list.X)+12, int32(list.Y)+fdRowH+8, 14, rl.NewColor(255, 120, 120, 255))
	} else if len(f.Entries) == 0 || (len(f.Entries) == 1 && f.Entries[0].Name == "..") {
		msg := "(empty folder)"
		if strings.TrimSpace(f.Pattern) != "" {
			msg = "(no matching files)"
		}
		rl.DrawText(msg, int32(list.X)+32, int32(list.Y)+fdRowH+6, 14, colTextDim)
	}
	rl.EndScissorMode()

	// scrollbar
	if maxScroll > 0 {
		thumbH := list.Height * float32(visible) / float32(len(f.Entries))
		thumbY := list.Y + (list.Height-thumbH)*float32(f.Scroll)/float32(maxScroll)
		rl.DrawRectangleRounded(rl.NewRectangle(list.X+list.Width-7, thumbY, 4, thumbH), 1, 4, rl.Fade(rl.White, 0.3))
	}

	// --- footer: selected name + buttons
	const bw, bh = 96, 34
	by := b.Y + b.Height - bh - 16
	label := "No file selected"
	if f.Selected >= 0 && f.Selected < len(f.Entries) {
		label = f.Entries[f.Selected].Name
	}
	rl.DrawText(fitText(label, int32(b.Width)-32-2*bw-30, 16), int32(b.X)+16, int32(by)+9, 16, colTextDim)

	rightX := b.X + b.Width - 16
	f.pick, f.result = "", ResNone
	f.openBtn.Bounds = rl.NewRectangle(rightX-bw, by, bw, bh)
	f.openBtn.Enabled = can && f.Selected >= 0
	f.openBtn.UpdateDraw()
	f.cancelBtn.Bounds = rl.NewRectangle(rightX-2*bw-10, by, bw, bh)
	f.cancelBtn.Enabled = can
	f.cancelBtn.UpdateDraw()
	if can && rl.IsKeyPressed(rl.KeyEscape) {
		f.result = ResCancel
	}
	if f.pick != "" {
		pick = f.pick
	}
	if f.result != ResNone {
		result = f.result
	}

	if pick != "" {
		result = ResOK
	}
	if result != ResNone {
		f.Active = false
	}
	return pick, result
}

// loadPreview reads the first lines of a file as text (or reports that it is binary).
func loadPreview(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		return []string{"Could not open file:", err.Error()}
	}
	defer file.Close()

	buf := make([]byte, 4096)
	n, _ := file.Read(buf)
	buf = buf[:n]
	if bytes.IndexByte(buf, 0) >= 0 {
		return []string{"(binary file - no preview)"}
	}
	if n == 0 {
		return []string{"(empty file)"}
	}
	lines := strings.Split(strings.ReplaceAll(string(buf), "\r", ""), "\n")
	if len(lines) > 18 {
		lines = lines[:18]
	}
	for i, l := range lines {
		l = strings.ReplaceAll(l, "\t", "    ")
		if len(l) > 70 {
			l = l[:70] + "..."
		}
		lines[i] = l
	}
	return lines
}

// ------------------------------------------------------------------ Rectangle grid tool

const cellSize = 40 // world units per grid cell; entities are drawn cellSize x cellSize

// RectTool fills a dragged rectangle of grid cells with the selected entity.
// Cells are aligned to the world grid, so repeated fills line up exactly.
type RectTool struct {
	Active     bool // tool is selected
	dragging   bool
	startX     int // cell coordinates where the drag began
	startY     int
	endX, endY int // cell under the mouse
}

// cellOf returns the grid cell containing a world position.
func cellOf(world rl.Vector2) (int, int) {
	return int(math.Floor(float64(world.X / cellSize))), int(math.Floor(float64(world.Y / cellSize)))
}

func (t *RectTool) Dragging() bool { return t.dragging }

func (t *RectTool) Begin(world rl.Vector2) {
	t.dragging = true
	t.startX, t.startY = cellOf(world)
	t.endX, t.endY = t.startX, t.startY
}

func (t *RectTool) Drag(world rl.Vector2) { t.endX, t.endY = cellOf(world) }
func (t *RectTool) Cancel()               { t.dragging = false }

// Cells returns the dragged range as inclusive, normalized cell coordinates.
func (t *RectTool) Cells() (x0, y0, x1, y1 int) {
	x0, x1 = t.startX, t.endX
	y0, y1 = t.startY, t.endY
	if x0 > x1 {
		x0, x1 = x1, x0
	}
	if y0 > y1 {
		y0, y1 = y1, y0
	}
	return
}

// DrawPreview draws the selection and ghost entities (call in world space).
func (t *RectTool) DrawPreview(e *Entity, tint rl.Color, zoom float32) {
	if !t.dragging {
		return
	}
	x0, y0, x1, y1 := t.Cells()
	area := rl.NewRectangle(float32(x0)*cellSize, float32(y0)*cellSize,
		float32(x1-x0+1)*cellSize, float32(y1-y0+1)*cellSize)
	rl.DrawRectangleRec(area, rl.Fade(colAccent, 0.12))

	ghost := rl.Fade(tint, 0.55)
	drawn := 0
	for cy := y0; cy <= y1 && drawn < maxPoints; cy++ {
		for cx := x0; cx <= x1 && drawn < maxPoints; cx++ {
			e.Draw(rl.NewRectangle(float32(cx)*cellSize, float32(cy)*cellSize, cellSize, cellSize), ghost)
			drawn++
		}
	}
	rl.DrawRectangleLinesEx(area, 2/zoom, colAccent)
}

// DrawLabel shows "W x H = N" next to the mouse (call in screen space).
func (t *RectTool) DrawLabel(mouse rl.Vector2, existing int) {
	if !t.dragging {
		return
	}
	x0, y0, x1, y1 := t.Cells()
	w, h := x1-x0+1, y1-y0+1
	txt := fmt.Sprintf("%d x %d = %d", w, h, w*h)
	col := colText
	if existing+w*h > maxPoints {
		txt += "  (limit reached)"
		col = rl.NewColor(255, 120, 120, 255)
	}
	r := rl.NewRectangle(mouse.X+16, mouse.Y+16, float32(rl.MeasureText(txt, 16))+16, 26)
	r.X = clamp(r.X, 4, screenW-r.Width-4)
	r.Y = clamp(r.Y, 4, screenH-r.Height-4)
	rl.DrawRectangleRounded(r, 0.3, 6, rl.NewColor(20, 20, 26, 230))
	rl.DrawText(txt, int32(r.X)+8, int32(r.Y)+5, 16, col)
}

// FillCellRect places one entity in every cell of the inclusive range. Cells that
// already hold an entity at that exact spot are replaced. Stops at maxPoints.
func FillCellRect(dots []Dot, x0, y0, x1, y1, entity int, col rl.Color) []Dot {
	index := make(map[rl.Vector2]int, len(dots))
	for i, d := range dots {
		index[d.Pos] = i
	}
	for cy := y0; cy <= y1; cy++ {
		for cx := x0; cx <= x1; cx++ {
			pos := rl.NewVector2((float32(cx)+0.5)*cellSize, (float32(cy)+0.5)*cellSize)
			d := Dot{Pos: pos, Color: col, Entity: entity}
			if i, ok := index[pos]; ok {
				dots[i] = d
				continue
			}
			if len(dots) >= maxPoints {
				return dots
			}
			dots = append(dots, d)
		}
	}
	return dots
}

// ------------------------------------------------------------------ View (zoom + pan)

const (
	minZoom  = 0.2
	maxZoom  = 8.0
	zoomStep = 1.1 // zoom factor per wheel notch
)

// View is a pannable, zoomable camera over the canvas world.
//   - Ctrl + mouse wheel          : zoom, anchored on the cursor
//   - hold Space, press left button: "grab" the canvas; keep the left button down
//     and move the mouse to pan
//   - Home                        : reset zoom and pan
type View struct {
	Camera  rl.Camera2D
	Grab    bool // Space is held: the canvas can be grabbed (open-hand cursor)
	Panning bool // left button is down after a grab: moving the mouse pans
}

func NewView() *View {
	v := &View{}
	v.Reset()
	return v
}

func (v *View) Reset() { v.Camera = rl.Camera2D{Zoom: 1} }

func (v *View) ScreenToWorld(p rl.Vector2) rl.Vector2 { return rl.GetScreenToWorld2D(p, v.Camera) }

func (v *View) Begin() { rl.BeginMode2D(v.Camera) }
func (v *View) End()   { rl.EndMode2D() }

// Update handles input. enabled is false while a modal dialog is open;
// overUI is true when the mouse is over a floating UI element (e.g. the toolbar).
func (v *View) Update(enabled, overUI bool) {
	if !enabled {
		v.Grab, v.Panning = false, false
		return
	}

	// --- pan: Space + left button down, then move the mouse
	v.Grab = rl.IsKeyDown(rl.KeySpace)
	if v.Grab && !overUI && rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
		v.Panning = true
	}
	if v.Panning {
		if rl.IsMouseButtonDown(rl.MouseButtonLeft) { // pans for as long as the button stays down
			d := rl.GetMouseDelta()
			v.Camera.Target.X -= d.X / v.Camera.Zoom
			v.Camera.Target.Y -= d.Y / v.Camera.Zoom
		} else {
			v.Panning = false
		}
	}

	if rl.IsKeyPressed(rl.KeyHome) {
		v.Reset()
	}

	// --- zoom: Ctrl + wheel, keeping the world point under the cursor fixed
	ctrl := rl.IsKeyDown(rl.KeyLeftControl) || rl.IsKeyDown(rl.KeyRightControl)
	if ctrl && !overUI {
		if wheel := rl.GetMouseWheelMove(); wheel != 0 {
			m := rl.GetMousePosition()
			world := v.ScreenToWorld(m)
			v.Camera.Offset = m
			v.Camera.Target = world

			z := v.Camera.Zoom * float32(math.Pow(zoomStep, float64(wheel)))
			v.Camera.Zoom = clamp(z, minZoom, maxZoom)
		}
	}
}

// DrawGrid draws an adaptive background grid (call between Begin and End).
func (v *View) DrawGrid() {
	z := v.Camera.Zoom
	spacing := float32(100)
	for spacing*z < 28 {
		spacing *= 2
	}
	for spacing*z > 220 {
		spacing /= 2
	}

	tl := v.ScreenToWorld(rl.NewVector2(0, 0))
	br := v.ScreenToWorld(rl.NewVector2(screenW, screenH))
	thick := 1 / z // constant 1px on screen
	line := rl.Fade(rl.White, 0.07)
	axis := rl.Fade(colAccent, 0.55)

	for x := float32(math.Floor(float64(tl.X/spacing))) * spacing; x <= br.X; x += spacing {
		c := line
		if x == 0 {
			c = axis
		}
		rl.DrawLineEx(rl.NewVector2(x, tl.Y), rl.NewVector2(x, br.Y), thick, c)
	}
	for y := float32(math.Floor(float64(tl.Y/spacing))) * spacing; y <= br.Y; y += spacing {
		c := line
		if y == 0 {
			c = axis
		}
		rl.DrawLineEx(rl.NewVector2(tl.X, y), rl.NewVector2(br.X, y), thick, c)
	}
}

// ------------------------------------------------------------------ Entities

const entityTexSize = 64

// Entity is one selectable item. Its look is a rl.Texture2D, which is either
// rendered procedurally (NewEntity) or loaded from an image file (LoadEntity).
type Entity struct {
	Name    string
	Texture rl.Texture2D
	Flipped bool // render-target textures are stored upside down by OpenGL

	target rl.RenderTexture2D // owned render target (zero value when loaded from a file)
}

// NewEntity renders draw() into a fresh texture. The shapes are drawn in white/grey
// so that a tint color can be applied when the entity is displayed.
// Must be called after InitWindow and outside BeginDrawing/EndDrawing.
func NewEntity(name string, draw func(c rl.Vector2)) Entity {
	rt := rl.LoadRenderTexture(entityTexSize, entityTexSize)
	rl.BeginTextureMode(rt)
	rl.ClearBackground(rl.Blank)
	draw(rl.NewVector2(entityTexSize/2, entityTexSize/2))
	rl.EndTextureMode()
	// rl.SetTextureFilter(rt.Texture, rl.FilterBilinear) // optional: smoother scaling
	return Entity{Name: name, Texture: rt.Texture, Flipped: true, target: rt}
}

// LoadEntity creates an entity from an image file (png, jpg, ...).
func LoadEntity(name, path string) Entity {
	return Entity{Name: name, Texture: rl.LoadTexture(path)}
}

// Draw renders the entity's texture into dest, multiplied by tint.
func (e *Entity) Draw(dest rl.Rectangle, tint rl.Color) {
	h := float32(e.Texture.Height)
	if e.Flipped {
		h = -h
	}
	src := rl.NewRectangle(0, 0, float32(e.Texture.Width), h)
	rl.DrawTexturePro(e.Texture, src, dest, rl.NewVector2(0, 0), 0, tint)
}

func (e *Entity) Unload() {
	if e.target.ID != 0 {
		rl.UnloadRenderTexture(e.target)
	} else {
		rl.UnloadTexture(e.Texture)
	}
}

// tri draws a triangle regardless of vertex winding order.
func tri(a, b, c rl.Vector2, col rl.Color) {
	rl.DrawTriangle(a, b, c, col)
	rl.DrawTriangle(a, c, b, col)
}

func DefaultEntities() []Entity {
	white := rl.White
	grey := rl.NewColor(160, 160, 175, 255)
	v := rl.NewVector2

	return []Entity{
		NewEntity("Circle", func(c rl.Vector2) {
			rl.DrawCircleV(c, 26, white)
		}),
		NewEntity("Square", func(c rl.Vector2) {
			rl.DrawRectangleRounded(rl.NewRectangle(c.X-24, c.Y-24, 48, 48), 0.15, 6, white)
		}),
		NewEntity("Triangle", func(c rl.Vector2) {
			tri(v(c.X, c.Y-28), v(c.X+28, c.Y+22), v(c.X-28, c.Y+22), white)
		}),
		NewEntity("Diamond", func(c rl.Vector2) {
			rl.DrawPoly(c, 4, 30, 0, white)
		}),
		NewEntity("Pentagon", func(c rl.Vector2) {
			rl.DrawPoly(c, 5, 28, -90, white)
		}),
		NewEntity("Hexagon", func(c rl.Vector2) {
			rl.DrawPoly(c, 6, 28, 0, white)
		}),
		NewEntity("Ring", func(c rl.Vector2) {
			rl.DrawRing(c, 14, 27, 0, 360, 48, white)
		}),
		NewEntity("Cross", func(c rl.Vector2) {
			rl.DrawRectangle(int32(c.X)-9, int32(c.Y)-26, 18, 52, white)
			rl.DrawRectangle(int32(c.X)-26, int32(c.Y)-9, 52, 18, white)
		}),
		NewEntity("Star", func(c rl.Vector2) {
			var pts [10]rl.Vector2
			for i := range pts {
				r := 28.0
				if i%2 == 1 {
					r = 12
				}
				a := (-90 + float64(i)*36) * math.Pi / 180
				pts[i] = v(c.X+float32(math.Cos(a)*r), c.Y+float32(math.Sin(a)*r))
			}
			for i := range pts {
				tri(c, pts[i], pts[(i+1)%10], white)
			}
		}),
		NewEntity("Target", func(c rl.Vector2) {
			rl.DrawCircleV(c, 27, white)
			rl.DrawCircleV(c, 18, grey)
			rl.DrawCircleV(c, 9, white)
		}),
		NewEntity("Heart", func(c rl.Vector2) {
			rl.DrawCircleV(v(c.X-13, c.Y-8), 14, white)
			rl.DrawCircleV(v(c.X+13, c.Y-8), 14, white)
			tri(v(c.X-25, c.Y-2), v(c.X+25, c.Y-2), v(c.X, c.Y+26), white)
		}),
		NewEntity("Gem", func(c rl.Vector2) {
			rl.DrawPoly(c, 6, 28, 0, white)
			rl.DrawPolyLines(c, 6, 18, 0, grey)
			rl.DrawPolyLines(c, 6, 9, 0, grey)
		}),
		NewEntity("Arrow", func(c rl.Vector2) {
			tri(v(c.X, c.Y-28), v(c.X+24, c.Y+2), v(c.X-24, c.Y+2), white)
			rl.DrawRectangle(int32(c.X)-9, int32(c.Y), 18, 28, white)
		}),
		NewEntity("Pill", func(c rl.Vector2) {
			rl.DrawRectangleRounded(rl.NewRectangle(c.X-12, c.Y-28, 24, 56), 1, 12, white)
		}),
	}
}

// ------------------------------------------------------------------ Item grid

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
	rl.DrawRectangleRounded(g.Bounds, 0.02, 4, colTitle)
	rl.BeginScissorMode(int32(g.Bounds.X), int32(g.Bounds.Y), int32(g.Bounds.Width), int32(g.Bounds.Height))
	for row := 0; row < visRows; row++ {
		for col := 0; col < cols; col++ {
			i := (g.Scroll+row)*cols + col
			if i >= len(g.Items) {
				break
			}
			r := g.cellRect(row, col, ox)

			if i == g.Selected {
				rl.DrawRectangleRounded(rl.NewRectangle(r.X-2, r.Y-2, r.Width+4, r.Height+4), 0.12, 6, colAccent)
			}
			bg := colBtn
			if i == hovered {
				bg = colBtnHover
			}
			rl.DrawRectangleRounded(r, 0.12, 6, bg)
			if i == g.Selected {
				rl.DrawRectangleRounded(r, 0.12, 6, rl.Fade(colAccent, 0.35))
			}

			icon := rl.NewRectangle(r.X+(r.Width-gridIcon)/2, r.Y+8, gridIcon, gridIcon)
			g.Items[i].Draw(icon, g.Tint)

			name := g.Items[i].Name
			tw := rl.MeasureText(name, 14)
			rl.DrawText(name, int32(r.X)+(int32(r.Width)-tw)/2, int32(r.Y)+gridIcon+16, 14, colText)
		}
	}
	rl.EndScissorMode()

	if maxScroll > 0 { // scrollbar
		thumbH := g.Bounds.Height * float32(visRows) / float32(totalRows)
		thumbY := g.Bounds.Y + (g.Bounds.Height-thumbH)*float32(g.Scroll)/float32(maxScroll)
		rl.DrawRectangleRounded(rl.NewRectangle(g.Bounds.X+g.Bounds.Width-7, thumbY, 4, thumbH), 1, 4, rl.Fade(rl.White, 0.3))
	}
}

// ------------------------------------------------------------------ Window + EntityDialog

// Window is a reusable draggable modal frame: dimmed backdrop, shadow, panel and title bar.
type Window struct {
	Bounds     rl.Rectangle
	Title      string
	dragging   bool
	dragOffset rl.Vector2
}

func NewWindow(title string, w, h float32) Window {
	return Window{Title: title, Bounds: rl.NewRectangle((screenW-w)/2, (screenH-h)/2, w, h)}
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
			w.Bounds.X = clamp(m.X-w.dragOffset.X, 0, screenW-w.Bounds.Width)
			w.Bounds.Y = clamp(m.Y-w.dragOffset.Y, 0, screenH-w.Bounds.Height)
		} else {
			w.dragging = false
		}
	}
}

func (w *Window) DrawFrame() {
	b := w.Bounds
	rl.DrawRectangle(0, 0, int32(screenW), int32(screenH), rl.Fade(rl.Black, 0.55))
	rl.DrawRectangleRounded(rl.NewRectangle(b.X+4, b.Y+6, b.Width, b.Height), 0.03, 8, rl.Fade(rl.Black, 0.35))
	rl.DrawRectangleRounded(b, 0.03, 8, colPanel)
	rl.DrawRectangleRounded(w.titleBar(), 0.4, 8, colTitle)
	rl.DrawRectangle(int32(b.X), int32(b.Y)+18, int32(b.Width), 12, colTitle)
	rl.DrawText(w.Title, int32(b.X)+12, int32(b.Y)+6, 18, colText)
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
	*d = EntityDialog{Window: NewWindow("Select entity", 500, 470), Active: true, openedFrame: frame}
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
	can := d.openedFrame != frame // ignore the click that opened it
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
	rl.DrawText(fitText(label, int32(b.Width)-32-2*bw-30, 16), int32(b.X)+16, int32(by)+9, 16, colTextDim)

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

// --------------------------------------------------------------- Floating toolbar

const (
	gripSize   = 18 // thickness of the drag handle
	tbPad      = 8
	tbGap      = 6
	tbBtnW     = 96
	tbBtnH     = 30
	tbSepExtra = 10 // handy value for Button.GapAfter (group separator)
	snapDist   = 16 // magnetic distance to screen edges
	snapMargin = 8  // resting distance from the edge once snapped
	statusH    = 24 // height of the status bar at the bottom
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
}

func NewToolbar(pos rl.Vector2) *Toolbar { return &Toolbar{Pos: pos} }

// Add appends a button and returns it so callers can tweak it further.
func (t *Toolbar) Add(b *Button) *Button {
	t.Buttons = append(t.Buttons, b)
	return b
}

// Flip switches between horizontal and vertical layout.
func (t *Toolbar) Flip() {
	t.Vertical = !t.Vertical
	t.clampToScreen()
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
		span := n*tbBtnH + gaps*tbGap + extra
		return tbBtnW + 2*tbPad, gripSize + tbPad + span + tbPad
	}
	span := n*tbBtnW + gaps*tbGap + extra
	return gripSize + tbPad + span + tbPad, tbBtnH + 2*tbPad
}

func (t *Toolbar) Bounds() rl.Rectangle {
	w, h := t.size()
	return rl.NewRectangle(t.Pos.X, t.Pos.Y, w, h)
}

func (t *Toolbar) gripRect() rl.Rectangle {
	w, h := t.size()
	if t.Vertical {
		return rl.NewRectangle(t.Pos.X, t.Pos.Y, w, gripSize)
	}
	return rl.NewRectangle(t.Pos.X, t.Pos.Y, gripSize, h)
}

// clampToScreen keeps the whole toolbar inside the window (above the status bar).
func (t *Toolbar) clampToScreen() {
	w, h := t.size()
	t.Pos.X = clamp(t.Pos.X, 0, screenW-w)
	t.Pos.Y = clamp(t.Pos.Y, 0, screenH-statusH-h)
}

// snap pulls the toolbar to a screen edge when it is close to one.
func (t *Toolbar) snap() {
	w, h := t.size()
	if t.Pos.X < snapDist {
		t.Pos.X = snapMargin
	}
	if t.Pos.X+w > screenW-snapDist {
		t.Pos.X = screenW - w - snapMargin
	}
	if t.Pos.Y < snapDist {
		t.Pos.Y = snapMargin
	}
	if t.Pos.Y+h > screenH-statusH-snapDist {
		t.Pos.Y = screenH - statusH - h - snapMargin
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
			t.clampToScreen()
			t.snap()
		} else {
			t.dragging = false
		}
	}
	if t.dragging || overGrip {
		rl.SetMouseCursor(rl.MouseCursorResizeAll)
	} else {
		rl.SetMouseCursor(rl.MouseCursorDefault)
	}

	// --- panel
	b := t.Bounds()
	rl.DrawRectangleRounded(rl.NewRectangle(b.X+3, b.Y+5, b.Width, b.Height), 0.25, 8, rl.Fade(rl.Black, 0.40))
	rl.DrawRectangleRounded(rl.NewRectangle(b.X-1, b.Y-1, b.Width+2, b.Height+2), 0.25, 8, rl.Fade(rl.White, 0.12))
	rl.DrawRectangleRounded(b, 0.25, 8, colToolbar)

	// grip dots
	dotCol := colTextDim
	if t.dragging || overGrip {
		dotCol = colAccent
	}
	for i := -1; i <= 1; i++ {
		if t.Vertical {
			rl.DrawCircle(int32(b.X+b.Width/2)+int32(i*8), int32(b.Y+gripSize/2), 2, dotCol)
		} else {
			rl.DrawCircle(int32(b.X+gripSize/2), int32(b.Y+b.Height/2)+int32(i*8), 2, dotCol)
		}
	}

	// --- buttons: the toolbar only positions them, each Button handles itself
	vertical := t.Vertical // a button callback may flip the layout mid-loop
	x, y := b.X+tbPad, b.Y+tbPad
	if vertical {
		y += gripSize
	} else {
		x += gripSize
	}

	var hovered *Button
	for _, btn := range t.Buttons {
		btn.Bounds = rl.NewRectangle(x, y, tbBtnW, tbBtnH)
		btn.Enabled = enabled && !t.dragging
		btn.UpdateDraw()
		if btn.Hovered() {
			hovered = btn
		}
		if vertical {
			y += tbBtnH + tbGap + btn.GapAfter
		} else {
			x += tbBtnW + tbGap + btn.GapAfter
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
		tr.X = clamp(tr.X, 4, screenW-tr.Width-4)
		tr.Y = clamp(tr.Y, 4, screenH-tr.Height-4)
		rl.DrawRectangleRounded(tr, 0.3, 6, rl.NewColor(20, 20, 26, 240))
		rl.DrawText(hovered.Tooltip, int32(tr.X)+8, int32(tr.Y)+5, 16, colText)
	}
}

// ------------------------------------------------------------------ App

type Dot struct {
	Pos    rl.Vector2
	Color  rl.Color
	Entity int // index into the entities slice
}

func main() {
	rl.SetConfigFlags(rl.FlagWindowResizable) // enables the maximize button and edge resizing
	rl.InitWindow(initialW, initialH, "raylib-go - toolbar & dialogs")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)
	rl.SetWindowMinSize(initialW, initialH) // the layout is designed for at least this size
	rl.SetExitKey(rl.KeyNull)               // Esc closes dialogs instead of the app

	// Entities are textures; they must be created after the window exists.
	entities := DefaultEntities()
	defer func() {
		for i := range entities {
			entities[i].Unload()
		}
	}()
	selectedEntity := 0

	palette := []rl.Color{rl.SkyBlue, rl.Orange, rl.Lime, rl.Pink, rl.Gold, rl.Violet}
	colorIdx := 0

	dots := make([]Dot, 0, maxPoints)
	docName := "Untitled"

	var dlg Dialog
	var fd FileDialog
	var ed EntityDialog
	var rectTool RectTool
	var rectBtn *Button
	var lastDir, openedPath string
	// Wildcard filter for the Open dialog. Use "" for all files, or several
	// patterns separated by ';', e.g. "*.entity.txt;*.png".
	const filePattern = "*.entity.txt"
	var preview []string
	// Build the toolbar: each button carries its own label, tooltip and click handler.
	toolbar := NewToolbar(rl.NewVector2(40, 40))
	toolbar.Add(NewButton("Open", "Browse for a file to open", func() {
		fd.Open(lastDir, filePattern)
	}))
	entityBtn := toolbar.Add(NewButton("Entity", "Choose which entity gets placed on the canvas", func() {
		ed.Open(entities, selectedEntity, palette[colorIdx])
	}))
	entityBtn.Decorate = func(b *Button) { // thumbnail of the current entity
		thumb := rl.NewRectangle(b.Bounds.X+b.Bounds.Width-30, b.Bounds.Y+5, 20, 20)
		entities[selectedEntity].Draw(thumb, palette[colorIdx])
	}
	toggleRect := func() {
		rectTool.Active = !rectTool.Active
		rectBtn.Active = rectTool.Active
		rectTool.Cancel()
	}
	rectBtn = toolbar.Add(NewButton("Rect", "Grid tool: drag a rectangle to fill it with the entity (R)", toggleRect))
	toolbar.Add(NewButton("Clear", "Remove all placed entities (asks for confirmation)", func() {
		dlg.Open(DlgConfirm, ActClear, "Clear canvas", "Remove all entities? This cannot be undone.", "")
	}))
	toolbar.Add(NewButton("Rename", "Change the document name", func() {
		dlg.Open(DlgInput, ActRename, "Rename document", "Enter a new name:", docName)
	}))
	colorBtn := toolbar.Add(NewButton("Color", "Cycle the dot color", func() {
		colorIdx = (colorIdx + 1) % len(palette)
	}))
	colorBtn.GapAfter = tbSepExtra
	colorBtn.Decorate = func(b *Button) { // color swatch on the button
		sw := rl.NewRectangle(b.Bounds.X+b.Bounds.Width-22, b.Bounds.Y+7, 16, 16)
		rl.DrawRectangleRounded(sw, 0.3, 4, palette[colorIdx])
	}
	toolbar.Add(NewButton("About", "About this demo", func() {
		dlg.Open(DlgMessage, ActAbout, "About", "raylib-go toolbar + dialogs demo.", "")
	}))
	toolbar.Add(NewButton("Flip", "Switch horizontal / vertical (or double-click the grip)", toolbar.Flip))

	view := NewView()

	for !rl.WindowShouldClose() {
		frame++

		// follow window resizing / maximizing
		screenW, screenH = float32(rl.GetScreenWidth()), float32(rl.GetScreenHeight())
		if rl.IsWindowResized() {
			toolbar.clampToScreen()
			keepOnScreen(&dlg.Bounds)
			keepOnScreen(&fd.Bounds)
			keepOnScreen(&ed.Bounds)
		}
		if rl.IsKeyPressed(rl.KeyF11) { // toggle maximized
			if rl.IsWindowMaximized() {
				rl.RestoreWindow()
			} else {
				rl.MaximizeWindow()
			}
		}
		modal := dlg.Kind != DlgNone || fd.Active || ed.Active

		// canvas input (blocked while a dialog is open or when over the toolbar)
		m := rl.GetMousePosition()
		overToolbar := rl.CheckCollisionPointRec(m, toolbar.Bounds())
		view.Update(!modal, overToolbar)
		if !modal && rl.IsKeyPressed(rl.KeyR) {
			toggleRect()
		}

		canDraw := !modal && !view.Grab && !view.Panning && !overToolbar && m.Y < screenH-statusH
		world := view.ScreenToWorld(m) // dots live in world coordinates, so zoom/pan do not move them
		switch {
		case rectTool.Active:
			if canDraw && rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
				rectTool.Begin(world)
			}
			if rectTool.Dragging() {
				rectTool.Drag(world)
				switch {
				case modal || rl.IsKeyPressed(rl.KeyEscape) || rl.IsMouseButtonPressed(rl.MouseButtonRight):
					rectTool.Cancel()
				case rl.IsMouseButtonReleased(rl.MouseButtonLeft):
					x0, y0, x1, y1 := rectTool.Cells()
					dots = FillCellRect(dots, x0, y0, x1, y1, selectedEntity, palette[colorIdx])
					rectTool.Cancel()
				}
			}
		case canDraw && rl.IsMouseButtonPressed(rl.MouseButtonLeft) && len(dots) < maxPoints:
			dots = append(dots, Dot{Pos: world, Color: palette[colorIdx], Entity: selectedEntity})
		}

		rl.BeginDrawing()
		rl.ClearBackground(colBG)

		// world: grid + placed entities, seen through the pan/zoom camera
		view.Begin()
		view.DrawGrid()
		for _, d := range dots {
			entities[d.Entity].Draw(rl.NewRectangle(d.Pos.X-20, d.Pos.Y-20, 40, 40), d.Color)
		}
		rectTool.DrawPreview(&entities[selectedEntity], palette[colorIdx], view.Camera.Zoom)
		view.End()
		rectTool.DrawLabel(m, len(dots))

		// file preview panel
		if openedPath != "" {
			px, py := float32(screenW-460), float32(130)
			ph := float32(len(preview)*18 + 44)
			rl.DrawRectangleRounded(rl.NewRectangle(px, py, 440, ph), 0.04, 6, colPanel)
			rl.DrawText(fitText(openedPath, 400, 16), int32(px)+14, int32(py)+10, 16, colText)
			for i, l := range preview {
				rl.DrawText(l, int32(px)+14, int32(py)+36+int32(i*18), 14, colTextDim)
			}
		}

		if len(dots) == 0 && openedPath == "" {
			rl.DrawText("Click anywhere to place the selected entity", int32(screenW)/2-200, int32(screenH)/2, 20, colTextDim)
		}

		// status bar
		rl.DrawRectangle(0, int32(screenH)-statusH, int32(screenW), statusH, colToolbar)
		status := fmt.Sprintf("%s  |  %d placed  |  entity: %s  |  zoom %d%%",
			docName, len(dots), entities[selectedEntity].Name, int(view.Camera.Zoom*100+0.5))
		rl.DrawText(status, 10, int32(screenH)-20, 16, colTextDim)
		help := "Ctrl+wheel: zoom  |  Space+drag: pan  |  Home: reset"
		rl.DrawText(help, int32(screenW)-rl.MeasureText(help, 16)-10, int32(screenH)-20, 16, colTextDim)

		// toolbar (its buttons run their own OnClick handlers)
		toolbar.UpdateDraw(!modal)
		// the toolbar resets the cursor each frame, so set the grab cursors after it
		if view.Panning {
			rl.SetMouseCursor(rl.MouseCursorResizeAll)
		} else if view.Grab && !overToolbar {
			rl.SetMouseCursor(rl.MouseCursorPointingHand)
		} else if rectTool.Active && !modal && !overToolbar {
			rl.SetMouseCursor(rl.MouseCursorCrosshair)
		}

		// dialog (drawn last so it sits on top)
		if res := dlg.UpdateDraw(); res != ResNone {
			if res == ResOK {
				switch dlg.Action {
				case ActClear:
					dots = dots[:0]
				case ActRename:
					if len(dlg.Input) > 0 {
						docName = string(dlg.Input)
					}
				}
			}
			dlg.Kind = DlgNone
		}

		// file selector (modal, drawn on top)
		if path, res := fd.UpdateDraw(); res != ResNone {
			lastDir = fd.Dir
			if res == ResOK {
				openedPath = path
				docName = filepath.Base(path)
				preview = loadPreview(path)
			}
		}

		// entity selector (modal, drawn on top)
		if idx, res := ed.UpdateDraw(); res == ResOK && idx >= 0 {
			selectedEntity = idx
		}

		rl.EndDrawing()
	}
}
