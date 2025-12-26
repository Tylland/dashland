package components

type Damage struct {
	Points int
}

func NewDamage(points int) *Damage {
	return &Damage{Points: points}
}
