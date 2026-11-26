package components

import (
	"github.com/tylland/dashland/internal/common"
	"github.com/tylland/dashland/internal/ecs"
	"github.com/tylland/dashland/internal/game/types"
)

// Animation holds a set of named animations and which is active.
type Animation struct {
	Animations map[string]common.Animation
	Current    string // name of the current animation
	active     string // internal: last-applied animation name
	// runtime state
	Frame      uint
	FrameTimer float32
}

func NewAnimation(anims map[string]common.Animation, initial string) *Animation {
	return &Animation{Animations: anims, Current: initial}
}

// ApplyAnimation copies animation properties into the given sprite and resets frame if switching.
func (ac *Animation) ApplyAnimation(s *common.Sprite) {
	if s == nil || ac == nil || ac.Current == "" {
		return
	}

	anim, ok := ac.Animations[ac.Current]
	if !ok {
		return
	}

	if ac.active != ac.Current {
		// switching animation -> update sprite and reset frame timer/frame
		s.SourceBase.X = anim.BaseX
		s.SourceBase.Y = anim.BaseY
		ac.Frame = 0
		ac.FrameTimer = 0
		s.UpdateFrame(ac.Frame)
		ac.active = ac.Current
	} else {
		// ensure base is kept in sync
		s.SourceBase.X = anim.BaseX
		s.SourceBase.Y = anim.BaseY
	}
}

// Advance progresses the animation timer and updates the sprite frame when needed.
func (ac *Animation) Advance(deltaTime float32, s *common.Sprite) {
	if ac == nil || s == nil || ac.Current == "" {
		return
	}

	anim, ok := ac.Animations[ac.Current]
	if !ok || anim.FrameCount <= 1 || anim.FrameDuration <= 0 {
		return
	}

	ac.FrameTimer += deltaTime
	if ac.FrameTimer >= anim.FrameDuration {
		ac.FrameTimer -= anim.FrameDuration

		if anim.Loop {
			ac.Frame = (ac.Frame + 1) % anim.FrameCount
		} else {
			ac.Frame = clamp(ac.Frame+1, anim.FrameCount-1)
		}

		s.UpdateFrame(ac.Frame)
	}
}

func clamp(val, max uint) uint {
	if val < max {
		return val
	}
	return max
}

func (c *Animation) WriteTo(w ecs.Writer) error {
	w.StartSection(Section(types.ComponentTypeAnimation))
	w.WriteInt("Animations", len(c.Animations), true)
	w.WriteString("Current", c.Current, false)
	w.WriteString("Active", c.active, false)
	w.WriteInt("Frame", int(c.Frame), false)
	w.WriteFloat("Frame Timer", c.FrameTimer, false)

	return nil
}

func (c *Animation) ReadFrom(r ecs.Reader) error {
	r.ConsumeSection(types.ComponentTypeAnimation.Name)
	r.ConsumeLine()
	r.ConsumeLine()
	r.ConsumeLine()
	r.ConsumeLine()
	r.ConsumeLine()

	return nil
}
