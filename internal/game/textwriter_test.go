package game

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/tylland/dashland/internal/common"
	"github.com/tylland/dashland/internal/components"
	"github.com/tylland/dashland/internal/ecs"
	"github.com/tylland/dashland/internal/game/types"
)

const saved = `# header
   exitdoor-1
   12
   0
   1 : 3f800000
# position
 * 10 : 41200000
 * 10 : 41200000
 * 320.000000 : 43a00000
 * 329.000000 : 43a48000`

func commonPrefix(expected, actual string) string {
	max := len(expected)

	if len(actual) > max {
		max = len(actual)
	}

	for i := 1; i < max; i++ {
		if !strings.HasPrefix(actual, expected[:i]) {
			return expected[:i-1]
		}
	}

	return expected[:max]
}

func TestWriter(t *testing.T) {
	EntityIdCounters.Reset()

	ent := ecs.NewEntity(NewEntityHeader(types.EntityExitDoor))
	ent.AddComponent(components.NewPositionComponent(common.NewBlockPosition(10, 10), rl.NewVector2(320.0, 329.0)))

	var buf bytes.Buffer
	writer := NewTextWriter(&buf)

	ent.WriteTo(writer)

	text := buf.String()

	if text != saved {
		t.Fatalf("%s<-- ERROR", commonPrefix(saved, text))
	}

	fmt.Print(text)
}

func TestRoundtrip(t *testing.T) {
	EntityIdCounters.Reset()
	ent := ecs.NewEntity(NewEntityHeader(types.EntityExitDoor))
	ent.AddComponent(components.NewPositionComponent(common.NewBlockPosition(10, 10), rl.NewVector2(320.0, 329.0)))

	var buf bytes.Buffer
	writer := NewTextWriter(&buf)

	ent.WriteTo(writer)

	text := buf.String()

	if text != saved {
		t.Fatalf("%s<-- ERROR", commonPrefix(saved, text))
	}

	fmt.Print(text)
}
