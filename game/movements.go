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
	finished func()
}

func (m *Movement) Start(startPos rl.Vector2, endPos rl.Vector2, speed float32, finished func()) {
	m.startPos = startPos
	m.endPos = endPos
	m.vector = rl.Vector2Subtract(endPos, startPos)
	m.length = rl.Vector2Length(m.vector)
	m.speed = speed
	m.progress = 0.0
	m.moving = true
	m.finished = finished
}

func (m *Movement) Update(deltaTime float32) {

	if m.moving {
		m.progress += m.speed * deltaTime

		progressVector := rl.Vector2Scale(m.vector, rl.Clamp(m.progress, 0, m.length)/m.length)

		m.position = rl.Vector2Add(m.startPos, progressVector)

		if m.progress >= m.length {
			m.progress = 0
			m.moving = false

			if m.finished != nil {
				m.finished()
				m.finished = nil
			}
		}
	}
}

type ProgressTimer struct {
	duration    float32
	elapsedTime float32
	progress    float32
	running     bool
	onFinished  func()
}

func (pt *ProgressTimer) StartTimer(duration float32, onFinished func()) {
	pt.duration = duration
	pt.elapsedTime = 0.0
	pt.progress = 0.0
	pt.running = true
	pt.onFinished = onFinished
}

func (pt *ProgressTimer) ResetTimer() {
	pt.duration = 0
	pt.elapsedTime = 0.0
	pt.progress = 0.0
	pt.running = false
	pt.onFinished = nil
}

func (pt *ProgressTimer) UpdateTimer(deltaTime float32) {

	if !pt.running {
		return
	}

	pt.elapsedTime += deltaTime
	pt.progress = pt.elapsedTime / pt.duration

	if pt.progress >= 1.0 {
		pt.progress = 1

		if pt.onFinished != nil {
			pt.onFinished()
		}

		pt.running = false
	}
}

type MovementTimer struct {
	ProgressTimer
	startPos rl.Vector2
	vector   rl.Vector2
}

func (mt *MovementTimer) StartMovement(startPos rl.Vector2, vector rl.Vector2, duration float32, onFinished func()) {
	mt.startPos = startPos
	mt.vector = vector

	mt.StartTimer(duration, onFinished)
}

func (mt *MovementTimer) Position() rl.Vector2 {
	progressVector := rl.Vector2Scale(mt.vector, mt.progress)

	return rl.Vector2Add(mt.startPos, progressVector)
}
