package utils

import (
	"github.com/lafriks/go-tiled"
)

func Find(tiles []*tiled.LayerTile, match func(*tiled.LayerTile) bool) (found []*tiled.LayerTile) {
	for _, tile := range tiles {
		if match(tile) {
			found = append(found, tile)
		}
	}

	return
}
