package game

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"galago/internal/config"
)

func (g *Game) drawHUD(screen *ebiten.Image) {
	hudFace := &text.GoTextFace{Source: g.images.UIFont, Size: 24}
	rightX := float64(config.ScreenWidth - 24)

	drawText(screen, fmt.Sprintf("%d %s", g.lives, lifeLabel(g.lives)), hudFace, rightX, 18, text.AlignEnd, color.White)
	drawText(screen, fmt.Sprintf("%d Pontos", g.score), hudFace, rightX, 50, text.AlignEnd, color.White)
	drawText(screen, fmt.Sprintf("x%02d", g.multiplier), hudFace, rightX, 82, text.AlignEnd, color.RGBA{R: 255, G: 220, B: 70, A: 255})
	drawText(screen, fmt.Sprintf("Recorde %d", g.highScore), hudFace, rightX, 114, text.AlignEnd, color.RGBA{R: 210, G: 225, B: 240, A: 255})
	drawText(screen, fmt.Sprintf("Berri %d/%d Nv %d", g.professorDefeats, config.ProfessorDefeatsToWin, g.professorVisitLevel()+1), hudFace, 24, 18, text.AlignStart, color.White)
	drawText(screen, fmt.Sprintf("Som %s (M)", soundState(g.muted)), hudFace, 24, 50, text.AlignStart, color.RGBA{R: 210, G: 225, B: 240, A: 255})
	g.drawPowerUpStatus(screen)

	if g.professorShip != nil {
		g.drawProfessorHealth(screen)
	} else if g.shouldShowProfessorWarning() {
		warningFace := &text.GoTextFace{Source: g.images.UIFont, Size: 24}
		drawText(screen, fmt.Sprintf("Berri Nv %d se aproxima...", g.professorVisitLevel()+1), warningFace, float64(config.ScreenWidth)/2, 112, text.AlignCenter, color.RGBA{R: 255, G: 215, B: 80, A: 255})
		g.drawProfessorSpawnProgress(screen)
	} else if g.professorShip == nil && g.professorSpawnTicks > 0 && !g.gameOver && !g.won {
		g.drawProfessorSpawnProgress(screen)
	}

	if g.penaltyTicks > 0 {
		alpha := uint8(255 * g.penaltyTicks / config.PenaltyMessageTicks)
		penaltyFace := &text.GoTextFace{Source: g.images.UIFont, Size: 30}
		drawText(
			screen,
			fmt.Sprintf("-%d Pontos", g.penaltyValue),
			penaltyFace,
			float64(config.ScreenWidth)/2,
			float64(config.ScreenHeight-54),
			text.AlignCenter,
			color.RGBA{R: 255, G: 65, B: 65, A: alpha},
		)
	}

	if g.won {
		titleFace := &text.GoTextFace{Source: g.images.UIFont, Size: 48}
		subtitleFace := &text.GoTextFace{Source: g.images.UIFont, Size: 24}
		centerX := float64(config.ScreenWidth) / 2

		drawText(screen, "Você venceu Berri!", titleFace, centerX, 185, text.AlignCenter, color.RGBA{R: 80, G: 220, B: 140, A: 255})
		drawText(screen, fmt.Sprintf("Pontuacao final: %d", g.score), subtitleFace, centerX, 255, text.AlignCenter, color.White)
		drawText(screen, fmt.Sprintf("Recorde: %d   Rank %s", g.highScore, finalRank(g.score)), subtitleFace, centerX, 292, text.AlignCenter, color.RGBA{R: 255, G: 220, B: 70, A: 255})
		drawText(screen, "Pressione R para jogar de novo", subtitleFace, centerX, 342, text.AlignCenter, color.White)
	} else if g.gameOver {
		titleFace := &text.GoTextFace{Source: g.images.UIFont, Size: 38}
		subtitleFace := &text.GoTextFace{Source: g.images.UIFont, Size: 24}
		centerX := float64(config.ScreenWidth) / 2

		drawText(screen, "Você foi derrotado por Berri", titleFace, centerX, 205, text.AlignCenter, color.RGBA{R: 255, G: 80, B: 80, A: 255})
		drawText(screen, "Pressione R para reiniciar", subtitleFace, centerX, 280, text.AlignCenter, color.White)
	} else if g.rewardMessageTicks > 0 {
		rewardFace := &text.GoTextFace{Source: g.images.UIFont, Size: 30}
		alpha := uint8(255 * g.rewardMessageTicks / config.ProfessorRewardMessageTicks)
		drawText(screen, "Berri repelida!", rewardFace, float64(config.ScreenWidth)/2, 170, text.AlignCenter, color.RGBA{R: 80, G: 220, B: 140, A: alpha})
	}
}

func (g *Game) drawProfessorHealth(screen *ebiten.Image) {
	const (
		barWidth  = 260
		barHeight = 12
		x         = 270
		y         = 24
	)

	fraction := float32(g.professorShip.Health) / float32(g.professorShip.MaxHealth)
	vector.DrawFilledRect(screen, x, y, barWidth, barHeight, color.RGBA{R: 50, G: 40, B: 55, A: 220}, false)
	vector.DrawFilledRect(screen, x, y, barWidth*fraction, barHeight, color.RGBA{R: 255, G: 80, B: 120, A: 255}, false)
	vector.StrokeLine(screen, x, y, x+barWidth, y, 2, color.White, false)
	vector.StrokeLine(screen, x, y+barHeight, x+barWidth, y+barHeight, 2, color.White, false)
}

