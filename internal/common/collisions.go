package common

type CollisionLayer uint16

type Collider struct {
	Layer CollisionLayer
	Mask  CollisionLayer
}

func NewCollider(layer, mask CollisionLayer) *Collider {
	return &Collider{Layer: layer, Mask: mask}
}

func (c *Collider) CollidesWith(other *Collider) bool {
	return c.Mask&other.Layer != 0
}

// type Position struct {
// 	position rl.Vector2
// 	listners []func(rl.Vector2)
// }

// func (p *Position) Update(pos rl.Vector2) {
// 	p.position = pos
// 	p.notyfyUpdated()
// }

// func (p *Position) Register(listner func(rl.Vector2)) {
// 	p.listners = append(p.listners, listner)
// }

// func (p *Position) Unregister(listner func(rl.Vector2)) {
// 	//TODO:	p.listners = remove(p.listners, listner)
// }

// func (p *Position) notyfyUpdated() {

// 	for _, listner := range p.listners {
// 		listner(p.position)
// 	}
// }

// type Body interface {
// 	Position() rl.Vector2
// 	SetPosition(pos rl.Vector2)
// 	IsColliding(body Body) bool
// }

// type Collider interface {
// 	Body() Body
// 	//IsColliding(collider Collider) bool
// }

// type Box struct {
// 	rl.Rectangle
// }

// func NewBox(x float32, y float32, width float32, height float32) Box {
// 	return Box{Rectangle: rl.Rectangle{X: x, Y: y, Width: width, Height: height}}
// }

// func (tb *Box) Overlaps(ob *Box) bool {
// 	return tb.X > ob.X+ob.Width && tb.X+tb.Width > ob.X && tb.Y > ob.Y+ob.Height && tb.Y+tb.Height < ob.Y
// }

// func (b *Box) Position() rl.Vector2 {
// 	return rl.Vector2{X: b.X, Y: b.Y}
// }

// func (b *Box) SetPosition(pos rl.Vector2) {
// 	b.X = pos.X - b.Width/2
// 	b.Y = pos.Y - b.Height/2
// }

// func (b *Box) IsColliding(body Body) bool {
// 	box, ok := body.(*Box)

// 	if ok {
// 		return b.Overlaps(box)
// 	}

// 	return false
// }

// type BoxBody struct {
// 	Width  float32
// 	Height float32
// }

// func NewBoxBody(width, height float32) *BoxBody {
// 	return &BoxBody{
// 		Width:  width,
// 		Height: height,
// 	}
// }

// func (b *BoxBody) Rectangle(pos rl.Vector2) rl.Rectangle {
// 	return rl.Rectangle{
// 		X:      pos.X,
// 		Y:      pos.Y,
// 		Width:  b.Width,
// 		Height: b.Height,
// 	}
// }

// func (b *BoxBody) Overlaps(pos1 rl.Vector2, pos2 rl.Vector2, other *BoxBody) bool {
// 	return pos1.X < pos2.X+other.Width &&
// 		pos1.X+b.Width > pos2.X &&
// 		pos1.Y < pos2.Y+other.Height &&
// 		pos1.Y+b.Height > pos2.Y
// }
