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
//   - Canvas view: Ctrl+wheel zooms at the cursor; hold Space, press the left
//     button and drag to pan; Home resets the view
//   - Modal dialogs: message (OK), confirm (Yes/No), text input (OK/Cancel)
//   - Dialogs are draggable by their title bar, dim the background,
//     and block toolbar/canvas input while open
//   - Keyboard: Enter = confirm, Esc = cancel
package main

import (
	"fmt"
	"path/filepath"

	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/tylland/dashland/internal/gui"
)

type Dot struct {
	Pos    rl.Vector2
	Color  rl.Color
	Entity int // index into the entities slice
}

func main() {
	rl.InitWindow(gui.ScreenW, gui.ScreenH, "raylib-go - toolbar & dialogs")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)
	rl.SetExitKey(rl.KeyNull) // Esc closes dialogs instead of the app

	// Entities are textures; they must be created after the window exists.
	entities := gui.DefaultEntities()
	defer func() {
		for i := range entities {
			entities[i].Unload()
		}
	}()
	selectedEntity := 0

	palette := []rl.Color{rl.SkyBlue, rl.Orange, rl.Lime, rl.Pink, rl.Gold, rl.Violet}
	colorIdx := 0

	dots := make([]Dot, 0, gui.MaxPoints)
	docName := "Untitled"

	var dlg gui.Dialog
	var fd gui.FileDialog
	var ed gui.EntityDialog
	var lastDir, openedPath string
	var preview []string
	// Build the toolbar: each button carries its own label, tooltip and click handler.
	toolbar := gui.NewToolbar(rl.NewVector2(40, 40))
	toolbar.Add(gui.NewButton("Open", "Browse for a file to open", func() {
		fd.Open(lastDir)
	}))
	entityBtn := toolbar.Add(gui.NewButton("Entity", "Choose which entity gets placed on the canvas", func() {
		ed.Open(entities, selectedEntity, palette[colorIdx])
	}))
	entityBtn.Decorate = func(b *gui.Button) { // thumbnail of the current entity
		thumb := rl.NewRectangle(b.Bounds.X+b.Bounds.Width-30, b.Bounds.Y+5, 20, 20)
		entities[selectedEntity].Draw(thumb, palette[colorIdx])
	}
	toolbar.Add(gui.NewButton("Clear", "Remove all placed entities (asks for confirmation)", func() {
		dlg.Open(gui.DlgConfirm, gui.ActClear, "Clear canvas", "Remove all entities? This cannot be undone.", "")
	}))
	toolbar.Add(gui.NewButton("Rename", "Change the document name", func() {
		dlg.Open(gui.DlgInput, gui.ActRename, "Rename document", "Enter a new name:", docName)
	}))
	colorBtn := toolbar.Add(gui.NewButton("Color", "Cycle the dot color", func() {
		colorIdx = (colorIdx + 1) % len(palette)
	}))
	colorBtn.GapAfter = gui.TbSepExtra
	colorBtn.Decorate = func(b *gui.Button) { // color swatch on the button
		sw := rl.NewRectangle(b.Bounds.X+b.Bounds.Width-22, b.Bounds.Y+7, 16, 16)
		rl.DrawRectangleRounded(sw, 0.3, 4, palette[colorIdx])
	}
	toolbar.Add(gui.NewButton("About", "About this demo", func() {
		dlg.Open(gui.DlgMessage, gui.ActAbout, "About", "raylib-go toolbar + dialogs demo.", "")
	}))
	toolbar.Add(gui.NewButton("Flip", "Switch horizontal / vertical (or double-click the grip)", toolbar.Flip))

	view := gui.NewView()

	for !rl.WindowShouldClose() {
		gui.Frame++
		modal := dlg.Kind != gui.DlgNone || fd.Active || ed.Active

		// canvas input (blocked while a dialog is open or when over the toolbar)
		m := rl.GetMousePosition()
		overToolbar := rl.CheckCollisionPointRec(m, toolbar.Bounds())
		view.Update(!modal, overToolbar)
		if !modal && !view.Grab && !view.Panning && !overToolbar && m.Y < gui.ScreenH-gui.StatusH &&
			rl.IsMouseButtonPressed(rl.MouseButtonLeft) && len(dots) < gui.MaxPoints {
			// dots live in world coordinates, so zoom/pan do not move them
			dots = append(dots, Dot{Pos: view.ScreenToWorld(m), Color: palette[colorIdx], Entity: selectedEntity})
		}

		rl.BeginDrawing()
		rl.ClearBackground(gui.ColBG)

		// world: grid + placed entities, seen through the pan/zoom camera
		view.Begin()
		view.DrawGrid()
		for _, d := range dots {
			entities[d.Entity].Draw(rl.NewRectangle(d.Pos.X-20, d.Pos.Y-20, 40, 40), d.Color)
		}
		view.End()

		// file preview panel
		if openedPath != "" {
			px, py := float32(gui.ScreenW-460), float32(130)
			ph := float32(len(preview)*18 + 44)
			rl.DrawRectangleRounded(rl.NewRectangle(px, py, 440, ph), 0.04, 6, gui.ColPanel)
			rl.DrawText(gui.FitText(openedPath, 400, 16), int32(px)+14, int32(py)+10, 16, gui.ColText)
			for i, l := range preview {
				rl.DrawText(l, int32(px)+14, int32(py)+36+int32(i*18), 14, gui.ColTextDim)
			}
		}

		if len(dots) == 0 && openedPath == "" {
			rl.DrawText("Click anywhere to place the selected entity", gui.ScreenW/2-200, gui.ScreenH/2, 20, gui.ColTextDim)
		}

		// status bar
		rl.DrawRectangle(0, gui.ScreenH-gui.StatusH, gui.ScreenW, gui.StatusH, gui.ColToolbar)
		status := fmt.Sprintf("%s  |  %d placed  |  entity: %s  |  zoom %d%%",
			docName, len(dots), entities[selectedEntity].Name, int(view.Camera.Zoom*100+0.5))
		rl.DrawText(status, 10, gui.ScreenH-20, 16, gui.ColTextDim)
		help := "Ctrl+wheel: zoom  |  Space+drag: pan  |  Home: reset"
		rl.DrawText(help, gui.ScreenW-rl.MeasureText(help, 16)-10, gui.ScreenH-20, 16, gui.ColTextDim)

		// toolbar (its buttons run their own OnClick handlers)
		toolbar.UpdateDraw(!modal)
		// the toolbar resets the cursor each frame, so set the grab cursors after it
		if view.Panning {
			rl.SetMouseCursor(rl.MouseCursorResizeAll)
		} else if view.Grab && !overToolbar {
			rl.SetMouseCursor(rl.MouseCursorPointingHand)
		}

		// dialog (drawn last so it sits on top)
		if res := dlg.UpdateDraw(); res != gui.ResNone {
			if res == gui.ResOK {
				switch dlg.Action {
				case gui.ActClear:
					dots = dots[:0]
				case gui.ActRename:
					if len(dlg.Input) > 0 {
						docName = string(dlg.Input)
					}
				}
			}
			dlg.Kind = gui.DlgNone
		}

		// file selector (modal, drawn on top)
		if path, res := fd.UpdateDraw(); res != gui.ResNone {
			lastDir = fd.Dir
			if res == gui.ResOK {
				openedPath = path
				docName = filepath.Base(path)
				preview = gui.LoadPreview(path)
			}
		}

		// entity selector (modal, drawn on top)
		if idx, res := ed.UpdateDraw(); res == gui.ResOK && idx >= 0 {
			selectedEntity = idx
		}

		rl.EndDrawing()
	}
}
