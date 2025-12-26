package components

type Inventory struct {
	Items map[string]int
}

func NewInventory() *Inventory {
	return &Inventory{Items: make(map[string]int)}
}
