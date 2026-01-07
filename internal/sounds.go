package internal

import (
	"path/filepath"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type SoundPlayer interface {
	PlayFx(name string)
}

type Sounds struct {
	sounds map[string]rl.Sound
}

func (s *Sounds) LoadSounds(soundDir string) {
	s.sounds = make(map[string]rl.Sound, 2)

	s.sounds["diamond_collision"] = rl.LoadSound(filepath.Join(soundDir, "effects", "diamond_collision.mp3"))
	s.sounds["diamond_collected"] = rl.LoadSound(filepath.Join(soundDir, "effects", "diamond_collected.mp3"))

	s.sounds["player_hurt"] = rl.LoadSound(filepath.Join(soundDir, "effects", "player_hurt.mp3"))

	s.sounds["boulder_collision"] = rl.LoadSound(filepath.Join(soundDir, "effects", "boulder_collision.mp3"))
	s.sounds["boulder_collision"] = rl.LoadSound(filepath.Join(soundDir, "effects", "boulder_collision.mp3"))

	s.sounds["stage_exit_opened"] = rl.LoadSound(filepath.Join(soundDir, "effects", "bang.mp3"))
	s.sounds["explosion"] = rl.LoadSound(filepath.Join(soundDir, "effects", "explosion.mp3"))
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
