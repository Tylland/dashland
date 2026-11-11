package main

import (
	"math"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	W = 1280
	H = 720
)

type particle struct {
	x, y, vx, vy float32
	size, alpha  float32
	phase        float32
}

type tree struct{ x, h, w float32 }

type layer struct {
	parts      []particle
	trees      []tree
	baseY      float32 // ground line for the trees
	yMin, yMax float32 // band the mist lives in
	speed      float32 // parallax / wind multiplier
	treeColor  rl.Color
	mistTint   rl.Color
}

func f64(v float32) float64 { return float64(v) }

// Smooth pseudo-noise flow field -> angle in radians.
func flow(x, y, t float32) float32 {
	a := math.Sin(f64(x)*0.004+f64(t)*0.15)*1.2 + math.Cos(f64(y)*0.006-f64(t)*0.10)*1.2
	return float32(a)
}

func newLayer(baseY, treeH float32, count int, speed, sizeMin, sizeMax, alpha float32, tree, mist rl.Color) layer {
	l := layer{
		baseY:     baseY,
		yMin:      baseY - treeH*1.1,
		yMax:      baseY + 40,
		speed:     speed,
		treeColor: tree,
		mistTint:  mist,
	}

	// Pine silhouettes
	for x := float32(-30); x < W+30; x += 18 + rand.Float32()*34 {
		h := treeH * (0.45 + rand.Float32()*0.75)
		l.trees = append(l.trees, tree_(x, h, h*0.35))
	}

	// Mist particles
	for i := 0; i < count; i++ {
		l.parts = append(l.parts, particle{
			x:     rand.Float32() * W,
			y:     l.yMin + rand.Float32()*(l.yMax-l.yMin),
			size:  sizeMin + rand.Float32()*(sizeMax-sizeMin),
			alpha: alpha * (0.6 + rand.Float32()*0.6),
			phase: rand.Float32() * 6.28,
		})
	}
	return l
}

func tree_(x, h, w float32) tree { return tree{x, h, w} }

func (l *layer) update(dt, t, wind float32, mouse rl.Vector2) {
	const mouseR = 170
	for i := range l.parts {
		p := &l.parts[i]

		ang := flow(p.x, p.y, t)
		tx := float32(math.Cos(f64(ang)))*10*l.speed + wind*l.speed
		ty := float32(math.Sin(f64(ang))) * 3

		p.vx += (tx - p.vx) * dt * 0.6
		p.vy += (ty - p.vy) * dt * 0.6

		// Mist parts around the cursor
		dx, dy := p.x-mouse.X, p.y-mouse.Y
		d2 := dx*dx + dy*dy
		if d2 < mouseR*mouseR && d2 > 1 {
			d := float32(math.Sqrt(f64(d2)))
			f := (1 - d/mouseR) * 500 * dt
			p.vx += dx / d * f
			p.vy += dy / d * f
		}

		p.x += p.vx * dt
		p.y += p.vy * dt

		// Keep inside the band
		if p.y < l.yMin {
			p.vy += 25 * dt
		}
		if p.y > l.yMax {
			p.vy -= 25 * dt
		}

		// Wrap horizontally
		if p.x > W+p.size {
			p.x = -p.size
		} else if p.x < -p.size {
			p.x = W + p.size
		}
	}
}

func (l *layer) drawTrees() {
	rl.DrawRectangleRec(rl.NewRectangle(0, l.baseY, W, H-l.baseY), l.treeColor)
	for _, t := range l.trees {
		for k := 0; k < 3; k++ {
			fk := float32(k)
			top := l.baseY - t.h + fk*t.h*0.22
			bottom := top + t.h*0.45
			half := t.w * 0.5 * (0.55 + 0.22*fk)
			rl.DrawTriangle(
				rl.NewVector2(t.x, top),
				rl.NewVector2(t.x-half, bottom),
				rl.NewVector2(t.x+half, bottom),
				l.treeColor,
			)
		}
	}
}

func (l *layer) drawMist(tex rl.Texture2D, t float32) {
	src := rl.NewRectangle(0, 0, float32(tex.Width), float32(tex.Height))
	for _, p := range l.parts {
		a := p.alpha * (0.75 + 0.25*float32(math.Sin(f64(t)*0.5+f64(p.phase))))
		c := l.mistTint
		c.A = uint8(a * 255)
		dst := rl.NewRectangle(p.x, p.y, p.size, p.size*0.6)
		rl.DrawTexturePro(tex, src, dst, rl.NewVector2(p.size/2, p.size*0.3), 0, c)
	}
}

func main() {
	rl.SetConfigFlags(rl.FlagMsaa4xHint)
	rl.InitWindow(W, H, "Mist - raylib + Go")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	// Soft radial puff texture used by every mist particle
	img := rl.GenImageGradientRadial(256, 256, 0.0, rl.White, rl.Blank)
	puff := rl.LoadTextureFromImage(img)
	rl.UnloadImage(img)
	defer rl.UnloadTexture(puff)
	rl.SetTextureFilter(puff, rl.FilterBilinear)

	layers := []layer{
		newLayer(H*0.62, 130, 70, 0.5, 300, 560, 0.10, rl.NewColor(62, 78, 100, 255), rl.NewColor(190, 205, 225, 255)),
		newLayer(H*0.76, 190, 80, 0.8, 320, 600, 0.12, rl.NewColor(36, 50, 68, 255), rl.NewColor(175, 192, 215, 255)),
		newLayer(H*0.92, 270, 80, 1.3, 360, 680, 0.13, rl.NewColor(14, 22, 33, 255), rl.NewColor(160, 180, 205, 255)),
	}

	var t, wind float32
	skyTop := rl.NewColor(18, 26, 44, 255)
	skyBot := rl.NewColor(120, 140, 165, 255)

	for !rl.WindowShouldClose() {
		dt := rl.GetFrameTime()
		t += dt

		if rl.IsKeyDown(rl.KeyRight) {
			wind += 20 * dt
		}
		if rl.IsKeyDown(rl.KeyLeft) {
			wind -= 20 * dt
		}
		mouse := rl.GetMousePosition()

		for i := range layers {
			layers[i].update(dt, t, wind, mouse)
		}

		rl.BeginDrawing()
		rl.ClearBackground(skyTop)
		rl.DrawRectangleGradientV(0, 0, W, H, skyTop, skyBot)

		// Pale moon glow
		moon := rl.NewRectangle(W*0.72, H*0.05, 360, 360)
		rl.DrawTexturePro(puff, rl.NewRectangle(0, 0, 256, 256), moon, rl.NewVector2(180, 180), 0, rl.NewColor(235, 240, 255, 120))

		// Far -> near: trees, then a mist veil in front of them
		for i := range layers {
			layers[i].drawTrees()
			layers[i].drawMist(puff, t)
		}

		rl.DrawText("Move mouse: part the mist   |   Left/Right: wind", 16, 14, 18, rl.NewColor(255, 255, 255, 170))
		rl.EndDrawing()
	}
}
