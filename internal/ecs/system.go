package ecs

type System interface {
	Update(world *World, deltaTime float32)
}

type SoundPlayer interface {
	PlayFx(name string)
}
