package common

// Animation describes a named animation sequence inside a sprite sheet.
type Animation struct {
	BaseX         float32
	BaseY         float32
	FrameCount    uint
	FrameDuration float32
	Loop          bool
}
