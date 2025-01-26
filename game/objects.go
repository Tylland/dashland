package game

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/lafriks/go-tiled"
)

type GroundObject interface {
	isObstacle() bool
	IsObstacleForPlayer(player *Player) bool
	GetBlockType() BlockType
	GetBlockPosition() BlockPosition
	SetBlockPosition(position BlockPosition)
	HasBehavior(behavior BlockBehavior) bool
	update(deltaTime float32)
	render()
}

type FallingObject interface {
	GroundObject
	StartFalling(startPos rl.Vector2, endPos rl.Vector2)
	IsFalling() bool
	UpdateFalling(float32)
}

type CollectableObject interface {
	Collected()
}

type PushableObject interface {
	GroundObject
	Pushed(player *Player, position BlockPosition)
	// StartPushing(startPos rl.Vector2, endPos rl.Vector2)
	// IsPushing() bool
	// UpdatePushing(float32)
}

type CollectableBehavior struct {
	collected bool
}

func (cb *CollectableBehavior) Collected() {
	cb.collected = true
}

// type PushableBehavior struct {
// 	pushing bool
// }

// func (p *PushableBehavior) StartPushing(startPos rl.Vector2, endPos rl.Vector2){
// 	p.pushing = true

// }

// func (p *PushableBehavior) IsFalling() bool
// UpdateFalling(float32)

func NewGroundObject(world *world, blockType BlockType, x int, y int) (GroundObject, error) {
	blockPosition := BlockPosition{X: x, Y: y}
	//body := BoxBody{position: world.GetPosition(blockPosition), Width: world.blockWidth, Height: world.blockHeight}

	switch blockType {
	// case Bedrock:
	// 	return &BedrockObject{BlockObject: BlockObject{world: world, blockType: blockType, blockPosition: BlockPosition{X: x, Y: y}, behavior: Obstacle}}, nil
	// case Void:
	// 	return &VoidObject{BlockObject: BlockObject{world: world, blockType: blockType, blockPosition: BlockPosition{X: x, Y: y}}}, nil
	// case Soil:
	// 	return &SoilObject{BlockObject: BlockObject{world: world, blockType: blockType, blockPosition: BlockPosition{X: x, Y: y}}}, nil
	case Boulder:
		return NewBoulderObject(world, blockPosition), nil
		// &BoulderObject{GravityObject: GravityObject{
		// 	BlockObject: BlockObject{world: world, blockType: blockType, blockPosition: blockPosition, behavior: CanFall | Obstacle | Pushable},
		// 	falling:     MovementTimer{ProgressTimer: ProgressTimer{}}}}, nil
	case Diamond:
		return NewDiamondObject(world, blockPosition), nil
		//  {GravityObject: GravityObject{
		// 	BlockObject: BlockObject{world: world, blockType: blockType, blockPosition: BlockPosition{X: x, Y: y}, behavior: },
		// 	falling:     MovementTimer{ProgressTimer: ProgressTimer{}}}}, nil
	default:
		return nil, fmt.Errorf("%v (%d) is unknown blocktype", blockType, int(blockType))
	}
}

type GroundMap struct {
	MapSize
	objectTextures rl.Texture2D
	groundCorners  rl.Texture2D
	objects        []GroundObject
}

func (gm *GroundMap) InitObjects(world *world, tiles []*tiled.LayerTile) {
	gm.objects = make([]GroundObject, len(tiles))

	for index, tile := range tiles {
		object, err := NewGroundObject(world, BlockType(tile.ID), index%world.width, index/world.width)

		if err != nil {

		}

		gm.objects[index] = object
	}
}

func (gm *GroundMap) GetObject(position BlockPosition) GroundObject {

	if position.X < 0 || position.X >= gm.width || position.Y < 0 || position.Y >= gm.height {
		return nil
	}

	return gm.objects[position.Y*gm.width+position.X]
}

func (gm *GroundMap) CheckObjectAtPosition(blockType BlockType, position BlockPosition) bool {
	return gm.objects[position.Y*gm.width+position.X].GetBlockType() == blockType
}

func (gm *GroundMap) MoveObject(source GroundObject, targetPos BlockPosition) {
	sourcePosition := source.GetBlockPosition()

	gm.objects[targetPos.Y*gm.width+targetPos.X] = source
	gm.objects[sourcePosition.Y*gm.width+sourcePosition.X] = nil

	source.SetBlockPosition(targetPos)
}

func (gm *GroundMap) RemoveObject(doomed GroundObject) {
	sourcePosition := doomed.GetBlockPosition()

	gm.objects[sourcePosition.Y*gm.width+sourcePosition.X] = nil
}

func (gm *GroundMap) SwapObject(source GroundObject, targetPos BlockPosition) {

	target := gm.GetObject(targetPos)

	if target == nil {
		return
	}

	sourcePos := source.GetBlockPosition()

	gm.objects[targetPos.Y*gm.width+targetPos.X], gm.objects[sourcePos.Y*gm.width+sourcePos.X] = gm.objects[sourcePos.Y*gm.width+sourcePos.X], gm.objects[targetPos.Y*gm.width+targetPos.X]

	target.SetBlockPosition(sourcePos)
	source.SetBlockPosition(targetPos)
}

