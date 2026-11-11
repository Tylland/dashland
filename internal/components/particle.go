package components

import (
	"math/rand/v2"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Particle struct {
	X     float32
	Y     float32
	Vx    float32
	Vy    float32
	Size  float32
	Alpha float32
	Phase float32
}

type ParticleComponent struct {
	Type  string
	Parts []*Particle
	Area  rl.Rectangle
	Speed float32
	Tint  rl.Color
}

func NewParticleComponent(pt string, area rl.Rectangle, count int, speed, sizeMin, sizeMax, alpha float32, tint rl.Color) *ParticleComponent {
	var parts []*Particle

	// Mist particles
	for i := 0; i < count; i++ {
		parts = append(parts, &Particle{
			X:     area.X + rand.Float32()*area.Width,
			Y:     area.Y + rand.Float32()*area.Height,
			Size:  sizeMin + rand.Float32()*(sizeMax-sizeMin),
			Alpha: alpha * (0.6 + rand.Float32()*0.6),
			Phase: rand.Float32() * 6.28,
		})
	}

	return &ParticleComponent{Type: pt, Parts: parts, Area: area, Speed: speed, Tint: tint}
}
