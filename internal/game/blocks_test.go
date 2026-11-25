package game

import (
	"bytes"
	"strings"
	"testing"

	"github.com/tylland/dashland/internal/characteristics"
	"github.com/tylland/dashland/internal/common"
	"github.com/tylland/dashland/internal/components"
	"github.com/tylland/dashland/internal/ecs"
)

func TestBlockIsAnEntityAndRoundTrips(t *testing.T) {
	size := NewMapSize(2, 2, 16, 16)
	stage := NewStage("test", size, nil, nil, nil)
	world := ecs.NewWorld()
	position := common.NewBlockPosition(1, 0)
	block, err := NewEntityAndHeader(EntitySoil, world, stage, position)
	if err != nil {
		t.Fatal(err)
	}
	blockComponent := ecs.GetComponent[components.BlockComponent](block)
	blockComponent.Corners = [4]uint8{1, 2, 3, 4}

	stage.BlockMap.SetBlock(block, position)

	if got := len(world.Entities()); got != 1 || world.Entities()[0] != block {
		t.Fatalf("world contains %d entities; want the block entity", got)
	}
	if got, ok := stage.GetBlockAtPosition(position); !ok || got != block {
		t.Fatal("block map does not point to the block entity")
	}
	if !BlockHasCharacteristic(block, characteristics.Destructable) {
		t.Fatal("soil entity is missing its destructable characteristic")
	}

	var data bytes.Buffer
	if err := WriteBlockTo(block, NewTextWriter(&data)); err != nil {
		t.Fatalf("write block: %v", err)
	}

	header := NewEntityHeader(EntitySoil)
	loaded, err := NewEntity(header, world, stage, common.NewBlockPosition(0, 0))
	if err != nil {
		t.Fatal(err)
	}
	if err := ReadBlockFrom(loaded, NewTextReader(strings.NewReader(data.String())), stage.BlockMap); err != nil {
		t.Fatalf("read block: %v", err)
	}

	loadedPosition := ecs.GetComponent[components.PositionComponent](loaded)
	loadedBlock := ecs.GetComponent[components.BlockComponent](loaded)
	if loadedPosition.CurrentBlockPosition != position {
		t.Fatalf("loaded position is %v; want %v", loadedPosition.CurrentBlockPosition, position)
	}
	if loadedBlock.Corners != blockComponent.Corners {
		t.Fatalf("loaded corners are %v; want %v", loadedBlock.Corners, blockComponent.Corners)
	}
	if loaded.Type != block.Type || loaded.ID != block.ID {
		t.Fatalf("loaded header is %q/%d; want %q/%d", loaded.ID, loaded.Type, block.ID, block.Type)
	}
}

func TestSetBlockReplacesWorldEntity(t *testing.T) {
	size := NewMapSize(1, 1, 16, 16)
	stage := NewStage("test", size, nil, nil, nil)
	world := ecs.NewWorld()
	position := common.NewBlockPosition(0, 0)
	first, err := NewEntityAndHeader(EntitySoil, world, stage, position)
	if err != nil {
		t.Fatal(err)
	}
	second, err := NewEntityAndHeader(EntityVoid, world, stage, position)
	if err != nil {
		t.Fatal(err)
	}

	stage.BlockMap.SetBlock(first, position)
	old := stage.BlockMap.SetBlock(second, position)
	if old != nil {
		world.RemoveEntity(old)
	}

	if got := world.Entities(); len(got) != 1 || got[0] != second {
		t.Fatalf("world entities after replacement = %v; want only replacement", got)
	}
	if got, ok := stage.GetBlockAtPosition(position); !ok || got != second {
		t.Fatal("block map does not point to the replacement entity")
	}
}
