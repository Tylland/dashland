package common

var (
	DirectionLeft      BlockVector = NewBlockVector(-1, 0)
	DirectionLeftUp    BlockVector = NewBlockVector(-1, -1)
	DirectionUp        BlockVector = NewBlockVector(0, -1)
	DirectionRightUp   BlockVector = NewBlockVector(1, -1)
	DirectionRight     BlockVector = NewBlockVector(1, 0)
	DirectionRightDown BlockVector = NewBlockVector(1, 1)
	DirectionDown      BlockVector = NewBlockVector(0, 1)
	DirectionLeftDown  BlockVector = NewBlockVector(-1, 1)
)
