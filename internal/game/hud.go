package game

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"galago/internal/config"
)

func (g *Game) drawHUD(screen *ebiten.Image) {
	hudFace := &text.GoTextFace{Source: g.images.UIFont, Size: 24}
	rightX := float64(config.ScreenWidth - 24)

	drawText(screen, fmt.Sprintf("%d %s", g.lives, lifeLabel(g.lives)), hudFace, rightX, 18, text.AlignEnd, color.White)
	drawText(screen, fmt.Sprintf("%d Pontos", g.score), hudFace, rightX, 50, text.AlignEnd, color.White)
	drawText(screen, fmt.Sprintf("x%02d", g.multiplier), hudFace, rightX, 82, text.AlignEnd, color.RGBA{R: 255, G: 220, B: 70, A: 255})

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

	if g.gameOver {
		titleFace := &text.GoTextFace{Source: g.images.UIFont, Size: 58}
		subtitleFace := &text.GoTextFace{Source: g.images.UIFont, Size: 24}
		centerX := float64(config.ScreenWidth) / 2

		drawText(screen, "GAME OVER", titleFace, centerX, 205, text.AlignCenter, color.RGBA{R: 255, G: 80, B: 80, A: 255})
		drawText(screen, "Pressione R para reiniciar", subtitleFace, centerX, 280, text.AlignCenter, color.White)
	}
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
