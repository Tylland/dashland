package game

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/tylland/dashland/internal/characteristics"
	"github.com/tylland/dashland/internal/common"
	"github.com/tylland/dashland/internal/components"
	"github.com/tylland/dashland/internal/ecs"
	"github.com/tylland/dashland/internal/game/types"
	"github.com/tylland/dashland/internal/utils"
)

type SoundPlayer interface {
	PlayFx(name string)
}

type Stage struct {
	Name string
	common.MapSize
	*BlockMap
	*EntityMap
	EnterPosition      common.BlockPosition
	ExitPosition       common.BlockPosition
	ExitCondition      bool
	DiamondsRequired   int
	DiamondPoints      int
	DiamondBonusPoints int
	TimeLimit          int
	TimeRemaining      float32
	TimeExpired        bool
	MagicWallActive    bool
	MagicWallTimer     float32
	MagicWallDuration  float32
	counters           *EntityCounters
}

func NewStage(name string, size common.MapSize, blockTexture, entityTextures, groundCorners *rl.Texture2D) *Stage {
	return &Stage{Name: name,
		MapSize:           size,
		BlockMap:          NewBlockMap(size, blockTexture),
		EntityMap:         NewEntityMap(size, entityTextures, groundCorners),
		MagicWallDuration: 20.0,
		counters:          NewEntityCounters(),
	}
}

func (s *Stage) Area(pos common.BlockPosition) rl.Rectangle {
	return rl.NewRectangle(float32(pos.X)*s.BlockWidth, float32(pos.Y)*s.BlockHeight, s.BlockWidth, s.BlockHeight)
}

func (s *Stage) SetTimeLimit(seconds int) {
	s.TimeLimit = seconds
	s.TimeRemaining = float32(seconds)
	s.TimeExpired = false
}

func (s *Stage) GetPosition(position common.BlockPosition) rl.Vector2 {
	return rl.NewVector2(float32(position.X)*s.BlockMap.BlockWidth, float32(position.Y)*s.BlockMap.BlockHeight)
}

func (s *Stage) GetBlockPosition(position rl.Vector2) common.BlockPosition {
	bx := int(position.X / float32(s.BlockWidth))
	by := int(position.Y / float32(s.BlockHeight))

	return common.BlockPosition{X: bx, Y: by}
}

func (s *Stage) ClampedBlockPosition(position rl.Vector2) (common.BlockPosition, bool) {
	bx := int(position.X / float32(s.BlockWidth))
	by := int(position.Y / float32(s.BlockHeight))

	x, xi := utils.ClampInside(bx, 0, s.Width-1)
	y, yi := utils.ClampInside(by, 0, s.Height-1)

	return common.NewBlockPosition(x, y), xi && yi
}

// VisitBlock stores the visited block at position and returns the block it
// replaced. The replaced block is nil when the cell was empty.
func (s *Stage) VisitBlock(block *ecs.Entity, position common.BlockPosition) *ecs.Entity {
	if old, ok := s.GetBlockAtPosition(position); ok {
		fmt.Printf("Block at position %d,%d changed from %s", position.X, position.Y, types.EntityTypeName(old.Type))
	}

	replaced := s.BlockMap.SetBlock(block, position)
	fmt.Printf(" to %s \n", types.EntityTypeName(block.Type))
	return replaced
}

func (s *Stage) ActivateMagicWall() {
	if !s.MagicWallActive {
		s.MagicWallActive = true
		s.MagicWallTimer = s.MagicWallDuration
		fmt.Println("Magic wall activated!")
	}
}

func (s *Stage) UpdateMagicWall(deltaTime float32) {
	if s.MagicWallActive {
		s.MagicWallTimer -= deltaTime
		if s.MagicWallTimer <= 0 {
			s.MagicWallActive = false
			s.MagicWallTimer = 0
			fmt.Println("Magic wall expired!")
		}
	}
}

func (s *Stage) IsMagicWallActive() bool {
	return s.MagicWallActive
}

func (s *Stage) CheckCharacteristics(position common.BlockPosition, character characteristics.Characteristics) bool {
	if block, ok := s.GetBlock(position.X, position.Y); ok {
		if component, ok := block.Component[components.CharacteristicComponent](); ok {
			if component.Has(character) {
				return true
			}
		}
	}

	if entity, ok := s.GetEntityAtPosition(position); ok {
		characteristics := ecs.GetComponent[components.CharacteristicComponent](entity)

		return characteristics != nil && characteristics.Has(character)
	}

	return false
}

