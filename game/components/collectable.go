package components

type CollectableType uint32

const (
	CollectableNone    CollectableType = 0
	CollectableDiamond CollectableType = 1
)

type CollectableComponent struct {
	Type      CollectableType
	Amount    int
	Collected bool
}

func NewCollectableComponent(t CollectableType, amount int) *CollectableComponent {
	return &CollectableComponent{Type: t, Amount: amount, Collected: false}
}
