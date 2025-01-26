package game

import rl "github.com/gen2brain/raylib-go/raylib"

type Position struct {
	position rl.Vector2
	listners []func(rl.Vector2)
}

func (p *Position) Update(pos rl.Vector2) {
	p.position = pos
	p.notyfyUpdated()
}

func (p *Position) Register(listner func(rl.Vector2)) {
	p.listners = append(p.listners, listner)
}

func (p *Position) Unregister(listner func(rl.Vector2)) {
	//TODO:	p.listners = remove(p.listners, listner)
}

func (p *Position) notyfyUpdated() {

	for _, listner := range p.listners {
		listner(p.position)
	}
}

type Body interface {
	Position() rl.Vector2
	SetPosition(pos rl.Vector2)
	IsColliding(body Body) bool
}

type Collider interface {
	Body() Body
	//IsColliding(collider Collider) bool
}

type Box struct {
	rl.Rectangle
}

func NewBox(x float32, y float32, width float32, height float32) Box {
	return Box{Rectangle: rl.Rectangle{X: x, Y: y, Width: width, Height: height}}
}

func (tb *Box) Overlaps(ob *Box) bool {
	return tb.X > ob.X+ob.Width && tb.X+tb.Width > ob.X && tb.Y > ob.Y+ob.Height && tb.Y+tb.Height < ob.Y
}

func (b *Box) Position() rl.Vector2 {
	return rl.Vector2{X: b.X, Y: b.Y}
}

func (b *Box) SetPosition(pos rl.Vector2) {
	b.X = pos.X - b.Width/2
	b.Y = pos.Y - b.Height/2
}

func (b *Box) IsColliding(body Body) bool {
	box, ok := body.(*Box)

	if ok {
		return b.Overlaps(box)
	}

	return false
}

type BoxBody struct {
	position rl.Vector2
	Width    float32
	Height   float32
}

func (b *BoxBody) Position() rl.Vector2 {
	return b.position
}

func (b *BoxBody) SetPosition(pos rl.Vector2) {
	b.position.X = pos.X
	b.position.Y = pos.Y
}

func (b *BoxBody) Rectangle() rl.Rectangle {
	return rl.Rectangle{X: b.position.X, Y: b.position.Y, Width: b.Width, Height: b.Height}
}

func (b *BoxBody) IsColliding(body Body) bool {
	box, ok := body.(*BoxBody)

	if ok {
		return b.Overlaps(box)
	}

	return false
}

func (tb *BoxBody) Overlaps(ob *BoxBody) bool {
	return tb.position.X > ob.position.X+ob.Width && tb.position.X+tb.Width > ob.position.X && tb.position.Y > ob.position.Y+ob.Height && tb.position.Y+tb.Height < ob.position.Y
}