func (s *Stage) CheckBlocked(world *ecs.World, position common.BlockPosition, collider *components.ColliderComponent) bool {

	if block, ok := s.GetBlock(position.X, position.Y); ok {
		blockCollider := ecs.GetComponent[components.ColliderComponent](block)
		if blockCollider != nil {
			obstructs, _ := collider.Collide(blockCollider)
			if obstructs {
				return true
			}
		}
	}

	if entity, ok := s.GetEntityAtPosition(position); ok {
		entCharacter := ecs.GetComponent[components.CharacteristicComponent](entity)
		entCollider := ecs.GetComponent[components.ColliderComponent](entity)

		if entCharacter != nil && entCollider != nil {
			// Falling entities are not considered obstacles
			if entCharacter.Has(characteristics.Falling) {
				return false
			}
			obstructs, _ := collider.Collide(entCollider)
			if obstructs {
				return true
			}
		}
	}

	return false
}

func (s *Stage) CheckPositionOccupied(position common.BlockPosition) bool {
	if !s.CheckBlockAtPosition(types.EntityVoid, position) {
		return true
	}

	return s.GetEntity(position) != nil
}

func (s *Stage) Save(path string) ([]string, error) {
	start := time.Now()
	counters := NewEntityCounters()

	writer, err := NewTextFileWriter(path + "/" + "stage.txt")

	if err != nil {
		return nil, err
	}

	s.WriteTo(writer)

	writer.Close()

	for _, b := range s.blocks {
		if b == nil {
			continue
		}

		fileName := BlockFileName(b)

		writer, err := NewTextFileWriter(path + "/" + fileName)

		if err != nil {
			return nil, err
		}

		// if err := WriteBlockTo(b, writer); err != nil {
		// 	writer.Close()
		// 	return nil, err
		// }

		b.WriteTo(writer)

		writer.Close()

		counters.Increment(types.EntityTypeName(b.Type))
	}

	for _, e := range s.entities {
		if e == nil {
			continue
		}

		fileName := EntityFileName(e)
		writer, err := NewTextFileWriter(path + "/" + fileName)

		if err != nil {
			return nil, err
		}

		e.WriteTo(writer)

		writer.Close()

		counters.Increment(types.EntityTypeName(e.Type))
	}

	buf := &strings.Builder{}

	fmt.Fprintf(buf, "\nSaving stage %q in directory %q\n", s.Name, path)
	tw := tabwriter.NewWriter(buf, 0, 0, 2, ' ', 0)

	total := 1

	fmt.Fprintf(tw, "%s\t:\t%d\t\n", "stage", 1)

	for name, count := range counters.Counters {
		fmt.Fprintf(tw, "%s\t:\t%d\t\n", name, count)
		total += count
	}
	fmt.Fprintf(tw, "%s\t:\t%d\t\n", "total", total)

	tw.Flush()

	fmt.Fprintf(buf, "Elapsed time %s\n\n", time.Since(start))

	fmt.Print(buf.String())

	return strings.Split(buf.String(), "\n"), nil
}

func (s *Stage) Load(world *ecs.World, path string) error {
	start := time.Now()
	counters := NewEntityCounters()

	reader, err := NewTextFileReader(path + "/" + "stage.txt")

	if err != nil {
		return err
	}

	s.ReadFrom(reader)

	//reader.Close()

	fileNames, err := filepath.Glob(filepath.Join(path, "*.entity.txt"))

	if err != nil {
		return err
	}

	for _, fileName := range fileNames {
		blockPosition := common.NewBlockPosition(0, 0)

		reader, err := NewTextFileReader(fileName)

		if err != nil {
			return err
		}

		header := &ecs.EntityHeader{}
		header.ReadFrom(reader)
		reader.Rewind()

		entity, err := NewEntity(header, world, s, blockPosition)
		if err != nil {
			return fmt.Errorf("creating block %q from %q: %w", types.EntityTypeName(header.Type), fileName, err)
		}
		if err := entity.ReadFrom(reader); err != nil {
			world.RemoveEntity(entity)
			return fmt.Errorf("reading entity %q from %q: %w", types.EntityTypeName(header.Type), fileName, err)
		}

		if position, ok := entity.Component[components.PositionComponent](); ok {
			if header.Category == types.EntityCategoryBlock {
				s.BlockMap.SetBlock(entity, position.CurrentBlockPosition)
			} else {
				s.EntityMap.SetEntity(entity, position.CurrentBlockPosition)
			}
		}

		counters.Increment(types.EntityTypeName(header.Type))
	}

	fmt.Printf("\nLoading stage %q from directory %q\n", s.Name, path)
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	total := 1

	fmt.Fprintf(tw, "%s\t:\t%d\t\n", "stage", 1)

	for name, count := range counters.Counters {
		fmt.Fprintf(tw, "%s\t:\t%d\t\n", name, count)
		total += count
	}
	fmt.Fprintf(tw, "%s\t:\t%d\t\n", "total", total)

	tw.Flush()
	fmt.Printf("Elapsed time %s\n\n", time.Since(start))

	return nil
}

