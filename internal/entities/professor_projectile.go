package entities

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"

	"galago/internal/config"
)

type ProfessorProjectile struct {
	X      float64
	Y      float64
	VX     float64
	VY     float64
	Width  float64
	Height float64
	Angle  float64
	Sprite *ebiten.Image
}

func NewProfessorProjectile(x, y, velocityX, velocityY float64, sprite *ebiten.Image) ProfessorProjectile {
	width, height := sprite.Size()

	return ProfessorProjectile{
		X:      x - float64(width)/2,
		Y:      y - float64(height)/2,
		VX:     velocityX,
		VY:     velocityY,
		Width:  float64(width),
		Height: float64(height),
		Angle:  math.Atan2(velocityY, velocityX) + math.Pi/2,
		Sprite: sprite,
	}
}

func (projectile *ProfessorProjectile) Update() {
	projectile.X += projectile.VX * config.DeltaTime
	projectile.Y += projectile.VY * config.DeltaTime
}

func (projectile ProfessorProjectile) Draw(screen *ebiten.Image) {
	options := &ebiten.DrawImageOptions{}
	options.GeoM.Translate(-projectile.Width/2, -projectile.Height/2)
	options.GeoM.Rotate(projectile.Angle)
	options.GeoM.Translate(projectile.CenterX(), projectile.CenterY())
	options.ColorScale.Scale(1, 0.45, 0.45, 1)
	screen.DrawImage(projectile.Sprite, options)
}

func (projectile ProfessorProjectile) IsOffScreen() bool {
	return projectile.X+projectile.Width < 0 ||
		projectile.X > config.ScreenWidth ||
		projectile.Y+projectile.Height < 0 ||
		projectile.Y > config.ScreenHeight
}

func (projectile ProfessorProjectile) CenterX() float64 {
	return projectile.X + projectile.Width/2
}

func (projectile ProfessorProjectile) CenterY() float64 {
	return projectile.Y + projectile.Height/2
}

func (projectile ProfessorProjectile) Radius() float64 {
	return min(projectile.Width, projectile.Height) * 0.5
}
