package utils

import rl "github.com/gen2brain/raylib-go/raylib"

// RectangleOverlaps checks if two rectangles overlap
func RectangleOverlaps(r1, r2 *rl.Rectangle) bool {
	// Check if one rectangle is completely to the left of the other
	// or completely to the right
	if r1.X >= r2.X+r2.Width || r2.X >= r1.X+r1.Width {
		return false
	}

	// Check if one rectangle is completely above the other
	// or completely below
	if r1.Y >= r2.Y+r2.Height || r2.Y >= r1.Y+r1.Height {
		return false
	}

	// If neither of the above is true, the rectangles must overlap
	return true
}
