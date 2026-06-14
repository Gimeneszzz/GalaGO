package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type floatingText struct {
	X        float64
	Y        float64
	Value    string
	Color    color.RGBA
	MaxTicks int
	Ticks    int
}

func newFloatingText(x, y float64, value string, clr color.RGBA) floatingText {
	const ticks = 42
	return floatingText{
		X:        x,
		Y:        y,
		Value:    value,
		Color:    clr,
		MaxTicks: ticks,
		Ticks:    ticks,
	}
}

func (textItem *floatingText) Update() {
	if textItem.Ticks > 0 {
		textItem.Ticks--
		textItem.Y -= 0.55
	}
}

func (textItem floatingText) Draw(screen *ebiten.Image, source *text.GoTextFaceSource) {
	if textItem.Ticks <= 0 {
		return
	}

	clr := textItem.Color
	clr.A = uint8(int(clr.A) * textItem.Ticks / textItem.MaxTicks)
	face := &text.GoTextFace{Source: source, Size: 18}
	drawText(screen, textItem.Value, face, textItem.X, textItem.Y, text.AlignCenter, clr)
}

func (textItem floatingText) IsDone() bool {
	return textItem.Ticks <= 0
}
