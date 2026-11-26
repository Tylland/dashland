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
		t.ws.SetHoveredTarget(nil)
	}
}

func (t *DefaultTool) Deactivate() {
	t.active = false
}

func (t *DefaultTool) Cancel() {
	t.WorldMouseLeave()
}

// Mouse Tracking
func (t *DefaultTool) WorldMouseMove(world rl.Vector2) {
	if t.ws == nil || t.ws.Stage() == nil {
		t.WorldMouseLeave()
		return
	}

	pos, inside := t.ws.Stage().ClampedBlockPosition(world)
	var target *Target
	if inside {
		if _, ok := t.ws.Stage().GetEntityAtPosition(pos); ok {
			target = &Target{Kind: TargetEntity, Position: pos}
		} else if _, ok := t.ws.Stage().GetBlockAtPosition(pos); ok {
			target = &Target{Kind: TargetBlock, Position: pos}
		}
	}

	t.ws.SetHoveredTarget(target)
}

func (t *DefaultTool) WorldMouseLeave() {
	if t.ws != nil {
		t.ws.SetHoveredTarget(nil)
	}
}

func (t *DefaultTool) WorldMouseDown(world rl.Vector2, button rl.MouseButton) {
	if t.ws == nil || t.ws.Stage() == nil {
		return
	}

	if button == rl.MouseButtonRight {
		t.ws.SetSelectedTarget(nil)
		pos, inside := t.ws.Stage().ClampedBlockPosition(world)
		if !inside {
			return
		}
		if block, ok := t.ws.Stage().GetBlockAtPosition(pos); ok {
			t.ws.SetSelectedTarget(&Target{Kind: TargetBlock, Position: pos})
			path := filepath.Join(t.ws.StagePath(), game.EntityFileName(block))
			preview := gui.LoadPreview(path)

			t.ws.Inspect(&gui.EntityInspection{Openpath: path, Preview: preview})

			return
		}
	}

	if button == rl.MouseButtonLeft {
		t.ws.SetSelectedTarget(nil)
		pos, inside := t.ws.Stage().ClampedBlockPosition(world)
		if !inside {
			return
		}
		if entity, ok := t.ws.Stage().GetEntityAtPosition(pos); ok {
			t.ws.SetSelectedTarget(&Target{Kind: TargetEntity, Position: pos})
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
