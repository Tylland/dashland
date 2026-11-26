package game

import (
	"strings"
	"testing"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/tylland/dashland/internal/common"
	"github.com/tylland/dashland/internal/components"
	"github.com/tylland/dashland/internal/ecs"
	"github.com/tylland/dashland/internal/game/types"
)

func TestReader(t *testing.T) {

	ent := ecs.NewEntity(NewEntityHeader(types.EntityUnknown))
	ent.AddComponent(components.NewPosition(common.NewBlockPosition(10, 10), rl.NewVector2(320.0, 329.0)))

	reader := NewTextReader(strings.NewReader(saved))

	ent.ReadFrom(reader)

	if ent.ID != "exitdoor-1" {
		t.Fatalf("<-- ERROR")
	}

	if ent.Type != types.EntityExitDoor {
		t.Fatalf("<-- ERROR")
	}
}
