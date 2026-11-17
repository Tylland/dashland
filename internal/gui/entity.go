package gui

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const entityTexSize = 64

// Entity is one selectable item. Its look is a rl.Texture2D, which is either
// rendered procedurally (NewEntity) or loaded from an image file (LoadEntity).
type Entity struct {
	Name    string
	Texture rl.Texture2D
	Flipped bool // render-target textures are stored upside down by OpenGL

	target rl.RenderTexture2D // owned render target (zero value when loaded from a file)
}

// NewEntity renders draw() into a fresh texture. The shapes are drawn in white/grey
// so that a tint color can be applied when the entity is displayed.
// Must be called after InitWindow and outside BeginDrawing/EndDrawing.
func NewEntity(name string, draw func(c rl.Vector2)) Entity {
	rt := rl.LoadRenderTexture(entityTexSize, entityTexSize)
	rl.BeginTextureMode(rt)
	rl.ClearBackground(rl.Blank)
	draw(rl.NewVector2(entityTexSize/2, entityTexSize/2))
	rl.EndTextureMode()
	// rl.SetTextureFilter(rt.Texture, rl.FilterBilinear) // optional: smoother scaling
	return Entity{Name: name, Texture: rt.Texture, Flipped: true, target: rt}
}

// LoadEntity creates an entity from an image file (png, jpg, ...).
func LoadEntity(name, path string) Entity {
	return Entity{Name: name, Texture: rl.LoadTexture(path)}
}

// Draw renders the entity's texture into dest, multiplied by tint.
func (e *Entity) Draw(dest rl.Rectangle, tint rl.Color) {
	h := float32(e.Texture.Height)
	if e.Flipped {
		h = -h
	}
	src := rl.NewRectangle(0, 0, float32(e.Texture.Width), h)
	rl.DrawTexturePro(e.Texture, src, dest, rl.NewVector2(0, 0), 0, tint)
}

func (e *Entity) Unload() {
	if e.target.ID != 0 {
		rl.UnloadRenderTexture(e.target)
	} else {
		rl.UnloadTexture(e.Texture)
	}
}

// tri draws a triangle regardless of vertex winding order.
func tri(a, b, c rl.Vector2, col rl.Color) {
	rl.DrawTriangle(a, b, c, col)
	rl.DrawTriangle(a, c, b, col)
}

func DefaultEntities() []Entity {
	white := rl.White
	grey := rl.NewColor(160, 160, 175, 255)
	v := rl.NewVector2

	return []Entity{
		NewEntity("Circle", func(c rl.Vector2) {
			rl.DrawCircleV(c, 26, white)
		}),
		NewEntity("Square", func(c rl.Vector2) {
			rl.DrawRectangleRounded(rl.NewRectangle(c.X-24, c.Y-24, 48, 48), 0.15, 6, white)
		}),
		NewEntity("Triangle", func(c rl.Vector2) {
			tri(v(c.X, c.Y-28), v(c.X+28, c.Y+22), v(c.X-28, c.Y+22), white)
		}),
		NewEntity("Diamond", func(c rl.Vector2) {
			rl.DrawPoly(c, 4, 30, 0, white)
		}),
		NewEntity("Pentagon", func(c rl.Vector2) {
			rl.DrawPoly(c, 5, 28, -90, white)
		}),
		NewEntity("Hexagon", func(c rl.Vector2) {
			rl.DrawPoly(c, 6, 28, 0, white)
		}),
		NewEntity("Ring", func(c rl.Vector2) {
			rl.DrawRing(c, 14, 27, 0, 360, 48, white)
		}),
		NewEntity("Cross", func(c rl.Vector2) {
			rl.DrawRectangle(int32(c.X)-9, int32(c.Y)-26, 18, 52, white)
			rl.DrawRectangle(int32(c.X)-26, int32(c.Y)-9, 52, 18, white)
		}),
		NewEntity("Star", func(c rl.Vector2) {
			var pts [10]rl.Vector2
			for i := range pts {
				r := 28.0
				if i%2 == 1 {
					r = 12
				}
				a := (-90 + float64(i)*36) * math.Pi / 180
				pts[i] = v(c.X+float32(math.Cos(a)*r), c.Y+float32(math.Sin(a)*r))
			}
			for i := range pts {
				tri(c, pts[i], pts[(i+1)%10], white)
			}
		}),
		NewEntity("Target", func(c rl.Vector2) {
			rl.DrawCircleV(c, 27, white)
			rl.DrawCircleV(c, 18, grey)
			rl.DrawCircleV(c, 9, white)
		}),
		NewEntity("Heart", func(c rl.Vector2) {
			rl.DrawCircleV(v(c.X-13, c.Y-8), 14, white)
			rl.DrawCircleV(v(c.X+13, c.Y-8), 14, white)
			tri(v(c.X-25, c.Y-2), v(c.X+25, c.Y-2), v(c.X, c.Y+26), white)
		}),
		NewEntity("Gem", func(c rl.Vector2) {
			rl.DrawPoly(c, 6, 28, 0, white)
			rl.DrawPolyLines(c, 6, 18, 0, grey)
			rl.DrawPolyLines(c, 6, 9, 0, grey)
		}),
		NewEntity("Arrow", func(c rl.Vector2) {
			tri(v(c.X, c.Y-28), v(c.X+24, c.Y+2), v(c.X-24, c.Y+2), white)
			rl.DrawRectangle(int32(c.X)-9, int32(c.Y), 18, 28, white)
		}),
		NewEntity("Pill", func(c rl.Vector2) {
			rl.DrawRectangleRounded(rl.NewRectangle(c.X-12, c.Y-28, 24, 56), 1, 12, white)
		}),
	}
}
