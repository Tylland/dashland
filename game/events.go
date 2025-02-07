package game

type GameEvent interface {
	IsEvent()
}

type GameEventListner interface {
	OnEvent(event GameEvent)
}

type GameEventDispatcher struct {
	listeners []GameEventListner
}

func (d *GameEventDispatcher) AddListener(listener GameEventListner) {
	d.listeners = append(d.listeners, listener)
}
