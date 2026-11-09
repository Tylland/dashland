package systems

import (
	"github.com/tylland/dashland/internal/ecs"
	"github.com/tylland/dashland/internal/game"
)

const timeoutDeathDelay = 10.0

type StageTimerSystem struct {
	stage               *game.Stage
	sound               game.SoundPlayer
	expired             bool
	deathDelayRemaining float32
	deathIssued         bool
	beepTime            float32
}

func NewStageTimerSystem(stage *game.Stage, sound game.SoundPlayer) *StageTimerSystem {
	return &StageTimerSystem{stage: stage, sound: sound}
}

func (s *StageTimerSystem) Update(world *ecs.World, deltaTime float32) {
	if s.stage.TimeLimit <= 0 || s.deathIssued {
		return
	}

	if !s.expired {
		s.stage.TimeRemaining -= deltaTime

		if s.stage.TimeRemaining < 10 {
			if s.beepTime == 0.0 || (s.beepTime-s.stage.TimeRemaining) >= 1.0 {
				s.sound.PlayFx("countdown_beep")
				s.beepTime = s.stage.TimeRemaining
			}

		}

		if s.stage.TimeRemaining > 0 {
			return
		}

		s.stage.TimeRemaining = 0
		s.stage.TimeExpired = true
		s.expired = true
		s.deathDelayRemaining = timeoutDeathDelay
		world.AddEvent("stageexpired", game.NewStageExpiredEvent())
		world.SetPaused(true)

		s.sound.PlayFx("player_died")

		return
	}

	s.deathDelayRemaining -= deltaTime
	if s.deathDelayRemaining <= 0 {
		s.deathIssued = true
		world.AddEvent("timeoutdeath", game.NewTimeoutDeathEvent())
	}
}

func (s *StageTimerSystem) RunWhenPaused() bool {
	return true
}
