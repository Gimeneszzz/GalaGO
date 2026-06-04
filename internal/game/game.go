package game

import (
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"galago/internal/assets"
	"galago/internal/config"
	"galago/internal/entities"
)

// Game eh o estado principal do jogo.
//
// A Ebitengine espera que esse tipo implemente tres metodos:
// - Update: atualiza regras, fisica simples e inputs.
// - Draw: desenha o frame atual.
// - Layout: informa a resolucao logica usada pelo jogo.
type Game struct {
	images     assets.Images
	player     *entities.Player
	bullets    []entities.Bullet
	enemies    []entities.Enemy
	background []entities.BackgroundObject

	rng            *rand.Rand
	lastEnemySpawn time.Time
}

func New() (*Game, error) {
	images, warnings, err := assets.LoadImages()
	if err != nil {
		return nil, err
	}

	if formatted := assets.FormatWarnings(warnings); formatted != "" {
		log.Println(formatted)
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	return &Game{
		images:         images,
		player:         entities.NewPlayer(images.Player),
		bullets:        make([]entities.Bullet, 0, 16),
		enemies:        make([]entities.Enemy, 0, 16),
		background:     createBackground(images, rng),
		rng:            rng,
		lastEnemySpawn: time.Now(),
	}, nil
}

func (g *Game) Update() error {
	g.player.Update()
	g.handleShooting()
	g.spawnEnemies()
	g.updateBackground()
	g.updateBullets()
	g.updateEnemies()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	screen.Fill(color.RGBA{R: 8, G: 10, B: 20, A: 255})

	for _, object := range g.background {
		object.Draw(screen)
	}
	for _, bullet := range g.bullets {
		bullet.Draw(screen)
	}
	for _, enemy := range g.enemies {
		enemy.Draw(screen)
	}
	g.player.Draw(screen)

	
	ebitenutil.DebugPrint(screen, "Mover: A/D ou Setas | Atirar: Espaco")
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	
	return config.ScreenWidth, config.ScreenHeight
}

func (g *Game) handleShooting() {
	// IsKeyPressed seria verdadeiro enquanto a tecla estivesse segurada, criando
	// tiros em todos os frames. IsKeyJustPressed dispara somente no frame em que
	// a tecla acabou de ser pressionada, que e o comportamento esperado aqui.
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.bullets = append(g.bullets, entities.NewBullet(
			g.player.CenterX(),
			g.player.Y-12,
			g.images.Bullet,
		))
	}
}

func (g *Game) spawnEnemies() {
	now := time.Now()
	if now.Sub(g.lastEnemySpawn) < time.Duration(config.EnemySpawnMs)*time.Millisecond {
		return
	}
	g.lastEnemySpawn = now

	sprite := g.images.Meteors[g.rng.Intn(len(g.images.Meteors))]
	enemyWidth, enemyHeight := sprite.Size()
	maxX := config.ScreenWidth - enemyWidth
	x := float64(g.rng.Intn(maxX + 1))
	y := -float64(enemyHeight)
	speed := config.EnemySpeed + g.rng.Float64()*85
	rotationSpeed := (g.rng.Float64()*2 - 1) * 1.4

	g.enemies = append(g.enemies, entities.NewEnemy(x, y, speed, rotationSpeed, sprite))
}

func (g *Game) updateBackground() {
	for i := range g.background {
		g.background[i].Update()
	}
}

func (g *Game) updateBullets() {
	alive := g.bullets[:0]
	for i := range g.bullets {
		g.bullets[i].Update()
		if !g.bullets[i].IsOffScreen() {
			alive = append(alive, g.bullets[i])
		}
	}
	g.bullets = alive
}

func (g *Game) updateEnemies() {
	alive := g.enemies[:0]
	for i := range g.enemies {
		g.enemies[i].Update()
		if !g.enemies[i].IsOffScreen() {
			alive = append(alive, g.enemies[i])
		}
	}
	g.enemies = alive
}

func createBackground(images assets.Images, rng *rand.Rand) []entities.BackgroundObject {
	objects := make([]entities.BackgroundObject, 0, 10)

	// Estrelas pequenas se movem devagar para dar sensacao de profundidade 
	for i := 0; i < 8 && len(images.Stars) > 0; i++ {
		sprite := images.Stars[rng.Intn(len(images.Stars))]
		width, height := sprite.Size()
		scale := 0.35 + rng.Float64()*0.35
		x := rng.Float64() * float64(config.ScreenWidth-int(float64(width)*scale))
		y := rng.Float64()*float64(config.ScreenHeight+height) - float64(height)
		speed := 8 + rng.Float64()*18

		objects = append(objects, entities.NewBackgroundObject(x, y, scale, speed, sprite))
	}

	// Planetas sao sprites grandes, entao usamos escala baixa. Eles ficam no
	// fundo como decoracao, nao como obstaculos.
	for i := 0; i < 2 && len(images.Planets) > 0; i++ {
		sprite := images.Planets[rng.Intn(len(images.Planets))]
		width, height := sprite.Size()
		scale := 0.18 + rng.Float64()*0.12
		x := rng.Float64() * float64(config.ScreenWidth-int(float64(width)*scale))
		y := rng.Float64()*float64(config.ScreenHeight+height) - float64(height)
		speed := 3 + rng.Float64()*6

		objects = append(objects, entities.NewBackgroundObject(x, y, scale, speed, sprite))
	}

	return objects
}
