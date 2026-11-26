package gui

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/tylland/dashland/internal/ecs"
	"github.com/tylland/dashland/internal/game/types"
)

// EntityItem is one selectable item. Its look is a rl.Texture2D, which is either
// rendered procedurally (NewEntity) or loaded from an image file (LoadEntity).
type EntityItem struct {
	Name    string
	Type    ecs.EntityType
	Texture *rl.Texture2D
	Source  rl.Rectangle
	Flipped bool // render-target textures are stored upside down by OpenGL

	//	target rl.RenderTexture2D // owned render target (zero value when loaded from a file)
}

// NewEntity renders draw() into a fresh texture. The shapes are drawn in white/grey
// so that a tint color can be applied when the entity is displayed.
// Must be called after InitWindow and outside BeginDrawing/EndDrawing.
func NewEntityItem(entityType ecs.EntityType, texture *rl.Texture2D, source rl.Rectangle) *EntityItem {
	return &EntityItem{
		Name:    types.EntityTypeName(entityType),
		Type:    entityType,
		Texture: texture,
		Source:  source,
		Flipped: true}
}

// Draw renders the entity's texture into dest, multiplied by tint.
func (e *EntityItem) Draw(dest rl.Rectangle, tint rl.Color) {
	h := float32(e.Texture.Height)
	if e.Flipped {
		h = -h
	}

	rl.DrawTexturePro(*e.Texture, e.Source, dest, rl.Vector2Zero(), 0, rl.White)
}

// tri draws a triangle regardless of vertex winding order.
func tri(a, b, c rl.Vector2, col rl.Color) {
	rl.DrawTriangle(a, b, c, col)
	rl.DrawTriangle(a, c, b, col)
}

type EntityInspection struct {
	Entity   *ecs.Entity
	Block    *ecs.Entity
	Document *Document
	Openpath string
	Preview  []string
}
