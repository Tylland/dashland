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

	"github.com/tylland/dashland/cmd/stagemaker/renderer"
	"github.com/tylland/dashland/cmd/stagemaker/tool"
	"github.com/tylland/dashland/internal/assets"
	"github.com/tylland/dashland/internal/ecs"
	"github.com/tylland/dashland/internal/game"
	"github.com/tylland/dashland/internal/gui"
)

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

	var lastDir string = "C:\\src\\dashland\\assets\\stages\\start"

	// Build the toolbar: each button carries its own label, tooltip and click handler.
	var workspace tool.ToolbarWorkspace = ws
	toolbar := gui.NewToolbar(rl.NewVector2(40, 40), workspace)
	toolbar.Add(gui.NewButton("Open...", "Browse for a file to open", func() {
		fd.Open(lastDir, "stage.txt")
		fd.OnSelected = func(path string) {
			if fd.Mode == gui.FileDialogModeFile {
				ws.OpenStage(path)
			}
		}
	}))

	toolbar.Add(gui.NewButton("New...", "Create new stage", func() {
		fd.OpenDirectory(lastDir)
		fd.OnSelected = func(path string) {
			if fd.Mode == gui.FileDialogModeDirectory {
				stage := game.NewStage("<new stage>", game.NewMapSize(30, 20, 32, 32), assets.LoadTexture("blocks"), assets.LoadTexture("entities"), assets.LoadTexture("groundCorners"))
				doc := gui.NewDocument(stage)

				vd.Open(stage, doc, filepath.Join(path, "stage.txt"))
				vd.OnSaved = func(stageFilePath string) {
					ws.OpenStage(stageFilePath)
				}
			}
		}
	}))

	toolbar.Add(gui.NewButton("Save", "Save stage", func() {
		preview := game.Must(ws.stage.Save(ws.stagePath))

		ws.Inspect(&gui.EntityInspection{Openpath: "saved", Preview: preview})
	}))

	propBtn := toolbar.Add(gui.NewButton("Props...", "View and edit properties of current stage", func() {
		doc := &gui.Document{}
		writer := gui.NewDocumentWriter(doc)

		ws.stage.WriteTo(writer)

		vd.Open(ws.stage, doc, filepath.Join(ws.stagePath))
	}))

	propBtn.GapAfter = gui.TbSepExtra

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

	toolbar.Add(gui.NewButton("Create", "Grid tool: drag a rectangle to fill it with the entity (R)", toggleRect))

	toolbar.Add(gui.NewButton("Inspect", "Click on entity to inspect properties", func() {
		toolbar.ActivateTool(tool.NewInspectTool())
	}))

	toolbar.Add(gui.NewButton("Edit", "Click on entity to edit properties", func() {
		toolbar.ActivateTool(tool.NewDefaultTool())
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
	ws.world.AddSystem(renderer.NewEditorRenderer(ws, &ViewCamera{view: ws.view}))

	for !rl.WindowShouldClose() {
		gui.Frame++

		// follow window resizing / maximizing
		gui.ScreenW, gui.ScreenH = float32(rl.GetScreenWidth()), float32(rl.GetScreenHeight())
		if rl.IsWindowResized() {
			toolbar.ClampToScreen()
			gui.KeepOnScreen(&dlg.Bounds)
			gui.KeepOnScreen(&fd.Bounds)
			gui.KeepOnScreen(&ed.Bounds)
			gui.KeepOnScreen(&vd.Bounds)
		}
		if rl.IsKeyPressed(rl.KeyF11) { // toggle maximized
			if rl.IsWindowMaximized() {
				rl.RestoreWindow()
			} else {
				rl.MaximizeWindow()
			}
		}
		modal := dlg.Kind != gui.DlgNone || fd.Active || ed.Active || vd.Active

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
		}

		//ws.view.DrawGrid()

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
			ph := float32(len(ws.inspection.Preview)*18 + 44)
			rl.DrawRectangleRounded(rl.NewRectangle(px, py, 440, ph), 0.04, 6, rl.Fade(gui.ColPanel, 0.5))
			rl.DrawText(gui.FitText(ws.inspection.Openpath, 400, 16), int32(px)+14, int32(py)+10, 16, gui.ColText)
			for i, l := range ws.inspection.Preview {
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

		// Set the canvas cursor first; the toolbar overrides it while hovered.
		if ws.view.Panning || (ws.view.Grab && !overToolbar) {
			rl.SetMouseCursor(rl.MouseCursorResizeAll)
		} else {
			rl.SetMouseCursor(rl.MouseCursorDefault)
		}

		// toolbar (its buttons run their own OnClick handlers)
		toolbar.UpdateDraw(!modal)

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
		if _, res := fd.UpdateDraw(); res != gui.ResNone {
			lastDir = fd.Dir
			// if res == gui.ResOK && fd.Mode == gui.FileDialogModeFile {
			// 	docName = filepath.Base(path)
			// 	info := gui.LoadPreview(path)

			// 	ws.Inspect(&gui.EntityInspection{Openpath: path, Preview: info})
			// 	ws.world.Clear()

			// 	stagePath := strings.TrimSuffix(path, docName)
			// 	st, err := openStage(ws.world, stagePath)

			// 	if err != nil {
			// 		panic(err)
			// 	}

			// 	ws.stage = st
			// 	ws.stagePath = stagePath
			// 	ws.world.AddSystem(systems.NewEditorRenderer(ws.stage, NewViewCamera(ws.view)))
			// }
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

func createEntityItem(entityType ecs.EntityType) *gui.EntityItem {
	return gui.NewEntityItem(entityType, assets.LoadTexture("entities"), entitySource(entityType))
}
func createEnemyItem(entityType ecs.EntityType) *gui.EntityItem {
	return gui.NewEntityItem(entityType, assets.LoadTexture("enemies"), entitySource(entityType))
}

func CreateEntityItems() []*gui.EntityItem {
	return []*gui.EntityItem{
		createEntityItem(game.EntityDiamond),
		createEntityItem(game.EntityBoulder),
		createEntityItem(game.EntityMagicWall),
		createEntityItem(game.EntityExitDoor),
		createEnemyItem(game.EntityFirefly),
		createEnemyItem(game.EntityButterfly),
	}
}

type StageMakerWorkspace struct {
	world             *ecs.World
	view              *gui.View
	stage             *game.Stage
	stagePath         string
	entities          []*gui.EntityItem
	selectedEntity    int
	hoveredTarget     *tool.Target
	selectedTarget    *tool.Target
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
func (ws *StageMakerWorkspace) HoveredTarget() *tool.Target { return ws.hoveredTarget }
func (ws *StageMakerWorkspace) SetHoveredTarget(target *tool.Target) {
	ws.hoveredTarget = target
}
func (ws *StageMakerWorkspace) SelectedTarget() *tool.Target { return ws.selectedTarget }
func (ws *StageMakerWorkspace) SetSelectedTarget(target *tool.Target) {
	ws.selectedTarget = target
}

func (ws *StageMakerWorkspace) OpenStage(path string) {
	docName := filepath.Base(path)
	stagePath := strings.TrimSuffix(path, docName)

	ws.world.Clear()

	st, err := openStage(ws.world, stagePath)

	if err != nil {
		panic(err)
	}

	ws.stage = st
	ws.stagePath = stagePath
	ws.hoveredTarget = nil
	ws.selectedTarget = nil
	ws.world.AddSystem(renderer.NewEditorRenderer(ws, NewViewCamera(ws.view)))
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
