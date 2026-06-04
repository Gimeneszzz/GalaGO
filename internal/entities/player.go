package entities

import (
	"github.com/hajimehoshi/ebiten/v2"

	"galago/internal/config"
)

// Player representa a nave controlada pelo usuario.
//
// X e Y sao coordenadas em pixels no espaco logico do jogo. A Ebitengine usa um
// sistema de coordenadas 2D em que (0,0) fica no canto superior esquerdo; X
// cresce para a direita e Y cresce para baixo.
type Player struct {
	X      float64
	Y      float64
	Width  float64
	Height float64
	Sprite *ebiten.Image
}

func NewPlayer(sprite *ebiten.Image) *Player {
	width, height := sprite.Size()

	return &Player{
		X:      float64(config.ScreenWidth-width) / 2,
		Y:      float64(config.ScreenHeight - height - 28),
		Width:  float64(width),
		Height: float64(height),
		Sprite: sprite,
	}
}

func (p *Player) Update() {
	// Verificar input do jogador e mover a nave. 
	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		p.X -= config.PlayerSpeed * config.DeltaTime
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		p.X += config.PlayerSpeed * config.DeltaTime
	}

	// Clamp: impede que a nave atravesse as bordas da tela.
	if p.X < 0 {
		p.X = 0
	}
	if p.X+p.Width > config.ScreenWidth {
		p.X = config.ScreenWidth - p.Width
	}
}

func (p *Player) Draw(screen *ebiten.Image) {
	options := &ebiten.DrawImageOptions{}
	options.GeoM.Translate(p.X, p.Y)
	screen.DrawImage(p.Sprite, options)
}

func (p *Player) CenterX() float64 {
	return p.X + p.Width/2
}
