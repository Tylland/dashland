package tool

import (
	"path/filepath"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/tylland/dashland/internal/game"
	"github.com/tylland/dashland/internal/gui"
)

type DefaultTool struct {
	active bool
	ws     ToolbarWorkspace
}

func NewDefaultTool() *DefaultTool {
	return &DefaultTool{}
}

func (t *DefaultTool) IsActive() bool {
	return t.active
}

func (t *DefaultTool) Activate(tb gui.ToolBar) {
	t.active = true

	if ws, ok := tb.Workspace().(ToolbarWorkspace); ok {
		t.ws = ws
	}
}

func (t *DefaultTool) Deactivate() {
	t.active = false
}

func (t *DefaultTool) Cancel() {
}

// Mouse Tracking
func (t *DefaultTool) WorldMouseDown(world rl.Vector2, button rl.MouseButton) {

	if button == rl.MouseButtonRight {
		pos := t.ws.Stage().GetBlockPosition(world)
		if block, ok := t.ws.Stage().GetBlockAtPosition(pos); ok {
			path := filepath.Join(t.ws.StagePath(), game.BlockFileName(block))
			preview := gui.LoadPreview(path)

			t.ws.Inspect(&gui.EntityInspection{Openpath: path, Preview: preview})

			return
		}
	}

	if button == rl.MouseButtonLeft {
		pos := t.ws.Stage().GetBlockPosition(world)
		if entity, ok := t.ws.Stage().GetEntityAtPosition(pos); ok {
			doc := &gui.Document{}
			writer := gui.NewDocumentWriter(doc)

			entity.WriteTo(writer)

			path := filepath.Join(t.ws.StagePath(), game.EntityFileName(entity))
			preview := gui.LoadPreview(path)

			t.ws.Inspect(&gui.EntityInspection{Entity: entity, Document: doc, Openpath: path, Preview: preview})

			return
		}
	}
}