func (s *Stage) WriteTo(w ecs.Writer) error {
	const version = 1

	w.StartSection("stage", version)
	w.WriteString("Name", s.Name, false)

	s.MapSize.WriteTo(w)

	w.StartSection("enter position", version)
	w.WriteInt("X", s.EnterPosition.X, false)
	w.WriteInt("Y", s.EnterPosition.Y, false)

	w.StartSection("exit position", version)
	w.WriteInt("X", s.ExitPosition.X, false)
	w.WriteInt("Y", s.ExitPosition.Y, false)

	w.StartSection("exit condition", version)
	w.WriteBool("Condition", s.ExitCondition, false)

	w.StartSection("diamonds", version)
	w.WriteInt("Required", s.DiamondsRequired, false)
	w.WriteInt("Points", s.DiamondPoints, false)
	w.WriteInt("Bonus Points", s.DiamondBonusPoints, false)

	w.StartSection("time", version)
	w.WriteInt("Limit", s.TimeLimit, false)
	w.WriteFloat("Remaining", s.TimeRemaining, false)
	w.WriteBool("Expired", s.TimeExpired, false)

	w.StartSection("magic wall", version)
	w.WriteBool("Active", s.MagicWallActive, false)
	w.WriteFloat("Diration", s.MagicWallDuration, false)
	w.WriteFloat("Timer", s.MagicWallTimer, false)

	s.counters.WriteTo(w)

	return nil
}

func (s *Stage) ReadFrom(r ecs.Reader) error {
	r.ConsumeSection("stage")

	if err := ecs.WithChangedValue(r.ReadString, func(name string) {
		s.Name = name
	}); err != nil {
		return err
	}

	if err := s.MapSize.ReadFrom(r); err != nil {
		return err
	}

	if s.BlockMap.Width != s.MapSize.Width || s.BlockMap.Height != s.MapSize.Height {
		s.BlockMap = NewBlockMap(s.MapSize, s.blockTextures)
	}

	if s.EntityMap.Width != s.MapSize.Width || s.EntityMap.Height != s.MapSize.Height {
		s.EntityMap = NewEntityMap(s.MapSize, s.entityTextures, s.groundCorners)
	}

	r.ConsumeSection("enter position")
	if err := ecs.WithChangedValue(r.ReadInt, func(x int) {
		s.EnterPosition.X = x
	}); err != nil {
		return err
	}
	if err := ecs.WithChangedValue(r.ReadInt, func(y int) {
		s.EnterPosition.Y = y
	}); err != nil {
		return err
	}

	r.ConsumeSection("exit position")
	if err := ecs.WithChangedValue(r.ReadInt, func(x int) {
		s.ExitPosition.X = x
	}); err != nil {
		return err
	}
	if err := ecs.WithChangedValue(r.ReadInt, func(y int) {
		s.ExitPosition.Y = y
	}); err != nil {
		return err
	}

	r.ConsumeSection("exit condition")
	if err := ecs.WithChangedValue(r.ReadBool, func(exitCondition bool) {
		s.ExitCondition = exitCondition
	}); err != nil {
		return err
	}

	r.ConsumeSection("diamonds")
	if err := ecs.WithChangedValue(r.ReadInt, func(diamondsRequired int) {
		s.DiamondsRequired = diamondsRequired
	}); err != nil {
		return err
	}
	if err := ecs.WithChangedValue(r.ReadInt, func(diamondPoints int) {
		s.DiamondPoints = diamondPoints
	}); err != nil {
		return err
	}
	if err := ecs.WithChangedValue(r.ReadInt, func(diamondBonusPoints int) {
		s.DiamondBonusPoints = diamondBonusPoints
	}); err != nil {
		return err
	}

	r.ConsumeSection("time")
	if err := ecs.WithChangedValue(r.ReadInt, func(timeLimit int) {
		s.TimeLimit = timeLimit
	}); err != nil {
		return err
	}
	if err := ecs.WithChangedValue(r.ReadFloat, func(timeRemaining float32) {
		s.TimeRemaining = timeRemaining
	}); err != nil {
		return err
	}
	if err := ecs.WithChangedValue(r.ReadBool, func(timeExpired bool) {
		s.TimeExpired = timeExpired
	}); err != nil {
		return err
	}

	r.ConsumeSection("magic wall")
	if err := ecs.WithChangedValue(r.ReadBool, func(magicWallActive bool) {
		s.MagicWallActive = magicWallActive
	}); err != nil {
		return err
	}
	if err := ecs.WithChangedValue(r.ReadFloat, func(magicWallDuration float32) {
		s.MagicWallDuration = magicWallDuration
	}); err != nil {
		return err
	}
	if err := ecs.WithChangedValue(r.ReadFloat, func(magicWallTimer float32) {
		s.MagicWallTimer = magicWallTimer
	}); err != nil {
		return err
	}

	s.counters.ReadFrom(r)

	return nil
}
