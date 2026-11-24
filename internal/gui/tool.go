package gui

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type LifetimeManager interface {
	IsActive() bool
	Activate(ToolBar)
	Deactivate()
	Cancel()
}

type WorldMouseDownTracker interface {
	WorldMouseDown(world rl.Vector2, button rl.MouseButton)
}
type WorldMouseMoveTracker interface {
	WorldMouseMove(world rl.Vector2)
}

type WorldMouseLeaveTracker interface {
	WorldMouseLeave()
}

type WorldMouseUpTracker interface {
	WorldMouseUp(world rl.Vector2, button rl.MouseButton)
}

type ScreenMouseDownTracker interface {
	ScreenMouseDown(coord rl.Vector2, button rl.MouseButton)
}
type ScreenWorldMouseMoveTracker interface {
	ScreenMouseMove(coord rl.Vector2)
}

type ScreenMouseUpTracker interface {
	ScreenMouseUp(coord rl.Vector2, button rl.MouseButton)
}

type Tool interface {
	LifetimeManager
}

type ToolBar interface {
	ActivateTool(t Tool)
	Workspace() any
}
