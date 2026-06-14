package entities

import (
	"github.com/hajimehoshi/ebiten/v2"

	"galago/internal/config"
)

type ProfessorShip struct {
	X           float64
	Y           float64
	Width       float64
	Height      float64
	Health      int
	MaxHealth   int
	Reward      int
	DamageTicks int
	Sprite      *ebiten.Image
}

func NewProfessorShip(sprite *ebiten.Image, player *Player) *ProfessorShip {
	width, height := sprite.Size()
	ship := &ProfessorShip{
		Width:     float64(width),
		Height:    float64(height),
		Health:    config.ProfessorShipHealth,
		MaxHealth: config.ProfessorShipHealth,
		Reward:    config.ProfessorShipReward,
		Sprite:    sprite,
	}
	ship.X = player.CenterX() - ship.Width/2
	ship.Y = professorShipTargetY(player, ship.Height)
	return ship
}

func (ship *ProfessorShip) Update(player *Player) {
	targetX := player.CenterX() - ship.Width/2
	step := config.ProfessorShipFollowSpeed * config.DeltaTime
	if ship.X < targetX {
		ship.X = min(ship.X+step, targetX)
	}
	if ship.X > targetX {
		ship.X = max(ship.X-step, targetX)
	}

	ship.Y = professorShipTargetY(player, ship.Height)

	if ship.X < 0 {
		ship.X = 0
	}
	if ship.X+ship.Width > config.ScreenWidth {
		ship.X = config.ScreenWidth - ship.Width
	}
	if ship.DamageTicks > 0 {
		ship.DamageTicks--
	}
}

func (ship ProfessorShip) Draw(screen *ebiten.Image) {
	options := &ebiten.DrawImageOptions{}
	options.GeoM.Translate(ship.X, ship.Y)
	if ship.DamageTicks > 0 {
		options.ColorScale.Scale(1, 0.35, 0.35, 1)
	}
	screen.DrawImage(ship.Sprite, options)
}

func (ship ProfessorShip) CenterX() float64 {
	return ship.X + ship.Width/2
}

func (ship ProfessorShip) CenterY() float64 {
	return ship.Y + ship.Height/2
}

func (ship ProfessorShip) Radius() float64 {
	return min(ship.Width, ship.Height) * 0.42
}

func (ship *ProfessorShip) TakeDamage() bool {
	ship.Health--
	ship.DamageTicks = config.MeteorDamageFlashTicks
	return ship.Health <= 0
}

func professorShipTargetY(player *Player, height float64) float64 {
	targetY := player.Y - config.ProfessorShipVerticalDistance - height
	if targetY < 36 {
		return 36
	}
	return targetY
}
