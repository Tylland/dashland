package components

type Inventory struct {
	Points   int
	Diamonds int
	Items    map[string]int
}

func NewInventory() *Inventory {
	return &Inventory{Items: make(map[string]int)}
}
