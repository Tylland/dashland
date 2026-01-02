package components

type Lifetime struct {
	Remaining float32
}

func NewLifetime(remaining float32) *Lifetime {
	return &Lifetime{Remaining: remaining}
}
