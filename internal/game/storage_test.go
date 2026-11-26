package game

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/tylland/dashland/internal/characteristics"
	"github.com/tylland/dashland/internal/common"
	"github.com/tylland/dashland/internal/components"
	"github.com/tylland/dashland/internal/ecs"
	"github.com/tylland/dashland/internal/game/types"
	"github.com/tylland/dashland/specs"
)

func TestPositionComponentTextRoundtrip(t *testing.T) {
	var buffer bytes.Buffer
	writer := NewTextWriter(&buffer)

	blockPosition := common.NewBlockPosition(123, 456)
	position := rl.NewVector2(123.1, 456.1)

	source := components.NewPositionComponent(blockPosition, position)

	source.WriteTo(writer)

	reader := NewTextReader(&buffer)

	target := components.NewPositionComponentZero()

	target.ReadFrom(reader)

	if !target.CurrentBlockPosition.IsSame(blockPosition) {
		t.Fatalf("asdfasdas")
	}

	if target.Vector2.X != position.X || target.Vector2.Y != position.Y {
		t.Fatalf("asdfasdas")
	}

	if err := reader.ConsumeLine(); err != io.EOF {
		t.Fatal("Expected to reached end of file")
	}
}

func TestCharacteristicsComponentTextRoundtrip(t *testing.T) {
	var buffer bytes.Buffer
	writer := NewTextWriter(&buffer)

	character := characteristics.CanFall | characteristics.Pushable

	source := components.NewCharacteristicsComponent(character)

	source.WriteTo(writer)

	reader := NewTextReader(&buffer)

	target := components.NewCharacteristicsComponent(characteristics.None)

	target.ReadFrom(reader)

	if !target.Has(character) {
		t.Fatalf("asdfasdas")
	}

	if err := reader.ConsumeLine(); err != io.EOF {
		t.Fatal("Expected to reached end of file")
	}
}

func TestBoulderTextRoundtrip(t *testing.T) {
	var buffer bytes.Buffer
	writer := NewTextWriter(&buffer)

	blockPosition := common.NewBlockPosition(123, 456)

	world := ecs.NewWorld()
	stage := NewStage("test", common.NewMapSize(10, 10, 32, 32), nil, nil, nil)

	source := specs.Must(NewEntityAndHeader(types.EntityBoulder, world, stage, blockPosition))
	source.WriteTo(writer)

	fmt.Print(buffer.String())

	reader := NewTextReader(&buffer)

	target := specs.Must(NewEntityAndHeader(types.EntityBoulder, world, stage, blockPosition))
	target.ReadFrom(reader)

	if err := reader.ConsumeLine(); err != io.EOF {
		t.Fatal("Expected to reached end of file")
	}
}
