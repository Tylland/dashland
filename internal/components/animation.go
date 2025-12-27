package components

import "github.com/tylland/dashland/internal/common"

// Animation describes a named animation sequence inside a sprite sheet.
type Animation struct {
	BaseX         float32
	BaseY         float32
	FrameCount    uint
	FrameDuration float32
	Loop          bool
}

// AnimationComponent holds a set of named animations and which is active.
type AnimationComponent struct {
	Animations map[string]Animation
	Current    string // name of the current animation
	active     string // internal: last-applied animation name
	// runtime state
	Frame      uint
	FrameTimer float32
}

func NewAnimationComponent(anims map[string]Animation, initial string) *AnimationComponent {
	return &AnimationComponent{Animations: anims, Current: initial}
}

// ApplyAnimation copies animation properties into the given sprite and resets frame if switching.
func (ac *AnimationComponent) ApplyAnimation(s *common.Sprite) {
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
func (ac *AnimationComponent) Advance(deltaTime float32, s *common.Sprite) {
	if ac == nil || s == nil || ac.Current == "" {
		return
	}

	anim, ok := ac.Animations[ac.Current]
	if !ok || anim.FrameCount <= 1 || anim.FrameDuration <= 0 {
		return
	}

	ac.FrameTimer += deltaTime
	if ac.FrameTimer >= anim.FrameDuration {
		ac.FrameTimer = 0
		ac.Frame = (ac.Frame + 1) % anim.FrameCount
		s.UpdateFrame(ac.Frame)
	}
}
