package game

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type SoundPlayer interface {
	PlayFx(name string)
}

type Sounds struct {
	sounds map[string]rl.Sound
}

func (s *Sounds) LoadSounds(dir string) {
	s.sounds = make(map[string]rl.Sound, 2)

	sound := rl.LoadSound("sounds/effects/diamond_collected.mp3")

	s.sounds["diamond_collected"] = sound

	sound = rl.LoadSound("sounds/effects/player_hurt.mp3")
	s.sounds["player_hurt"] = sound
	//rl.PlaySound(s.sounds[0])
}

func (s *Sounds) PlayFx(name string) {
	rl.PlaySound(s.sounds[name])
}

func (s *Sounds) UnloadSounds() {

	for _, sound := range s.sounds {
		rl.UnloadSound(sound)
	}

	s.sounds = nil
}
