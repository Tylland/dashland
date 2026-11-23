package tool

import (
	"bytes"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/tylland/dashland/internal/common"
	"github.com/tylland/dashland/internal/game"
	"github.com/tylland/dashland/internal/gui"
)

type InspectTool struct {
	active bool
	last   common.BlockPosition
	ws     ToolbarWorkspace
}

func NewInspectTool() *InspectTool {
	return &InspectTool{}
}

func (t *InspectTool) IsActive() bool {
	return t.active
}

func (t *InspectTool) Activate(tb gui.ToolBar) {
	t.active = true

	if ws, ok := tb.Workspace().(ToolbarWorkspace); ok {
		t.ws = ws
	}
}

func (t *InspectTool) Deactivate() {
	t.active = false
}

func (t *InspectTool) Cancel() {
}

// Mouse Tracking
func (t *InspectTool) WorldMouseMove(world rl.Vector2) {

	pos := t.ws.Stage().GetBlockPosition(world)

	if !pos.IsSame(t.last) {
		if block, ok := t.ws.Stage().GetBlockAtPosition(pos); ok {
			var buffer bytes.Buffer
			writer := game.NewTextWriter(&buffer)

			block.WriteTo(writer)

			t.ws.Inspect(&gui.EntityInspection{Openpath: "inmemory", Preview: strings.Split(buffer.String(), "\n")})

			return
		}

		if entity, ok := t.ws.Stage().GetEntityAtPosition(pos); ok {
			var buffer bytes.Buffer
			writer := game.NewTextWriter(&buffer)

			entity.WriteTo(writer)

			t.ws.Inspect(&gui.EntityInspection{Openpath: "inmemory", Preview: strings.Split(buffer.String(), "\n")})

			return
		}
	}
}
