package components

type PushableComponent struct {
	Mass float32
}

func NewPushableComponent(mass float32) *PushableComponent {
	return &PushableComponent{Mass: mass}
}
