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
	IsColliding(body Body) bool
}

type Collider interface {
	Body
}

type Box struct {
	rl.Rectangle
}

func (b1 *Box) Overlaps(b2 *Box) bool {
	return b1.X > b2.X+b2.Width && b1.X+b1.Width > b2.X && b1.Y > b2.Y+b2.Height && b1.Y+b1.Height < b2.Y
}

func (b *Box) IsColliding(body Body) bool {
	box := body.(*Box)

	if box != nil {
		return b.Overlaps(box)
	}

	return false
}
