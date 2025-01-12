package game

import (
	"fmt"
)

type actor interface {
	update(deltaTime float32)
	render()
}

type world struct {
	blockMap *BlockMap
	player   *Player
	actors   []actor
}

func NewWorld(player *Player) *world {
	actors := []actor{player}

	return &world{player: player, actors: actors}
}

func (w *world) update(deltaTime float32) {
	w.blockMap.update(deltaTime)

	for _, act := range w.actors {
		act.update(deltaTime)
	}
}

func (w *world) render() {
	w.blockMap.render()

	for _, act := range w.actors {
		act.render()
	}
}

func (w *world) addActor(actor actor) {
	w.actors = append(w.actors, actor)
}

func (w *world) obstacleForPlayer(player *Player, position BlockPosition) bool {
	block, success := w.blockMap.GetBlock(position.X, position.Y)

	if !success {
		return true
	}

	return block.ObstacleForPlayer(player)
}

func (w *world) VisitBlock(position BlockPosition) {

	fmt.Printf("Block at position %d,%d changed type from %d", position.X, position.Y, w.blockMap.blocks[position.Y*w.blockMap.width+position.X].blockType)
	w.blockMap.blocks[position.Y*w.blockMap.width+position.X] = NewBlock(w, Void, position.X, position.Y)
	fmt.Printf(" to %d \n", w.blockMap.blocks[position.Y*w.blockMap.width+position.X].blockType)
}

func (w *world) checkPositionOccupied(position BlockPosition) bool {
	return !w.blockMap.CheckTypeAtPosition(Void, position) || w.player.blockPosition.IsSame(position)
}

func (w *world) checkPlayerAtPosition(position BlockPosition) bool {
	return w.player.blockPosition.IsSame(position)
}
