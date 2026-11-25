package tool

import (
	"bytes"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/tylland/dashland/internal/game"
	"github.com/tylland/dashland/internal/gui"
)

type InspectTool struct {
	active     bool
	lastTarget *Target
	ws         ToolbarWorkspace
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
		t.ws.SetHoveredTarget(nil)
	}
}

func (t *InspectTool) Deactivate() {
	t.active = false
}

func (t *InspectTool) Cancel() {
	t.WorldMouseLeave()
}

// Mouse Tracking
func (t *InspectTool) WorldMouseMove(world rl.Vector2) {
	if t.ws == nil || t.ws.Stage() == nil {
		t.WorldMouseLeave()
		return
	}

	pos, inside := t.ws.Stage().ClampedBlockPosition(world)

	var target *Target
	if inside {
		if _, ok := t.ws.Stage().GetBlockAtPosition(pos); ok {
			target = &Target{Kind: TargetBlock, Position: pos}
		} else if _, ok := t.ws.Stage().GetEntityAtPosition(pos); ok {
			target = &Target{Kind: TargetEntity, Position: pos}
		}
	}

	t.ws.SetHoveredTarget(target)
	if SameTarget(target, t.lastTarget) {
		return
	}

	t.lastTarget = target
	if target == nil {
		return
	}

	// TODO: Remove block handling like WriteBlockTo, all are entities
	if target.Kind == TargetBlock {
		if block, ok := t.ws.Stage().GetBlockAtPosition(pos); ok {
			var buffer bytes.Buffer
			writer := game.NewTextWriter(&buffer)

			if err := game.WriteBlockTo(block, writer); err != nil {
				panic(err)
			}

			t.ws.Inspect(&gui.EntityInspection{Openpath: "inmemory", Preview: strings.Split(buffer.String(), "\n")})
		}
	} else if target.Kind == TargetEntity {
		if entity, ok := t.ws.Stage().GetEntityAtPosition(pos); ok {
			var buffer bytes.Buffer
			writer := game.NewTextWriter(&buffer)

			entity.WriteTo(writer)

			t.ws.Inspect(&gui.EntityInspection{Openpath: "inmemory", Preview: strings.Split(buffer.String(), "\n")})
		}
	}
}

func (t *InspectTool) WorldMouseLeave() {
	t.lastTarget = nil
	if t.ws != nil {
		t.ws.SetHoveredTarget(nil)
	}
}
