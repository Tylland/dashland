package game

import (
	"strings"
	"testing"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/tylland/dashland/internal/common"
	"github.com/tylland/dashland/internal/components"
	"github.com/tylland/dashland/internal/ecs"
)

func TestReader(t *testing.T) {

	ent := ecs.NewEntity(NewEntityId(""), EntityUnknown)
	ent.AddComponent(components.NewPositionComponent(common.NewBlockPosition(10, 10), rl.NewVector2(320.0, 329.0)))

	reader := NewTextReader(strings.NewReader(saved))

	ent.ReadFrom(reader)

	if ent.ID != "test-1" {
		t.Fatalf("<-- ERROR")
	}

	if ent.Type != EntityDoor {
		t.Fatalf("<-- ERROR")
	}
}
