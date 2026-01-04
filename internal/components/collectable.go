package components

type CollectableType uint32

const (
	CollectableNone    CollectableType = 0
	CollectableDiamond CollectableType = 1
)

type CollectableComponent struct {
	Name      string
	Amount    int
	Collected bool
}

func NewCollectableComponent(name string, amount int) *CollectableComponent {
	return &CollectableComponent{Name: name, Amount: amount, Collected: false}
}
