package game

import (
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
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
	shootCooldown  int

	score      int
	lives      int
	multiplier int
	gameOver   bool

	penaltyValue int
	penaltyTicks int
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

	g := &Game{
		images:     images,
		background: createBackground(images, rng),
		rng:        rng,
	}
	g.reset()
	return g, nil
}

func (g *Game) Update() error {
	// Retornar ebiten.Termination encerra o jogo de forma limpa, sem fazer o
	// RunGame tratar a saida solicitada pelo jogador como um erro.
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		return ebiten.Termination
	}

	g.updateBackground()
	g.updatePenaltyMessage()

	if g.gameOver {
		if inpututil.IsKeyJustPressed(ebiten.KeyR) {
			g.reset()
		}
		return nil
	}

	g.player.Update()
	g.handleShooting()
	g.spawnEnemies()
	g.updateBullets()
	g.updateEnemies()
	g.handleCollisions()
	g.removeEscapedEnemies()

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
	g.drawHUD(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return config.ScreenWidth, config.ScreenHeight
}

func (g *Game) handleShooting() {
	if g.shootCooldown > 0 {
		g.shootCooldown--
	}

	// O pequeno cooldown permite segurar Espaco sem criar 60 tiros por segundo.
	if ebiten.IsKeyPressed(ebiten.KeySpace) && g.shootCooldown == 0 {
		g.bullets = append(g.bullets, entities.NewBullet(
			g.player.CenterX(),
			g.player.Y-12,
			g.images.Bullet,
		))
		g.shootCooldown = config.ShootCooldownTicks
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
	minX, maxX := meteorSpawnBounds(enemyWidth)
	x := float64(minX + g.rng.Intn(maxX-minX+1))
	y := -float64(enemyHeight)
	speed := config.EnemySpeed + g.rng.Float64()*85
	rotationSpeed := (g.rng.Float64()*2 - 1) * 1.4

	g.enemies = append(g.enemies, entities.NewEnemy(x, y, speed, rotationSpeed, sprite))
}

func meteorSpawnBounds(meteorWidth int) (minX, maxX int) {
	minX = config.SpawnSideMargin
	maxX = config.ScreenWidth - config.SpawnSideMargin - meteorWidth

	// Esta protecao mantem o calculo valido caso uma sprite futura seja maior
	// do que toda a area disponivel entre as margens.
	if maxX < minX {
		maxX = minX
	}
	return minX, maxX
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
	for i := range g.enemies {
		g.enemies[i].Update()
	}
}

func (g *Game) removeEscapedEnemies() {
	remaining := g.enemies[:0]
	for _, enemy := range g.enemies {
		if enemy.IsOffScreen() {
			g.score -= enemy.EscapePenalty
			g.multiplier = 1
			g.showPenalty(enemy.EscapePenalty)
			continue
		}
		remaining = append(remaining, enemy)
	}
	g.enemies = remaining
}

func (g *Game) updatePenaltyMessage() {
	if g.penaltyTicks > 0 {
		g.penaltyTicks--
	}
}

func (g *Game) showPenalty(value int) {
	g.penaltyValue = value
	g.penaltyTicks = config.PenaltyMessageTicks
}

func (g *Game) reset() {
	g.player = entities.NewPlayer(g.images.Player)
	g.bullets = make([]entities.Bullet, 0, 24)
	g.enemies = make([]entities.Enemy, 0, 24)
	g.lastEnemySpawn = time.Now()
	g.shootCooldown = 0
	g.score = 0
	g.lives = config.InitialLives
	g.multiplier = 1
	g.gameOver = false
	g.penaltyValue = 0
	g.penaltyTicks = 0
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
