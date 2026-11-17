package gui

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

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

// FitText shortens s from the left ("...tail") until it fits maxW pixels.
func FitText(s string, maxW, fontSize int32) string {
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

func (f *FileDialog) Open(startDir string) {
	if startDir == "" {
		if wd, err := os.Getwd(); err == nil {
			startDir = wd
		} else {
			startDir = "."
		}
	}
	*f = FileDialog{Active: true, openedFrame: Frame, lastClickIdx: -1}
	w, h := float32(560), float32(430)
	f.Bounds = rl.NewRectangle((ScreenW-w)/2, (ScreenH-h)/2, w, h)
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
		} else {
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
	can := f.openedFrame != Frame // ignore the click that opened it
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
			f.Bounds.X = clamp(m.X-f.dragOffset.X, 0, ScreenW-f.Bounds.Width)
			f.Bounds.Y = clamp(m.Y-f.dragOffset.Y, 0, ScreenH-f.Bounds.Height)
		} else {
			f.dragging = false
		}
	}
	b := f.Bounds

	// --- backdrop + panel
	rl.DrawRectangle(0, 0, ScreenW, ScreenH, rl.Fade(rl.Black, 0.55))
	rl.DrawRectangleRounded(rl.NewRectangle(b.X+4, b.Y+6, b.Width, b.Height), 0.03, 8, rl.Fade(rl.Black, 0.35))
	rl.DrawRectangleRounded(b, 0.03, 8, ColPanel)
	rl.DrawRectangleRounded(titleBar, 0.4, 8, ColTitle)
	rl.DrawRectangle(int32(b.X), int32(b.Y)+18, int32(b.Width), 12, ColTitle)
	rl.DrawText("Open file", int32(b.X)+12, int32(b.Y)+6, 18, ColText)

	// --- path bar + Up button
	pathRect := rl.NewRectangle(b.X+16, b.Y+40, b.Width-32-60, 28)
	upRect := rl.NewRectangle(pathRect.X+pathRect.Width+8, b.Y+40, 52, 28)
	rl.DrawRectangleRounded(pathRect, 0.2, 6, ColTitle)
	rl.DrawText(FitText(f.Dir, int32(pathRect.Width)-16, 16), int32(pathRect.X)+8, int32(pathRect.Y)+6, 16, ColTextDim)
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
	rl.DrawRectangleRounded(list, 0.02, 4, ColTitle)
	rl.BeginScissorMode(int32(list.X), int32(list.Y), int32(list.Width), int32(list.Height))
	for i := f.Scroll; i < len(f.Entries) && i < f.Scroll+visible+1; i++ {
		e := f.Entries[i]
		row := rl.NewRectangle(list.X, list.Y+float32((i-f.Scroll)*fdRowH), list.Width, fdRowH)
		switch {
		case i == f.Selected:
			rl.DrawRectangleRec(row, rl.Fade(ColAccent, 0.45))
		case can && rl.CheckCollisionPointRec(m, row) && rl.CheckCollisionPointRec(m, list):
			rl.DrawRectangleRec(row, rl.Fade(rl.White, 0.08))
		}

		// icon: folder = accent tab+body, file = gray page
		ix, iy := int32(row.X)+8, int32(row.Y)+5
		if e.IsDir {
			rl.DrawRectangle(ix, iy, 8, 4, ColAccent)
			rl.DrawRectangle(ix, iy+3, 16, 11, ColAccent)
		} else {
			rl.DrawRectangle(ix+2, iy, 12, 14, ColTextDim)
		}

		rl.DrawText(e.Name, int32(row.X)+32, int32(row.Y)+4, 16, ColText)
		if !e.IsDir {
			sz := humanSize(e.Size)
			rl.DrawText(sz, int32(row.X+row.Width)-rl.MeasureText(sz, 14)-16, int32(row.Y)+5, 14, ColTextDim)
		}
	}
	if f.Err != "" {
		rl.DrawText(FitText("Cannot read folder: "+f.Err, int32(list.Width)-24, 14), int32(list.X)+12, int32(list.Y)+fdRowH+8, 14, rl.NewColor(255, 120, 120, 255))
	} else if len(f.Entries) == 0 || (len(f.Entries) == 1 && f.Entries[0].Name == "..") {
		rl.DrawText("(empty folder)", int32(list.X)+32, int32(list.Y)+fdRowH+6, 14, ColTextDim)
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
	rl.DrawText(FitText(label, int32(b.Width)-32-2*bw-30, 16), int32(b.X)+16, int32(by)+9, 16, ColTextDim)

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
func LoadPreview(path string) []string {
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
