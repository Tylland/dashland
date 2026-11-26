package game

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/tylland/dashland/internal/common"
	"github.com/tylland/dashland/internal/components"
	"github.com/tylland/dashland/internal/ecs"
	"github.com/tylland/dashland/internal/game/entity"
	"github.com/tylland/dashland/internal/game/types"
)

type EntityCounters struct {
	Counters map[string]int
}

func NewEntityCounters() *EntityCounters {
	return &EntityCounters{Counters: make(map[string]int)}
}

func (ec *EntityCounters) Increment(prefix string) int {
	ec.Counters[prefix]++

	return ec.Counters[prefix]
}

func (ec *EntityCounters) Reset() {
	ec.Counters = make(map[string]int)
}

func (ec *EntityCounters) WriteTo(w ecs.Writer) error {
	w.StartSection("counters", 1)

	w.WriteInt("Count", len(ec.Counters), false)

	for k, v := range ec.Counters {
		w.WriteString("Key", k, false)
		w.WriteInt("Value", v, false)
	}

	return nil
}

func (ec *EntityCounters) ReadFrom(r ecs.Reader) error {
	r.ConsumeSection("counters")

	ec.Reset()

	count, _, err := r.ReadInt()

	if err != nil {
		return err
	}

	for i := 0; i < count; i++ {
		key, _, _ := r.ReadString()
		val, _, _ := r.ReadInt()
		ec.Counters[key] = val
	}

	return nil
}

var EntityIdCounters *EntityCounters = NewEntityCounters()

func NewEntityId(prefix string) ecs.EntityID {
	count := EntityIdCounters.Increment(prefix)

	return ecs.EntityID(fmt.Sprintf("%s-%d", prefix, count))
}

func NewEntityHeader(entityType ecs.EntityType) *ecs.EntityHeader {
	prefix := types.EntityTypeName(entityType)
	count := EntityIdCounters.Increment(prefix)

	return ecs.NewEntityHeader(ecs.EntityID(fmt.Sprintf("%s-%d", prefix, count)), entityType, types.EntityCategory(entityType))
}

func NewEntity(header *ecs.EntityHeader, world *ecs.World, stage *Stage, bpos common.BlockPosition) (*ecs.Entity, error) {
	ent, err := newEntity(header, stage.MapSize, bpos, stage.GetPosition(bpos))

	if err == nil {
		world.AddEntity(ent)
	}

	return ent, err
}

func newEntity(header *ecs.EntityHeader, size common.MapSize, bpos common.BlockPosition, pos rl.Vector2) (*ecs.Entity, error) {
	switch header.Type {
	case types.EntityVoid:
		return entity.NewVoidBlock(header, size, bpos, pos)
	case types.EntitySoil:
		return entity.NewSoilBlock(header, size, bpos, pos)
	case types.EntityWall:
		return entity.NewWallBlock(header, size, bpos, pos)
	case types.EntityBedrock:
		return entity.NewBedrockBlock(header, size, bpos, pos)
	case types.EntityDiamond:
		return entity.NewDiamond(header, size, bpos, pos), nil
	case types.EntityBoulder:
		return entity.NewBoulder(header, size, bpos, pos), nil
	case types.EntityExplosion:
		return entity.NewExplosion(header, size, bpos, pos), nil
	case types.EntityMagicWall:
		return entity.NewMagicWall(header, size, bpos, pos), nil
	case types.EntityExitDoor:
		return entity.NewExitDoor(header, size, bpos, pos, "", bpos), nil
	case types.EntityFirefly:
		return entity.NewFirefly(header, size, bpos, pos), nil
	case types.EntityButterfly:
		return entity.NewButterfly(header, size, bpos, pos), nil
	default:
		return nil, fmt.Errorf("%q (%d) is unknown EntityType", types.EntityTypeName(header.Type), int(header.Type))
	}
}

func NewEntityAndHeader(entityType ecs.EntityType, world *ecs.World, stage *Stage, blockPosition common.BlockPosition) (*ecs.Entity, error) {
	header := NewEntityHeader(entityType)
	return NewEntity(header, world, stage, blockPosition)
}

func NewFlash(world *ecs.World) *ecs.Entity {
	entity := ecs.NewEntity(NewEntityHeader(types.EntityFlash))

	entity.AddComponent(components.NewFlash())
	entity.AddComponent(components.NewLifetime(0.1))

	world.AddEntity(entity)

	return entity
}
