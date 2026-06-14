package game

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"galago/internal/config"
)

func (g *Game) drawProfessorCharge(screen *ebiten.Image) {
	angle := clampProfessorShotAngle(g.professorAimAngle)
	originX := g.professorShip.CenterX()
	originY := g.professorShip.Y + g.professorShip.Height
	alpha := uint8(70 + 120*chargeProgress(g.professorChargeTicks))

	for segment := 0; segment < 9; segment++ {
		startDistance := float64(segment*70 + 18)
		endDistance := startDistance + 35
		startX := originX + math.Cos(angle)*startDistance
		startY := originY + math.Sin(angle)*startDistance
		endX := originX + math.Cos(angle)*endDistance
		endY := originY + math.Sin(angle)*endDistance

		vector.StrokeLine(
			screen,
			float32(startX),
			float32(startY),
			float32(endX),
			float32(endY),
			2,
			color.RGBA{R: 255, G: 70, B: 80, A: alpha},
			true,
		)
	}
	vector.DrawFilledCircle(screen, float32(originX), float32(originY), 8, color.RGBA{R: 255, G: 230, B: 90, A: alpha}, true)
}

func (g *Game) drawPlayerShield(screen *ebiten.Image) {
	if g.shieldTicks <= 0 {
		return
	}

	pulse := 0.75 + math.Sin(float64(g.shieldTicks)*0.22)*0.12
	radius := float32(g.player.Radius() * 2.1 * pulse)
	vector.DrawFilledCircle(
		screen,
		float32(g.player.CenterX()),
		float32(g.player.CenterY()),
		radius,
		color.RGBA{R: 95, G: 230, B: 150, A: 48},
		true,
	)
	vector.StrokeLine(
		screen,
		float32(g.player.CenterX()-float64(radius)),
		float32(g.player.CenterY()),
		float32(g.player.CenterX()+float64(radius)),
		float32(g.player.CenterY()),
		2,
		color.RGBA{R: 120, G: 255, B: 180, A: 150},
		true,
	)
}

func (g *Game) drawProfessorEntryShield(screen *ebiten.Image) {
	pulse := 0.82 + math.Sin(float64(g.professorSpawnTicks+g.professorActiveTicks)*0.18)*0.08
	radius := float32(g.professorShip.Radius() * 1.35 * pulse)
	vector.DrawFilledCircle(
		screen,
		float32(g.professorShip.CenterX()),
		float32(g.professorShip.CenterY()),
		radius,
		color.RGBA{R: 80, G: 190, B: 255, A: 42},
		true,
	)
	vector.DrawFilledCircle(
		screen,
		float32(g.professorShip.CenterX()),
		float32(g.professorShip.CenterY()),
		radius*0.72,
		color.RGBA{R: 80, G: 190, B: 255, A: 24},
		true,
	)
}

func (g *Game) drawScreenFlash(screen *ebiten.Image) {
	if g.screenFlashTicks <= 0 {
		return
	}

	alpha := uint8(90 * g.screenFlashTicks / 14)
	vector.DrawFilledRect(screen, 0, 0, config.ScreenWidth, config.ScreenHeight, color.RGBA{R: 255, G: 45, B: 45, A: alpha}, false)
}

func chargeProgress(ticks int) float64 {
	if ticks <= 0 {
		return 1
	}
	return 1 - float64(ticks)/float64(config.ProfessorShipChargeTicks)
}
