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
	X                    float64
	Y                    float64
	Width                float64
	Height               float64
	InvulnerabilityTicks int
	Sprite               *ebiten.Image
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
	if p.InvulnerabilityTicks > 0 {
		p.InvulnerabilityTicks--
	}

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
	// Durante a invulnerabilidade, alternamos frames visiveis e invisiveis.
	// Esse piscar comunica o estado de protecao sem precisar de outra sprite.
	if p.InvulnerabilityTicks > 0 && (p.InvulnerabilityTicks/5)%2 == 0 {
		return
	}

	options := &ebiten.DrawImageOptions{}
	options.GeoM.Translate(p.X, p.Y)
	screen.DrawImage(p.Sprite, options)
}

func (p *Player) CenterX() float64 {
	return p.X + p.Width/2
}

func (p *Player) CenterY() float64 {
	return p.Y + p.Height/2
}

func (p *Player) Radius() float64 {
	return min(p.Width, p.Height) * 0.36
}

func (p *Player) IsInvulnerable() bool {
	return p.InvulnerabilityTicks > 0
}

func (p *Player) TakeDamage() {
	p.InvulnerabilityTicks = config.InvulnerabilityTicks
}
