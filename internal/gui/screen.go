package gui

import rl "github.com/gen2brain/raylib-go/raylib"

const (
	InitialW      = 900
	InitialH      = 600
	MaxText   int = 32
	MaxPoints int = 2048
)

// Current window size. Refreshed every frame so the whole UI follows
// resizing and maximizing.
var ScreenW float32 = InitialW
var ScreenH float32 = InitialH

// keepOnScreen moves r back inside the window if a resize left it outside.
func KeepOnScreen(r *rl.Rectangle) {
	r.X = Clamp(r.X, 0, ScreenW-r.Width)
	r.Y = Clamp(r.Y, 0, ScreenH-r.Height)
}

func Clamp(v, lo, hi float32) float32 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