type GravityObject struct {
	BlockObject
	BoxBody
	falling MovementTimer
}

func NewGravityObject(world *world, blockPosition BlockPosition, blockType BlockType, behavior BlockBehavior) GravityObject {
	return GravityObject{
		BlockObject: BlockObject{world: world, blockType: blockType, blockPosition: blockPosition, behavior: behavior},
		BoxBody:     BoxBody{position: world.GetPosition(blockPosition), Width: world.blockWidth, Height: world.blockHeight},
		falling:     MovementTimer{ProgressTimer: ProgressTimer{}}}
}

func (g *GravityObject) StartFalling(startPos rl.Vector2, endPos rl.Vector2) {
	const gravitySpeed float32 = 0.4

	vector := rl.Vector2Subtract(endPos, g.position)

	g.falling.StartMovment(startPos, vector, gravitySpeed, nil)
}

func (g *GravityObject) IsFalling() bool {
	return g.falling.running
}

func (g *GravityObject) UpdateFalling(deltaTime float32) {
	if g.falling.running {
		g.falling.UpdateTimer(deltaTime)

		g.position = g.falling.Position()
	}
}

type BlockObject struct {
	world         *world
	blockType     BlockType
	blockPosition BlockPosition
	behavior      BlockBehavior
}

func (bo *BlockObject) isObstacle() bool {
	return bo.HasBehavior(Obstacle)
}

func (bo *BlockObject) IsObstacleForPlayer(player *Player) bool {
	return bo.HasBehavior(Obstacle)
}

func (bo *BlockObject) GetBlockType() BlockType {
	return bo.blockType
}

func (bo *BlockObject) GetBlockPosition() BlockPosition {
	return bo.blockPosition
}

func (bo *BlockObject) SetBlockPosition(position BlockPosition) {
	bo.blockPosition = position
}

func (bo *BlockObject) HasBehavior(b BlockBehavior) bool {
	return bo.behavior&b == b
}

func (bo *BlockObject) update(deltaTime float32) {

}

func (bo *BlockObject) render() {
	bm := bo.world

	rl.DrawTextureRec(bm.objectTextures, rl.NewRectangle(float32(bo.blockType)*bm.blockWidth, 0, bm.blockWidth, bm.blockHeight), rl.NewVector2(float32(bo.blockPosition.X)*bm.blockWidth, float32(bo.blockPosition.Y)*bm.blockHeight), rl.White)
}

type BoulderObject struct {
	GravityObject
	pushing MovementTimer
}

func NewBoulderObject(world *world, blockPosition BlockPosition) *BoulderObject {
	return &BoulderObject{
		GravityObject: NewGravityObject(world, blockPosition, Boulder, CanFall|Obstacle|Pushable),
		pushing:       MovementTimer{ProgressTimer: ProgressTimer{}},
	}
}

func (bo *BoulderObject) IsObstacleForPlayer(player *Player) bool {

	// if bo.falling. {
	// 	return true
	// }

	offset := bo.blockPosition.Subtract(player.blockPosition)

	if offset.Y != 0 {
		return true
	}

	return bo.world.checkPositionOccupied(bo.blockPosition.Add(offset))
}

func (bo *BoulderObject) Pushed(player *Player, position BlockPosition) {
	bo.pushing.StartMovment(bo.world.GetPosition(bo.blockPosition), bo.world.GetPosition(position), 64, func() { bo.world.MoveObject(bo, position) })
}

func (bo *BoulderObject) Body() Body {
	return &bo.BoxBody
}

func (bo *BoulderObject) update(deltaTime float32) {
	fmt.Println("Upadte BoulderObject")

	bo.world.ApplyGravity(bo, deltaTime)
}

func (bo *BoulderObject) render() {
	bm := bo.world

	rl.DrawTextureRec(bm.objectTextures, rl.NewRectangle(float32(bo.blockType)*bm.blockWidth, 0, bm.blockWidth, bm.blockHeight), bo.position, rl.White)
}

type DiamondObject struct {
	GravityObject
	CollectableBehavior
}

func NewDiamondObject(world *world, blockPosition BlockPosition) *DiamondObject {
	return &DiamondObject{
		GravityObject: NewGravityObject(world, blockPosition, Diamond, CanFall|Collectable),
	}
}

func (bo *DiamondObject) Collected() {
	bo.world.PlayFx("diamond_collected")
	bo.world.RemoveObject(bo)
}

func (bo *DiamondObject) StartFalling(startPos rl.Vector2, endPos rl.Vector2) {
	bo.GravityObject.StartFalling(startPos, endPos)

	//	bo.world.PlayFx("diamond_collected")
}

func (bo *DiamondObject) update(deltaTime float32) {

	bo.world.ApplyGravity(bo, deltaTime)
}

func (bo *DiamondObject) render() {
	bm := bo.world

	if !bo.collected {
		rl.DrawTextureRec(bm.objectTextures, rl.NewRectangle(float32(bo.blockType)*bm.blockWidth, 0, bm.blockWidth, bm.blockHeight), bo.position, rl.White)
	}
}
