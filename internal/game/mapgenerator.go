package game

import (
	"math/rand"

	"github.com/tylland/dashland/internal/common"
	"github.com/tylland/dashland/internal/ecs"
)

// CaveParams defines all parameters for procedural cave generation.
type CaveParams struct {
	Width              int
	Height             int
	BlockSize          int
	DiamondsRequired   int
	DiamondPoints      int
	DiamondBonusPoints int
	MagicWallDuration  float32

	// Object counts
	BoulderCount   int
	DiamondCount   int
	FireflyCount   int
	ButterflyCount int
	MagicWallCount int

	// Terrain: fraction of interior cells that become walls
	WallRatio float32
}

// EntityPlacement describes an entity to spawn during cave loading.
type EntityPlacement struct {
	Type     ecs.EntityType
	Position common.BlockPosition
}

// GeneratedCave holds the output of the procedural generator.
type GeneratedCave struct {
	Blocks         []BlockType
	Width          int
	Height         int
	Entities       []EntityPlacement
	PlayerPosition common.BlockPosition
	ExitPosition   common.BlockPosition
}

// DefaultCaveParams returns sensible defaults resembling classic Boulder Dash Cave A.
func DefaultCaveParams() CaveParams {
	return CaveParams{
		Width:              40,
		Height:             22,
		BlockSize:          32,
		DiamondsRequired:   12,
		DiamondPoints:      10,
		DiamondBonusPoints: 20,
		MagicWallDuration:  20.0,
		BoulderCount:       30,
		DiamondCount:       15,
		FireflyCount:       3,
		ButterflyCount:     0,
		MagicWallCount:     0,
		WallRatio:          0.10,
	}
}

// GenerateCave creates a random cave layout based on the given parameters.
//
// Algorithm:
//  1. Border cells = Bedrock, interior = Soil
//  2. Scatter Wall blocks in interior (WallRatio)
//  3. Clear a safe room around the player spawn (top-left)
//  4. Place exit door position (bottom-right)
//  5. Place boulders and diamonds on random soil cells (convert to Void)
//  6. Place magic walls on random soil cells
//  7. Clear small rooms for enemies and place them
func GenerateCave(params CaveParams) *GeneratedCave {
	w, h := params.Width, params.Height
	blocks := make([]BlockType, w*h)

	// --- Step 1: Bedrock border, Soil interior ---
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if x == 0 || x == w-1 || y == 0 || y == h-1 {
				blocks[y*w+x] = Bedrock
			} else {
				blocks[y*w+x] = Soil
			}
		}
	}

	// --- Step 2: Scatter walls in interior ---
	interior := collectInteriorPositions(w, h)
	rand.Shuffle(len(interior), func(i, j int) {
		interior[i], interior[j] = interior[j], interior[i]
	})

	wallCount := int(float32(len(interior)) * params.WallRatio)
	for i := 0; i < wallCount && i < len(interior); i++ {
		pos := interior[i]
		blocks[pos.Y*w+pos.X] = Wall
	}

	// --- Step 3: Player spawn (top-left safe room) ---
	playerPos := common.NewBlockPosition(2, 2)
	clearArea(blocks, w, h, playerPos, 2)

	// --- Step 4: Exit position (bottom-right) ---
	exitPos := common.NewBlockPosition(w-3, h-3)
	blocks[exitPos.Y*w+exitPos.X] = Void

	// --- Step 5: Collect available soil positions (skip near player & exit) ---
	available := collectAvailableSoil(blocks, w, h, playerPos, exitPos, 3)
	rand.Shuffle(len(available), func(i, j int) {
		available[i], available[j] = available[j], available[i]
	})

	var entities []EntityPlacement
	idx := 0

	// Place boulders
	for i := 0; i < params.BoulderCount && idx < len(available); i++ {
		pos := available[idx]
		idx++
		blocks[pos.Y*w+pos.X] = Void
		entities = append(entities, EntityPlacement{Type: EntityBoulder, Position: pos})
	}

	// Place diamonds
	for i := 0; i < params.DiamondCount && idx < len(available); i++ {
		pos := available[idx]
		idx++
		blocks[pos.Y*w+pos.X] = Void
		entities = append(entities, EntityPlacement{Type: EntityDiamond, Position: pos})
	}

	// Place magic walls
	for i := 0; i < params.MagicWallCount && idx < len(available); i++ {
		pos := available[idx]
		idx++
		blocks[pos.Y*w+pos.X] = Void
		entities = append(entities, EntityPlacement{Type: EntityMagicWall, Position: pos})
	}

	// --- Step 6: Place enemies in cleared rooms ---
	enemySoil := collectAvailableSoil(blocks, w, h, playerPos, exitPos, 5)
	rand.Shuffle(len(enemySoil), func(i, j int) {
		enemySoil[i], enemySoil[j] = enemySoil[j], enemySoil[i]
	})

	enemyIdx := 0

	for i := 0; i < params.FireflyCount && enemyIdx < len(enemySoil); i++ {
		pos := enemySoil[enemyIdx]
		enemyIdx++
		clearArea(blocks, w, h, pos, 2)
		entities = append(entities, EntityPlacement{Type: EntityFirefly, Position: pos})
	}

	for i := 0; i < params.ButterflyCount && enemyIdx < len(enemySoil); i++ {
		pos := enemySoil[enemyIdx]
		enemyIdx++
		clearArea(blocks, w, h, pos, 2)
		entities = append(entities, EntityPlacement{Type: EntityButterfly, Position: pos})
	}

	return &GeneratedCave{
		Blocks:         blocks,
		Width:          w,
		Height:         h,
		Entities:       entities,
		PlayerPosition: playerPos,
		ExitPosition:   exitPos,
	}
}

// collectInteriorPositions returns all positions that are at least 2 cells from the border.
func collectInteriorPositions(w, h int) []common.BlockPosition {
	var positions []common.BlockPosition
	for y := 2; y < h-2; y++ {
		for x := 2; x < w-2; x++ {
			positions = append(positions, common.NewBlockPosition(x, y))
		}
	}
	return positions
}

// collectAvailableSoil returns soil cells that are at least minDist cells away
// from each excluded position.
func collectAvailableSoil(blocks []BlockType, w, h int, exclude1, exclude2 common.BlockPosition, minDist int) []common.BlockPosition {
	var positions []common.BlockPosition
	for y := 2; y < h-2; y++ {
		for x := 2; x < w-2; x++ {
			if blocks[y*w+x] != Soil {
				continue
			}
			if intAbs(x-exclude1.X) < minDist && intAbs(y-exclude1.Y) < minDist {
				continue
			}
			if intAbs(x-exclude2.X) < minDist && intAbs(y-exclude2.Y) < minDist {
				continue
			}
			positions = append(positions, common.NewBlockPosition(x, y))
		}
	}
	return positions
}

// clearArea sets a square region to Void, clamped to interior (avoids bedrock border).
func clearArea(blocks []BlockType, w, h int, center common.BlockPosition, radius int) {
	for dy := -radius; dy <= radius; dy++ {
		for dx := -radius; dx <= radius; dx++ {
			x, y := center.X+dx, center.Y+dy
			if x > 0 && x < w-1 && y > 0 && y < h-1 {
				blocks[y*w+x] = Void
			}
		}
	}
}

func intAbs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
