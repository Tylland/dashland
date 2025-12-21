package components

// type CollisionComponent struct {
// 	common.BoxBody
// 	IsColliding bool
// 	Layer       CollisionLayer
// 	Mask        CollisionLayer
// }

// const (
// 	LayerNone CollisionLayer = 0
// 	LayerAll  CollisionLayer = 0xFF
// 	// Common game layers
// 	LayerPlayer CollisionLayer = 1 << iota
// 	LayerEnemy
// 	LayerProjectile
// 	LayerWall
// 	LayerItem
// )

// // NewCollisionComponent creates a new collision component with default values
// func NewCollisionComponent(width, height float32, layer CollisionLayer) *CollisionComponent {
// 	return &CollisionComponent{
// 		BoxBody: common.BoxBody{
// 			Width:  width,
// 			Height: height,
// 		},
// 		Layer:       layer,
// 		Mask:        LayerAll,
// 		IsColliding: false,
// 	}
// }

// // CanCollideWith checks if this component can collide with the specified layer
// func (c *CollisionComponent) CanCollideWith(other CollisionLayer) bool {
// 	return (c.Mask & other) != 0
// }

// // Reset resets the collision state
// func (c *CollisionComponent) Reset() {
// 	c.IsColliding = false
// }
