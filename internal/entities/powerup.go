package entities

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"galago/internal/config"
)

type PowerUpKind int

const (
	PowerUpDoubleShot PowerUpKind = iota
	PowerUpShield
)

type PowerUp struct {
	X    float64
	Y    float64
	Kind PowerUpKind
}

func NewPowerUp(x, y float64, kind PowerUpKind) PowerUp {
	return PowerUp{X: x, Y: y, Kind: kind}
}

func (powerUp *PowerUp) Update() {
	powerUp.Y += config.PowerUpSpeed * config.DeltaTime
}

func (powerUp PowerUp) Draw(screen *ebiten.Image) {
	fill := color.RGBA{R: 80, G: 210, B: 255, A: 240}
	if powerUp.Kind == PowerUpShield {
		fill = color.RGBA{R: 100, G: 245, B: 145, A: 240}
	}

	vector.DrawFilledCircle(screen, float32(powerUp.X), float32(powerUp.Y), 13, fill, true)
	vector.StrokeLine(screen, float32(powerUp.X-7), float32(powerUp.Y), float32(powerUp.X+7), float32(powerUp.Y), 3, color.White, true)
	if powerUp.Kind == PowerUpShield {
		vector.StrokeLine(screen, float32(powerUp.X), float32(powerUp.Y-7), float32(powerUp.X), float32(powerUp.Y+7), 3, color.White, true)
	}
}

func (powerUp PowerUp) IsOffScreen() bool {
	return powerUp.Y-13 > config.ScreenHeight
}

func (powerUp PowerUp) CenterX() float64 {
	return powerUp.X
}

func (powerUp PowerUp) CenterY() float64 {
	return powerUp.Y
}

func (powerUp PowerUp) Radius() float64 {
	return 13
}
