package entities

import (
	"github.com/hajimehoshi/ebiten/v2"

	"galago/internal/config"
)

type Bullet struct {
	X      float64
	Y      float64
	Width  float64
	Height float64
	Sprite *ebiten.Image
}

func NewBullet(x, y float64, sprite *ebiten.Image) Bullet {
	width, height := sprite.Size()

	return Bullet{
		X:      x - float64(width)/2,
		Y:      y,
		Width:  float64(width),
		Height: float64(height),
		Sprite: sprite,
	}
}

func (b *Bullet) Update() {
	// Como Y cresce para baixo na tela, subtrair de Y faz o tiro subir.
	b.Y -= config.BulletSpeed * config.DeltaTime
}

func (b Bullet) Draw(screen *ebiten.Image) {
	options := &ebiten.DrawImageOptions{}
	options.GeoM.Translate(b.X, b.Y)
	screen.DrawImage(b.Sprite, options)
}

func (b Bullet) IsOffScreen() bool {
	return b.Y+b.Height < 0
}
