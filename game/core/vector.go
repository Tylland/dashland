package core

type Vector struct {
	X float32
	Y float32
}

func NewVector(x float32, y float32) Vector {
	return Vector{X: x, Y: y}
}

func (v *Vector) IsZero() bool {
	return v.X == 0 && v.Y == 0
}

func (v *Vector) Clear() {
	v.X = 0
	v.Y = 0
}

func (v Vector) Offset(deltaX float32, deltaY float32) Vector {
	return Vector{X: v.X + deltaX, Y: v.Y + deltaY}
}

func (v Vector) Add(vector Vector) Vector {
	return Vector{X: v.X + vector.X, Y: v.Y + vector.Y}
}

func (v Vector) Subtract(vector Vector) Vector {
	return Vector{X: v.X - vector.X, Y: v.Y - vector.Y}
}
