package components

import (
	"github.com/tylland/dashland/internal/ecs"
	"github.com/tylland/dashland/internal/game/types"
)

type BlockComponent struct {
	HasCorners bool     `version:"2"`
	Corners    [4]uint8 `version:"1"`
}

func (c *BlockComponent) WriteTo(w ecs.Writer) error {
	w.StartSection(Section(types.ComponentTypeBlock))

	w.WriteBool("Has corners", c.HasCorners, true)
	w.WriteInt("Upper left", int(c.Corners[0]), c.Corners[0] == 0)
	w.WriteInt("Upper right", int(c.Corners[1]), c.Corners[1] == 0)
	w.WriteInt("Lower left", int(c.Corners[2]), c.Corners[2] == 0)
	w.WriteInt("Lower right", int(c.Corners[3]), c.Corners[3] == 0)

	return nil
}

func (c *BlockComponent) ReadFrom(r ecs.Reader) error {
	ver := r.ConsumeSection(types.ComponentTypeBlock.Name)

	if ver >= 2 {
		if err := ecs.WithChangedValue(r.ReadBool, func(value bool) {
			c.HasCorners = value
		}); err != nil {
			return err
		}
	}

	for i := range c.Corners {
		if err := ecs.WithChangedValue(r.ReadInt, func(value int) { c.Corners[i] = uint8(value) }); err != nil {
			return err
		}
	}

	return nil
}
