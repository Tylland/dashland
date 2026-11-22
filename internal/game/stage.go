package game

import (
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/tylland/dashland/internal/characteristics"
	"github.com/tylland/dashland/internal/common"
	"github.com/tylland/dashland/internal/components"
	"github.com/tylland/dashland/internal/ecs"
)

type SoundPlayer interface {
	PlayFx(name string)
}

type Stage struct {
	Name string
	MapSize
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
}

func NewStage(name string, size MapSize, blockTexture, entityTextures, groundCorners *rl.Texture2D) *Stage {
	return &Stage{Name: name, MapSize: size, BlockMap: NewBlockMap(size, blockTexture), EntityMap: NewEntityMap(size, entityTextures, groundCorners), MagicWallDuration: 20.0}
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

func (s *Stage) Render(deltaTime float32) {
	for _, block := range s.blocks {
		block.Render()
	}
}

func (s *Stage) VisitBlock(position common.BlockPosition) {
	fmt.Printf("Block at position %d,%d changed from %s", position.X, position.Y, ToBlockName(s.blocks[position.Y*s.Width+position.X].Type))

	header := NewEntityHeader(EntityVoid)
	s.SetBlock(NewBlock(header, s.BlockMap, s.EntityMap, position), position)
	fmt.Printf(" to %s \n", ToBlockName(s.blocks[position.Y*s.Width+position.X].Type))
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
		if block.HasCharacteristic(character) {
			return true
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
		obstructs, _ := collider.Collide(block.Collider)
		if obstructs {
			return true
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
	if !s.CheckBlockAtPosition(EntityVoid, position) {
		return true
	}

	return s.GetEntity(position) != nil
}

// func (s *Stage) SpawnExitDoor(world *ecs.World, blockPos common.BlockPosition, targetStage string, targetPos common.BlockPosition) *ecs.Entity {
// 	door := NewExitDoor(world, s, blockPos, s.GetPosition(blockPos), targetStage, targetPos)
// 	s.EntityMap.SetEntity(door, blockPos)
// 	return door
// }

func (s *Stage) Save(path string) error {
	start := time.Now()

	writer, err := NewTextFileWriter(path + "/" + "stage.txt")

	if err != nil {
		return err
	}

	s.WriteTo(writer)

	writer.Close()

	for _, b := range s.blocks {
		fileName := BlockFileName(b)

		writer, err := NewTextFileWriter(path + "/" + fileName)

		if err != nil {
			return err
		}

		b.WriteTo(writer)

		writer.Close()
	}

	for _, e := range s.entities {
		if e == nil {
			continue
		}

		fileName := EntityFileName(e)
		writer, err := NewTextFileWriter(path + "/" + fileName)

		if err != nil {
			return err
		}

		e.WriteTo(writer)

		writer.Close()
	}

	fmt.Printf("\nSaving stage %q in directory %q\n", s.Name, path)
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	total := 1

	fmt.Fprintf(tw, "%s\t:\t%d\t\n", "stage", 1)

	for name, count := range EntityIdCounters.Counters {
		fmt.Fprintf(tw, "%s\t:\t%d\t\n", name, count)
		total += count
	}
	fmt.Fprintf(tw, "%s\t:\t%d\t\n", "total", total)

	tw.Flush()
	fmt.Printf("Elapsed time %s\n\n", time.Since(start))

	return nil
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

		if header.Category == EntityCategoryBlock {
			block := NewBlock(header, s.BlockMap, s.EntityMap, blockPosition)
			block.ReadFrom(reader)
			block.EntityHeader = header
			s.BlockMap.SetBlock(block, block.Position)
		} else {
			entity, err := NewEntity(header, world, s, blockPosition)

			if err != nil {
				return err
			}

			entity.ReadFrom(reader)
			entity.EntityHeader = header
			position := ecs.GetComponent[components.PositionComponent](entity)
			s.EntityMap.SetEntity(entity, position.CurrentBlockPosition)
		}

		counters.Increment(EntityTypeName(header.Type))
	}

	fmt.Printf("\nLoding stage %q from directory %q\n", s.Name, path)
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
	w.StartSection("stage")
	w.WriteString("Name", s.Name, false)

	s.MapSize.WriteTo(w)

	w.StartSection("enter position")
	w.WriteInt("X", s.EnterPosition.X, false)
	w.WriteInt("Y", s.EnterPosition.Y, false)

	w.StartSection("exit position")
	w.WriteInt("X", s.ExitPosition.X, false)
	w.WriteInt("Y", s.ExitPosition.Y, false)

	w.StartSection("exit condition")
	w.WriteBool("Condition", s.ExitCondition, false)

	w.StartSection("diamonds")
	w.WriteInt("Required", s.DiamondsRequired, false)
	w.WriteInt("Points", s.DiamondPoints, false)
	w.WriteInt("Bonus Points", s.DiamondBonusPoints, false)

	w.StartSection("time")
	w.WriteInt("Limit", s.TimeLimit, false)
	w.WriteFloat("Remaining", s.TimeRemaining, false)
	w.WriteBool("Expired", s.TimeExpired, false)

	w.StartSection("magic wall")
	w.WriteBool("Active", s.MagicWallActive, false)
	w.WriteFloat("Diration", s.MagicWallDuration, false)
	w.WriteFloat("Timer", s.MagicWallTimer, false)

	return nil
}

func (s *Stage) ReadFrom(r ecs.Reader) error {
	r.ConsumeLabel("stage")

	if err := ecs.WithChangedValue(r.ReadString, func(name string) {
		s.Name = name
	}); err != nil {
		return err
	}

	if err := s.MapSize.ReadFrom(r); err != nil {
		return err
	}

	s.BlockMap = NewBlockMap(s.MapSize, s.blockTextures)
	s.EntityMap = NewEntityMap(s.MapSize, s.entityTextures, s.groundCorners)

	r.ConsumeLabel("enter position")
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

	r.ConsumeLabel("exit position")
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

	r.ConsumeLabel("exit condition")
	if err := ecs.WithChangedValue(r.ReadBool, func(exitCondition bool) {
		s.ExitCondition = exitCondition
	}); err != nil {
		return err
	}

	r.ConsumeLabel("diamonds")
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

	r.ConsumeLabel("time")
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

	r.ConsumeLabel("magic wall")
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

	return nil
}
