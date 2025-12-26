package components

type Health struct {
	Points int
}

func NewHealth(points int) *Health {
	return &Health{Points: points}
}
