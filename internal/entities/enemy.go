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
	Health        int
	MaxHealth     int
	Reward        int
	EscapePenalty int
	DamageTicks   int
	Sprite        *ebiten.Image
}

func NewEnemy(x, y, speed, rotationSpeed float64, sprite *ebiten.Image) Enemy {
	width, height := sprite.Size()
	health, reward, escapePenalty := meteorStats(width, height)

	return Enemy{
		X:             x,
		Y:             y,
		Width:         float64(width),
		Height:        float64(height),
		Speed:         speed,
		RotationSpeed: rotationSpeed,
		Health:        health,
		MaxHealth:     health,
		Reward:        reward,
		EscapePenalty: escapePenalty,
		Sprite:        sprite,
	}
}

func (e *Enemy) Update() {
	e.Y += e.Speed * config.DeltaTime
	e.Angle += e.RotationSpeed * config.DeltaTime
	if e.DamageTicks > 0 {
		e.DamageTicks--
	}
}

func (e Enemy) Draw(screen *ebiten.Image) {
	options := &ebiten.DrawImageOptions{}
	options.GeoM.Translate(-e.Width/2, -e.Height/2)
	options.GeoM.Rotate(e.Angle)
	options.GeoM.Translate(e.X+e.Width/2, e.Y+e.Height/2)
	if e.DamageTicks > 0 {
		options.ColorScale.Scale(1, 0.35, 0.35, 1)
	}
	screen.DrawImage(e.Sprite, options)
}

func (e Enemy) IsOffScreen() bool {
	return e.Y > config.ScreenHeight
}

func (e Enemy) CenterX() float64 {
	return e.X + e.Width/2
}

func (e Enemy) CenterY() float64 {
	return e.Y + e.Height/2
}

func (e Enemy) Radius() float64 {
	return min(e.Width, e.Height) * 0.42
}

func (e *Enemy) TakeDamage() bool {
	e.Health--
	e.DamageTicks = config.MeteorDamageFlashTicks
	return e.Health <= 0
}

func meteorStats(width, height int) (health, reward, escapePenalty int) {
	size := max(width, height)
	switch {
	case size <= 35:
		return 1, config.SmallMeteorReward, config.SmallMeteorEscapePenalty
	case size <= 60:
		return 2, config.MediumMeteorReward, config.MediumMeteorEscapePenalty
	default:
		return 3, config.LargeMeteorReward, config.LargeMeteorEscapePenalty
	}
}
