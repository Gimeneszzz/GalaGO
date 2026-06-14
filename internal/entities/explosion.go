package entities

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Explosion struct {
	X        float64
	Y        float64
	Radius   float64
	MaxTicks int
	Ticks    int
}

func NewExplosion(x, y, radius float64) Explosion {
	const ticks = 24
	return Explosion{
		X:        x,
		Y:        y,
		Radius:   radius,
		MaxTicks: ticks,
		Ticks:    ticks,
	}
}

func (explosion *Explosion) Update() {
	if explosion.Ticks > 0 {
		explosion.Ticks--
	}
}

func (explosion Explosion) Draw(screen *ebiten.Image) {
	if explosion.Ticks <= 0 {
		return
	}

	progress := 1 - float64(explosion.Ticks)/float64(explosion.MaxTicks)
	alpha := uint8(220 * explosion.Ticks / explosion.MaxTicks)
	outerRadius := float32(explosion.Radius * (0.35 + progress))
	innerRadius := float32(explosion.Radius * (0.16 + progress*0.35))

	vector.DrawFilledCircle(screen, float32(explosion.X), float32(explosion.Y), outerRadius, color.RGBA{R: 255, G: 95, B: 45, A: alpha / 2}, true)
	vector.DrawFilledCircle(screen, float32(explosion.X), float32(explosion.Y), innerRadius, color.RGBA{R: 255, G: 230, B: 90, A: alpha}, true)
}

func (explosion Explosion) IsDone() bool {
	return explosion.Ticks <= 0
}
