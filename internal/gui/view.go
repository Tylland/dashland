package gui

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

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
			v.Camera.Zoom = Clamp(z, minZoom, maxZoom)
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
	br := v.ScreenToWorld(rl.NewVector2(ScreenW, ScreenH))
	thick := 1 / z // constant 1px on screen
	line := rl.Fade(rl.White, 0.07)
	axis := rl.Fade(ColAccent, 0.55)

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
