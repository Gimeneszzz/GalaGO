package entities

import (
	"github.com/hajimehoshi/ebiten/v2"

	"galago/internal/config"
)

// BackgroundObject representa um elemento visual do cenario, como estrelas e
// planetas. Ele nao participa da logica de colisao; serve para dar movimento e
// profundidade ao espaco sem misturar essa responsabilidade com Player/Enemy.
type BackgroundObject struct {
	X      float64
	Y      float64
	Scale  float64
	Speed  float64
	Sprite *ebiten.Image
}

func NewBackgroundObject(x, y, scale, speed float64, sprite *ebiten.Image) BackgroundObject {
	return BackgroundObject{
		X:      x,
		Y:      y,
		Scale:  scale,
		Speed:  speed,
		Sprite: sprite,
	}
}

func (o *BackgroundObject) Update() {
	o.Y += o.Speed * config.DeltaTime

	_, height := o.Sprite.Size()
	if o.Y > config.ScreenHeight+float64(height)*o.Scale {
		o.Y = -float64(height) * o.Scale
	}
}

func (o BackgroundObject) Draw(screen *ebiten.Image) {
	options := &ebiten.DrawImageOptions{}
	options.GeoM.Scale(o.Scale, o.Scale)
	options.GeoM.Translate(o.X, o.Y)
	screen.DrawImage(o.Sprite, options)
}
