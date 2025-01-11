package game

import rl "github.com/gen2brain/raylib-go/raylib"

type Movement struct {
	startPos rl.Vector2
	endPos   rl.Vector2
	speed    float32
	vector   rl.Vector2
	length   float32
	progress float32
	position rl.Vector2
	moving   bool
}

func (m *Movement) Start(startPos rl.Vector2, endPos rl.Vector2, speed float32) {
	m.startPos = startPos
	m.endPos = endPos
	m.vector = rl.Vector2Subtract(endPos, startPos)
	m.length = rl.Vector2Length(m.vector)
	m.speed = speed
	m.progress = 0.0
	m.moving = true
}

func (m *Movement) Update(deltaTime float32) {

	if m.moving {
		m.progress += m.speed * deltaTime

		progressVector := rl.Vector2Scale(m.vector, rl.Clamp(m.progress, 0, m.length)/m.length)

		m.position = rl.Vector2Add(m.startPos, progressVector)

		if m.progress >= m.length {
			m.progress = 0
			m.moving = false
		}
	}
}
