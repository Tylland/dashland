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
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/tylland/dashland/cmd/stagemaker/tool"
	"github.com/tylland/dashland/internal/assets"
	"github.com/tylland/dashland/internal/ecs"
	"github.com/tylland/dashland/internal/game"
	"github.com/tylland/dashland/internal/gui"
	"github.com/tylland/dashland/internal/systems"
)

// const (
// 	initialW  = 900 // starting (and minimum) window size
// 	initialH  = 600
// 	maxText   = 32
// 	maxPoints = 2048
// )

// // Current window size. Refreshed every frame so the whole UI follows
// // resizing and maximizing.
// var screenW, screenH float32 = initialW, initialH

func main() {
	rl.SetConfigFlags(rl.FlagWindowResizable) // enables the maximize button and edge resizing
	rl.InitWindow(int32(gui.ScreenW), int32(gui.ScreenH), "Dashland - stage maker")
	defer rl.CloseWindow()
	rl.SetWindowMinSize(int(gui.ScreenW), int(gui.ScreenH)) // the layout is designed for at least this size

	ws := &StageMakerWorkspace{
		world:    ecs.NewWorld(),
		entities: CreateEntityItems(),
	}

	palette := []rl.Color{rl.SkyBlue, rl.Orange, rl.Lime, rl.Pink, rl.Gold, rl.Violet}
	colorIdx := 0

	docName := "Untitled"

	var dlg gui.Dialog
	var fd gui.FileDialog
	var ed gui.EntityDialog
	var vd gui.ValuesDialog

	var lastDir string

	// Build the toolbar: each button carries its own label, tooltip and click handler.
	var workspace tool.ToolbarWorkspace = ws
	toolbar := gui.NewToolbar(rl.NewVector2(40, 40), workspace)
	toolbar.Add(gui.NewButton("Open", "Browse for a file to open", func() {
		fd.Open(lastDir, "stage.txt")
	}))

	toolbar.Add(gui.NewButton("Save", "Save stage", func() {
		ws.stage.Save(ws.stagePath)
	}))

	entityBtn := toolbar.Add(gui.NewButton("Entity", "Choose which entity gets placed on the canvas", func() {
		ed.Open(ws.entities, ws.selectedEntity, palette[colorIdx])
	}))
	entityBtn.Decorate = func(b *gui.Button) { // thumbnail of the current entity
		thumb := rl.NewRectangle(b.Bounds.X+b.Bounds.Width-30, b.Bounds.Y+5, 20, 20)
		ws.entities[ws.selectedEntity].Draw(thumb, palette[colorIdx])
	}

	toggleRect := func() {
		toolbar.ActivateTool(tool.NewFillRectTool(ws.stage))
	}

	toolbar.Add(gui.NewButton("Rect", "Grid tool: drag a rectangle to fill it with the entity (R)", toggleRect))

	toolbar.Add(gui.NewButton("Inspect", "Grid tool: drag a rectangle to fill it with the entity (R)", func() {
		toolbar.ActivateTool(tool.NewDefaultTool())
	}))

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

	ws.view = gui.NewView()
	ws.world.AddSystem(systems.NewRenderSystem(ws.stage, &ViewCamera{view: ws.view}))

	for !rl.WindowShouldClose() {
		gui.Frame++

		// follow window resizing / maximizing
		gui.ScreenW, gui.ScreenH = float32(rl.GetScreenWidth()), float32(rl.GetScreenHeight())
		if rl.IsWindowResized() {
			toolbar.ClampToScreen()
			gui.KeepOnScreen(&dlg.Bounds)
			gui.KeepOnScreen(&fd.Bounds)
			gui.KeepOnScreen(&ed.Bounds)
		}
		if rl.IsKeyPressed(rl.KeyF11) { // toggle maximized
			if rl.IsWindowMaximized() {
				rl.RestoreWindow()
			} else {
				rl.MaximizeWindow()
			}
		}
		modal := dlg.Kind != gui.DlgNone || fd.Active || ed.Active

		// canvas input (blocked while a dialog is open or when over the toolbar)
		m := rl.GetMousePosition()
		overToolbar := rl.CheckCollisionPointRec(m, toolbar.Bounds())
		ws.view.Update(!modal, overToolbar)

		rl.BeginDrawing()
		rl.ClearBackground(gui.ColBG)

		// world: grid + placed entities, seen through the pan/zoom camera
		ws.view.Begin()

		if ws.stage != nil {
			ws.world.Update(0)
			// ws.stage.Update(0)
			// ws.stage.Render(0)
		}

		ws.view.DrawGrid()

		canDraw := !modal && !ws.view.Grab && !ws.view.Panning && !overToolbar && m.Y < gui.ScreenH-gui.StatusH
		world := ws.view.ScreenToWorld(m) // dots live in world coordinates, so zoom/pan do not move them

		if toolbar.ActiveTool != nil && toolbar.ActiveTool.IsActive() {
			if tracker, ok := toolbar.ActiveTool.(gui.WorldMouseDownTracker); ok {
				if canDraw && rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
					tracker.WorldMouseDown(world, rl.MouseButtonLeft)
				}
				if canDraw && rl.IsMouseButtonPressed(rl.MouseButtonRight) {
					tracker.WorldMouseDown(world, rl.MouseButtonRight)
				}
			}

			if tracker, ok := toolbar.ActiveTool.(gui.WorldMouseMoveTracker); ok {
				tracker.WorldMouseMove(world)
			}

			if tracker, ok := toolbar.ActiveTool.(gui.WorldMouseUpTracker); ok {
				if rl.IsMouseButtonReleased(rl.MouseButtonLeft) {
					tracker.WorldMouseUp(world, rl.MouseButtonLeft)
				}
			}

			if modal || rl.IsKeyPressed(rl.KeyEscape) || rl.IsMouseButtonPressed(rl.MouseButtonRight) {
				toolbar.ActiveTool.Cancel()
			}
		}

		ws.view.End()

		if tracker, ok := toolbar.ActiveTool.(gui.ScreenWorldMouseMoveTracker); ok {
			tracker.ScreenMouseMove(m)
		}

		// file preview panel
		if ws.inspection != nil && ws.inspection.Openpath != "" {
			px, py := float32(gui.ScreenW-460), float32(130)
			ph := float32(len(ws.inspection.Info)*18 + 44)
			rl.DrawRectangleRounded(rl.NewRectangle(px, py, 440, ph), 0.04, 6, gui.ColPanel)
			rl.DrawText(gui.FitText(ws.inspection.Openpath, 400, 16), int32(px)+14, int32(py)+10, 16, gui.ColText)
			for i, l := range ws.inspection.Info {
				rl.DrawText(l, int32(px)+14, int32(py)+36+int32(i*18), 14, gui.ColTextDim)
			}
		}

		if ws.inspection != nil && !ws.inspectionHandled && ws.inspection.Entity != nil && ws.inspection.Document != nil {
			vd.Open(ws.inspection.Entity, ws.inspection.Document, ws.inspection.Openpath)
			ws.inspectionHandled = true
		}

		// status bar
		rl.DrawRectangle(0, int32(gui.ScreenH-gui.StatusH), int32(gui.ScreenW), gui.StatusH, gui.ColToolbar)
		status := fmt.Sprintf("%s  |  %d placed  |  entity: %s  |  zoom %d%%",
			docName, len(ws.entities), ws.entities[ws.selectedEntity].Name, int(ws.view.Camera.Zoom*100+0.5))
		rl.DrawText(status, 10, int32(gui.ScreenH-20), 16, gui.ColTextDim)
		help := "Ctrl+wheel: zoom  |  Space+drag: pan  |  Home: reset"
		rl.DrawText(help, int32(gui.ScreenW)-rl.MeasureText(help, 16)-10, int32(gui.ScreenH)-20, 16, gui.ColTextDim)

		// toolbar (its buttons run their own OnClick handlers)
		toolbar.UpdateDraw(!modal)
		// the toolbar resets the cursor each frame, so set the grab cursors after it
		if ws.view.Panning {
			rl.SetMouseCursor(rl.MouseCursorResizeAll)
		} else if ws.view.Grab && !overToolbar {
			rl.SetMouseCursor(rl.MouseCursorPointingHand)
		}

		vd.UpdateDraw()

		// dialog (drawn last so it sits on top)
		if res := dlg.UpdateDraw(); res != gui.ResNone {
			if res == gui.ResOK {
				switch dlg.Action {
				case gui.ActClear:
					//					dots = dots[:0]
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
				docName = filepath.Base(path)
				info := gui.LoadPreview(path)

				ws.Inspect(&gui.EntityInspection{Openpath: path, Info: info})
				ws.world.Clear()

				stagePath := strings.TrimSuffix(path, docName)
				st, err := openStage(ws.world, stagePath)

				if err != nil {
					panic(err)
				}

				ws.stage = st
				ws.stagePath = stagePath
				ws.world.AddSystem(systems.NewRenderSystem(ws.stage, NewViewCamera(ws.view)))
			}
		}

		// entity selector (modal, drawn on top)
		if idx, res := ed.UpdateDraw(); res == gui.ResOK && idx >= 0 {
			ws.selectedEntity = idx
		}

		rl.EndDrawing()
	}
}

