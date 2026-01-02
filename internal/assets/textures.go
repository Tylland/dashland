package assets

import (
	"path/filepath"
	"strings"
	"sync"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const assetbase string = "assets"
const textureDir string = "images"

var textures map[string]*rl.Texture2D = make(map[string]*rl.Texture2D)
var mu sync.Mutex = sync.Mutex{}

func LoadTexture(name string) *rl.Texture2D {
	path := filepath.Join(assetbase, textureDir, name+".png")

	return loadTexture(name, path)
}

func LoadTextureFromFile(path string) *rl.Texture2D {
	name := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))

	return loadTexture(name, path)
}

func loadTexture(name, path string) *rl.Texture2D {
	mu.Lock()
	if t, ok := textures[name]; ok {
		mu.Unlock()
		return t
	}
	mu.Unlock()

	tex := rl.LoadTexture(path)

	mu.Lock()
	textures[name] = &tex
	mu.Unlock()

	return &tex
}

func PreloadTextures(names ...string) {
	for _, n := range names {
		LoadTexture(n)
	}
}

func UnloadAll() {
	mu.Lock()
	defer mu.Unlock()
	for _, t := range textures {
		if t != nil {
			rl.UnloadTexture(*t)
		}
	}
	textures = make(map[string]*rl.Texture2D)
}
