package systems

import (
	"math"

	"github.com/tylland/dashland/internal/components"
	"github.com/tylland/dashland/internal/ecs"
)

type ParticleSystem struct {
	elapsedTime float32
}

func NewParticleSystem() *ParticleSystem {
	return &ParticleSystem{}
}
func (s *ParticleSystem) Update(world *ecs.World, deltaTime float32) {
	s.elapsedTime += deltaTime

	for _, entity := range world.Entities() {
		particle := ecs.GetComponent[components.ParticleComponent](entity)

		if particle != nil {
			s.updateMist(particle, deltaTime, 0)
		}
	}
}

// Smooth pseudo-noise flow field -> angle in radians.
func flow(x, y, t float64) float64 {
	a := math.Sin(x*0.004+t*0.15)*1.2 + math.Cos(y*0.006-t*0.10)*1.2
	return a
}

func (s *ParticleSystem) updateMist(particle *components.ParticleComponent, dt float32, wind float32) {
	for _, p := range particle.Parts {

		ang := flow(float64(p.X), float64(p.X), float64(s.elapsedTime))

		tx := float32(math.Cos(ang))*10*particle.Speed + wind*particle.Speed
		ty := float32(math.Sin(ang)) * 3

		p.Vx += (tx - p.Vx) * dt * 0.6
		p.Vy += (ty - p.Vy) * dt * 0.6

		p.X += p.Vx * dt
		p.Y += p.Vy * dt

		// // Keep inside the band
		// if p.y < l.yMin {
		// 	p.vy += 25 * dt
		// }
		// if p.y > l.yMax {
		// 	p.vy -= 25 * dt
		// }

		// // Wrap horizontally
		// if p.x > W+p.size {
		// 	p.x = -p.size
		// } else if p.x < -p.size {
		// 	p.x = W + p.size
		// }
	}
}
