package game

import (
	"testing"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/tylland/dashland/game/characteristics"
	"github.com/tylland/dashland/game/core"
)

type mockWorld struct {
	MapSize
	obstacles map[core.BlockPosition]bool
}

func (m *mockWorld) CheckCharacteristics(pos core.BlockPosition, c characteristics.Characteristics) bool {
	return m.obstacles[pos]
}

func (m *mockWorld) GetPosition(pos core.BlockPosition) rl.Vector2 {
	return rl.NewVector2(float32(pos.X)*m.blockWidth, float32(pos.Y)*m.blockHeight)
}

func (m *mockWorld) Entities() []*Entity {
	return []*Entity{}
}

func TestFindDirection_NoObstacles(t *testing.T) {
	w := &mockWorld{obstacles: make(map[core.BlockPosition]bool)}
	s := &WallWalkerSystem{world: w}
	start := core.BlockPosition{X: 0, Y: 0}
	dir, ok := s.findDirection(start, core.BlockVector{X: 1, Y: 0})
	if !ok || dir != (core.BlockVector{X: 1, Y: 0}) {
		t.Errorf("Expected to find direction right, got %v, ok=%v", dir, ok)
	}
}

func TestFindDirection_AllObstacles(t *testing.T) {
	w := &mockWorld{obstacles: map[core.BlockPosition]bool{
		{X: 1, Y: 0}:  true,
		{X: 0, Y: 1}:  true,
		{X: -1, Y: 0}: true,
		{X: 0, Y: -1}: true,
	}}
	s := &WallWalkerSystem{world: w}
	start := core.BlockPosition{X: 0, Y: 0}
	_, ok := s.findDirection(start, core.BlockVector{X: 1, Y: 0})
	if ok {
		t.Errorf("Expected not to find any direction, but got ok=true")
	}
}

// More tests can be added for UpdateTarget logic, including wall following and turning.
