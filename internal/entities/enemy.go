package entities

import (
	"github.com/hajimehoshi/ebiten/v2"

	"galago/internal/config"
)

type Enemy struct {
	X             float64
	Y             float64
	Width         float64
	Height        float64
	Speed         float64
	Angle         float64
	RotationSpeed float64
	Sprite        *ebiten.Image
}

func NewEnemy(x, y, speed, rotationSpeed float64, sprite *ebiten.Image) Enemy {
	width, height := sprite.Size()

	return Enemy{
		X:             x,
		Y:             y,
		Width:         float64(width),
		Height:        float64(height),
		Speed:         speed,
		RotationSpeed: rotationSpeed,
		Sprite:        sprite,
	}
}

func (e *Enemy) Update() {
	e.Y += e.Speed * config.DeltaTime
	e.Angle += e.RotationSpeed * config.DeltaTime
}

func (e Enemy) Draw(screen *ebiten.Image) {
	options := &ebiten.DrawImageOptions{}
	options.GeoM.Translate(-e.Width/2, -e.Height/2)
	options.GeoM.Rotate(e.Angle)
	options.GeoM.Translate(e.X+e.Width/2, e.Y+e.Height/2)
	screen.DrawImage(e.Sprite, options)
}

func (e Enemy) IsOffScreen() bool {
	return e.Y > config.ScreenHeight
}
