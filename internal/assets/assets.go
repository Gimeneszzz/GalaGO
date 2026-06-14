package assets

import (
	"bytes"
	"fmt"
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	rawassets "galago/assets"
)

// Guardar os sprites em uma unica estrutura facilita trocar a origem dos assets
type Images struct {
	Player        *ebiten.Image
	Bullet        *ebiten.Image
	Meteors       []*ebiten.Image
	Stars         []*ebiten.Image
	Planets       []*ebiten.Image
	ProfessorShip *ebiten.Image
	UIFont        *text.GoTextFaceSource
}

// LoadImages organiza os sprites carregados pelo pacote galago/assets.

func LoadImages() (Images, []string, error) {
	var warnings []string

	uiFont, err := text.NewGoTextFaceSource(bytes.NewReader(rawassets.UIFontData))
	if err != nil {
		return Images{}, nil, fmt.Errorf("nao foi possivel carregar fontui.ttf: %w", err)
	}

	meteors := rawassets.MeteorSprites
	if len(meteors) == 0 {
		warnings = append(warnings, "sprites de meteoros nao encontrados; usando placeholder")
		meteors = []*ebiten.Image{
			placeholder(44, 44, color.RGBA{R: 170, G: 120, B: 80, A: 255}),
			placeholder(74, 58, color.RGBA{R: 120, G: 120, B: 130, A: 255}),
		}
	}

	return Images{
		Player:        rawassets.PlayerSprite,
		Bullet:        rawassets.LaserSprite,
		Meteors:       meteors,
		Stars:         rawassets.StarsSprites,
		Planets:       rawassets.PlanetsSprites,
		ProfessorShip: rawassets.ProfessorShipSprite,
		UIFont:        uiFont,
	}, warnings, nil
}

func placeholder(width, height int, fill color.Color) *ebiten.Image {
	img := ebiten.NewImage(width, height)
	img.Fill(fill)
	return img
}

// FormatWarnings transforma os avisos em uma string para log

func FormatWarnings(warnings []string) string {
	if len(warnings) == 0 {
		return ""
	}

	return strings.Join(warnings, "; ")
}
