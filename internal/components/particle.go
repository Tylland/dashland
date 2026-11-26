package components

import (
	"math/rand/v2"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/tylland/dashland/internal/common"
)

type Particle struct {
	Type  string
	Parts []*common.Particle
	Area  rl.Rectangle
	Speed float32
	Tint  rl.Color
}

func NewParticle(pt string, area rl.Rectangle, count int, speed, sizeMin, sizeMax, alpha float32, tint rl.Color) *Particle {
	var parts []*common.Particle

	// Mist particles
	for i := 0; i < count; i++ {
		parts = append(parts, &common.Particle{
			X:     area.X + rand.Float32()*area.Width,
			Y:     area.Y + rand.Float32()*area.Height,
			Size:  sizeMin + rand.Float32()*(sizeMax-sizeMin),
			Alpha: alpha * (0.6 + rand.Float32()*0.6),
			Phase: rand.Float32() * 6.28,
		})
	}

	return &Particle{Type: pt, Parts: parts, Area: area, Speed: speed, Tint: tint}
}