func openStage(world *ecs.World, path string) (*game.Stage, error) {

	blockTexture := assets.LoadTexture("blocks")
	entityTextures := assets.LoadTexture("entities")
	groundCorners := assets.LoadTexture("ground_corners")

	mapSize := game.MapSize{}
	stage := game.NewStage("generated", mapSize, blockTexture, entityTextures, groundCorners)

	err := stage.Load(world, path)

	if err != nil {
		return nil, err
	}

	return stage, nil
}

func entitySource(entityType ecs.EntityType) rl.Rectangle {
	return rl.NewRectangle(float32(entityType)*32, 0, 32, 32)
}

func CreateEntityItems() []*gui.EntityItem {
	return []*gui.EntityItem{
		gui.NewEntityItem(game.EntityBoulder, assets.LoadTexture("entities"), entitySource(game.EntityBoulder)),
	}
}

type StageMakerWorkspace struct {
	world             *ecs.World
	view              *gui.View
	stage             *game.Stage
	stagePath         string
	entities          []*gui.EntityItem
	selectedEntity    int
	inspection        *gui.EntityInspection
	inspectionHandled bool
}

func (ws *StageMakerWorkspace) World() *ecs.World                 { return ws.world }
func (ws *StageMakerWorkspace) View() *gui.View                   { return ws.view }
func (ws *StageMakerWorkspace) Stage() *game.Stage                { return ws.stage }
func (ws *StageMakerWorkspace) StagePath() string                 { return ws.stagePath }
func (ws *StageMakerWorkspace) Inspection() *gui.EntityInspection { return ws.inspection }
func (ws *StageMakerWorkspace) Inspect(inspection *gui.EntityInspection) {
	ws.inspection = inspection
	ws.inspectionHandled = false
}
func (ws *StageMakerWorkspace) Entities() []*gui.EntityItem { return ws.entities }
func (ws *StageMakerWorkspace) SelectEntity(idx int) {
	ws.selectedEntity = idx
}
func (ws *StageMakerWorkspace) SelectedEntity() *gui.EntityItem {
	return ws.entities[ws.selectedEntity]
}

type ViewCamera struct {
	view *gui.View
}

func NewViewCamera(view *gui.View) *ViewCamera {
	return &ViewCamera{view: view}
}

func (c *ViewCamera) GetCamera() rl.Camera2D {
	return c.view.Camera
}

func (c *ViewCamera) Update(deltaTime float32) {

}
