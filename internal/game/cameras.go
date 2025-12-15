package game

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/tylland/dashland/internal/components"
)

type Camera interface {
	GetCamera() rl.Camera2D
	Update(deltaTime float32)
}

type PositionResolver interface {
	Resolve() *components.PositionComponent
}

type Screen struct {
	Width  int
	Height int
}

type SmoothFollowCamera struct {
	camera   rl.Camera2D
	screen   *Screen
	position *components.PositionComponent
}

func NewSmoothFollowCamera(screen *Screen, position *components.PositionComponent) *SmoothFollowCamera {
	camera := SmoothFollowCamera{screen: screen, position: position}
	camera.init()

	return &camera
}

func (c *SmoothFollowCamera) init() {
	camera := rl.Camera2D{}
	camera.Target = rl.NewVector2(float32(0), float32(0))
	camera.Offset = rl.NewVector2(float32(c.screen.Width/2), float32(c.screen.Height/2))
	camera.Rotation = 0.0
	camera.Zoom = 2.0

	c.camera = camera
}

func (c *SmoothFollowCamera) GetCamera() rl.Camera2D {
	return c.camera
}

func (c *SmoothFollowCamera) Update(deltaTime float32) {
	var minSpeed, minEffectLength, fractionSpeed float32 = 30.0, 10.0, 0.8

	diff := rl.Vector2Subtract(c.position.Vector2, c.camera.Target)
	length := rl.Vector2Length(diff)

	if length > minEffectLength {
		speed := max(fractionSpeed*length, minSpeed)
		c.camera.Target = rl.Vector2Add(c.camera.Target, rl.Vector2Scale(diff, speed*deltaTime/length))
	}
}