func (g *Game) drawProfessorSpawnProgress(screen *ebiten.Image) {
	if g.professorSpawnTotalTicks <= 0 {
		return
	}

	const (
		barWidth  = 220
		barHeight = 8
		x         = 290
		y         = 145
	)

	progress := 1 - float32(g.professorSpawnTicks)/float32(g.professorSpawnTotalTicks)
	progress = min(max(progress, 0), 1)
	vector.DrawFilledRect(screen, x, y, barWidth, barHeight, color.RGBA{R: 50, G: 40, B: 55, A: 180}, false)
	vector.DrawFilledRect(screen, x, y, barWidth*progress, barHeight, color.RGBA{R: 255, G: 215, B: 80, A: 230}, false)
}

func (g *Game) drawPowerUpStatus(screen *ebiten.Image) {
	statusFace := &text.GoTextFace{Source: g.images.UIFont, Size: 18}
	y := 82.0
	if g.doubleShotTicks > 0 {
		drawText(screen, fmt.Sprintf("Tiro duplo %.0fs", float64(g.doubleShotTicks)*config.DeltaTime), statusFace, 24, y, text.AlignStart, color.RGBA{R: 80, G: 210, B: 255, A: 255})
		y += 24
	}
	if g.shieldTicks > 0 {
		drawText(screen, fmt.Sprintf("Escudo %.0fs", float64(g.shieldTicks)*config.DeltaTime), statusFace, 24, y, text.AlignStart, color.RGBA{R: 100, G: 245, B: 145, A: 255})
	}
}

func (g *Game) shouldShowProfessorWarning() bool {
	return g.images.ProfessorShip != nil &&
		g.professorShip == nil &&
		g.professorSpawnTicks > 0 &&
		g.professorSpawnTicks <= config.ProfessorShipWarningTicks &&
		!g.gameOver &&
		!g.won
}

func (g *Game) drawStartScreen(screen *ebiten.Image) {
	titleFace := &text.GoTextFace{Source: g.images.UIFont, Size: 54}
	subtitleFace := &text.GoTextFace{Source: g.images.UIFont, Size: 24}
	controlsFace := &text.GoTextFace{Source: g.images.UIFont, Size: 18}
	centerX := float64(config.ScreenWidth) / 2

	drawText(screen, "GalaGO", titleFace, centerX, 178, text.AlignCenter, color.RGBA{R: 80, G: 210, B: 255, A: 255})
	drawText(screen, "Pressione Enter para iniciar o jogo", subtitleFace, centerX, 250, text.AlignCenter, color.White)
	drawText(screen, "Setas ou A/D: mover para os lados", controlsFace, centerX, 318, text.AlignCenter, color.RGBA{R: 210, G: 225, B: 240, A: 255})
	drawText(screen, "Espaco: atirar", controlsFace, centerX, 348, text.AlignCenter, color.RGBA{R: 210, G: 225, B: 240, A: 255})
	drawText(screen, "M: ligar ou desligar som", controlsFace, centerX, 378, text.AlignCenter, color.RGBA{R: 210, G: 225, B: 240, A: 255})
	drawText(screen, "P: pausar", controlsFace, centerX, 408, text.AlignCenter, color.RGBA{R: 210, G: 225, B: 240, A: 255})

	width, _ := g.images.Player.Size()
	options := &ebiten.DrawImageOptions{}
	options.GeoM.Translate(centerX-float64(width)/2, 450)
	screen.DrawImage(g.images.Player, options)
}

func (g *Game) drawPauseOverlay(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, 0, 0, config.ScreenWidth, config.ScreenHeight, color.RGBA{R: 5, G: 8, B: 18, A: 165}, false)
	titleFace := &text.GoTextFace{Source: g.images.UIFont, Size: 44}
	subtitleFace := &text.GoTextFace{Source: g.images.UIFont, Size: 22}
	centerX := float64(config.ScreenWidth) / 2

	drawText(screen, "PAUSADO", titleFace, centerX, 235, text.AlignCenter, color.RGBA{R: 255, G: 220, B: 70, A: 255})
	drawText(screen, "Pressione P para continuar", subtitleFace, centerX, 295, text.AlignCenter, color.White)
}

func drawText(screen *ebiten.Image, value string, face text.Face, x, y float64, align text.Align, clr color.Color) {
	options := &text.DrawOptions{}
	options.PrimaryAlign = align
	options.GeoM.Translate(x, y)
	options.ColorScale.ScaleWithColor(clr)
	text.Draw(screen, value, face, options)
}

func lifeLabel(lives int) string {
	if lives == 1 {
		return "vida"
	}
	return "vidas"
}

func soundState(muted bool) string {
	if muted {
		return "OFF"
	}
	return "ON"
}

func finalRank(score int) string {
	switch {
	case score >= 9000:
		return "S"
	case score >= 6000:
		return "A"
	default:
		return "B"
	}
}
