package game

import (
	"image/color"
	"log"
	"math"
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
	images               assets.Images
	sounds               *soundPlayer
	player               *entities.Player
	bullets              []entities.Bullet
	enemies              []entities.Enemy
	professorShip        *entities.ProfessorShip
	professorProjectiles []entities.ProfessorProjectile
	powerUps             []entities.PowerUp
	explosions           []entities.Explosion
	floatingTexts        []floatingText
	background           []entities.BackgroundObject

	rng                      *rand.Rand
	lastEnemySpawn           time.Time
	shootCooldown            int
	professorShootCooldown   int
	professorSpawnTicks      int
	professorActiveTicks     int
	professorChargeTicks     int
	professorAimAngle        float64
	professorSpawnTotalTicks int
	pauseTicks               int
	screenFlashTicks         int

	score            int
	highScore        int
	lives            int
	multiplier       int
	started          bool
	paused           bool
	gameOver         bool
	won              bool
	muted            bool
	professorDefeats int
	professorVisits  int
	doubleShotTicks  int
	shieldTicks      int

	penaltyValue       int
	penaltyTicks       int
	rewardMessageTicks int
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
		sounds:     newSoundPlayer(),
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
	if inpututil.IsKeyJustPressed(ebiten.KeyM) {
		g.muted = !g.muted
	}
	g.updateMusic()

	g.updateBackground()
	g.updatePenaltyMessage()
	g.updateRewardMessage()
	g.updateExplosions()
	g.updateFloatingTexts()
	g.updateScreenFlash()

	if !g.started {
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			g.reset()
			g.started = true
		}
		return nil
	}

	if g.gameOver || g.won {
		if inpututil.IsKeyJustPressed(ebiten.KeyR) {
			g.reset()
		}
		if g.won {
			g.updateWinScreenMotion()
		}
		return nil
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		g.paused = !g.paused
	}
	if g.paused {
		return nil
	}

	if g.pauseTicks > 0 {
		g.pauseTicks--
		return nil
	}

	g.player.Update()
	g.updatePowerUpTimers()
	g.handleShooting()
	g.updateProfessorShip()
	g.handleProfessorShooting()
	g.spawnEnemies()
	g.updateBullets()
	g.updateProfessorProjectiles()
	g.updatePowerUps()
	g.updateEnemies()
	g.handleCollisions()
	g.removeEscapedEnemies()
	g.updateHighScore()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 8, G: 10, B: 20, A: 255})

	for _, object := range g.background {
		object.Draw(screen)
	}

	if !g.started {
		g.drawStartScreen(screen)
		return
	}

	for _, bullet := range g.bullets {
		bullet.Draw(screen)
	}
	for _, enemy := range g.enemies {
		enemy.Draw(screen)
	}
	if g.professorShip != nil {
		g.professorShip.Draw(screen)
		if !g.professorShip.IsReady(g.player) {
			g.drawProfessorEntryShield(screen)
		}
	}
	for _, projectile := range g.professorProjectiles {
		projectile.Draw(screen)
	}
	for _, powerUp := range g.powerUps {
		powerUp.Draw(screen)
	}
	if g.professorShip != nil && g.professorChargeTicks > 0 {
		g.drawProfessorCharge(screen)
	}
	for _, explosion := range g.explosions {
		explosion.Draw(screen)
	}
	for _, floatingText := range g.floatingTexts {
		floatingText.Draw(screen, g.images.UIFont)
	}
	g.player.Draw(screen)
	g.drawPlayerShield(screen)
	g.drawScreenFlash(screen)
	g.drawHUD(screen)
	if g.paused {
		g.drawPauseOverlay(screen)
	}
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
		g.spawnPlayerBullets()
		g.shootCooldown = config.ShootCooldownTicks
		g.playSound(soundShoot)
	}
}

func (g *Game) spawnPlayerBullets() {
	if g.doubleShotTicks > 0 {
		g.bullets = append(g.bullets,
			entities.NewBullet(g.player.CenterX()-15, g.player.Y-12, g.images.Bullet),
			entities.NewBullet(g.player.CenterX()+15, g.player.Y-12, g.images.Bullet),
		)
		return
	}

	g.bullets = append(g.bullets, entities.NewBullet(
		g.player.CenterX(),
		g.player.Y-12,
		g.images.Bullet,
	))
}

func (g *Game) updateProfessorShip() {
	if g.professorShip == nil {
		g.updateProfessorSpawn()
		return
	}

	if g.professorActiveTicks > 0 {
		g.professorActiveTicks--
	}
	if g.professorActiveTicks == 0 {
		g.scheduleProfessorRespawn()
		return
	}

	g.professorShip.Update(g.player, g.professorFollowSpeed(), g.professorEntrySpeed())
}

func (g *Game) updateProfessorSpawn() {
	if g.images.ProfessorShip == nil || g.professorSpawnTicks < 0 {
		return
	}
	if g.professorSpawnTicks > 0 {
		g.professorSpawnTicks--
		return
	}

	g.professorVisits++
	g.professorShip = entities.NewProfessorShip(g.images.ProfessorShip, g.player, g.professorShipHealth())
	g.professorActiveTicks = g.professorActiveTicksForVisit()
	g.professorShootCooldown = config.ProfessorShipShootMinTicks
	g.professorSpawnTicks = -1
	g.playSound(soundAppear)
}

func (g *Game) scheduleProfessorRespawn() {
	g.professorShip = nil
	g.professorProjectiles = g.professorProjectiles[:0]
	g.professorShootCooldown = 0
	g.professorActiveTicks = 0
	g.professorChargeTicks = 0
	g.professorSpawnTicks = config.ProfessorShipRespawnTicks
	g.professorSpawnTotalTicks = config.ProfessorShipRespawnTicks
}

func (g *Game) handleProfessorShooting() {
	if g.professorShip == nil {
		return
	}
	if !g.professorShip.IsReady(g.player) {
		return
	}
	if g.professorChargeTicks > 0 {
		g.professorAimAngle = g.professorProjectileAngle()
		g.professorChargeTicks--
		if g.professorChargeTicks == 0 {
			g.professorProjectiles = append(g.professorProjectiles, g.newProfessorAttack(g.professorAimAngle)...)
			g.professorShootCooldown = g.nextProfessorShootCooldown()
			g.playSound(soundEnemyShoot)
		}
		return
	}
	if g.professorShootCooldown > 0 {
		g.professorShootCooldown--
		return
	}

	g.professorAimAngle = g.professorProjectileAngle()
	g.professorChargeTicks = config.ProfessorShipChargeTicks
}

func (g *Game) professorProjectileAngle() float64 {
	originX := g.professorShip.CenterX()
	originY := g.professorShip.Y + g.professorShip.Height
	targetX := g.player.CenterX()
	targetY := max(g.player.CenterY(), originY+90)
	angle := math.Atan2(targetY-originY, targetX-originX)
	return angle + (g.rng.Float64()*2-1)*g.professorProjectileSpread()
}

func (g *Game) newProfessorProjectile(angle float64) entities.ProfessorProjectile {
	angle = clampProfessorShotAngle(angle)
	originX := g.professorShip.CenterX() + math.Cos(angle)*g.professorShip.Radius()*0.75
	originY := g.professorShip.Y + g.professorShip.Height + math.Sin(angle)*g.professorShip.Radius()*0.75
	velocityX := math.Cos(angle) * g.professorProjectileSpeed()
	velocityY := math.Sin(angle) * g.professorProjectileSpeed()
	return entities.NewProfessorProjectile(originX, originY, velocityX, velocityY, g.images.Bullet)
}

func clampProfessorShotAngle(angle float64) float64 {
	const minDownwardAngle = 0.42
	maxDownwardAngle := math.Pi - minDownwardAngle

	if angle < minDownwardAngle {
		return minDownwardAngle
	}
	if angle > maxDownwardAngle {
		return maxDownwardAngle
	}
	return angle
}

func (g *Game) newProfessorAttack(angle float64) []entities.ProfessorProjectile {
	level := g.professorVisitLevel()
	switch {
	case level >= 3:
		return []entities.ProfessorProjectile{
			g.newProfessorProjectile(angle - 0.38),
			g.newProfessorProjectile(angle - 0.19),
			g.newProfessorProjectile(angle),
			g.newProfessorProjectile(angle + 0.19),
			g.newProfessorProjectile(angle + 0.38),
		}
	case level >= 1:
		return []entities.ProfessorProjectile{
			g.newProfessorProjectile(angle - 0.16),
			g.newProfessorProjectile(angle),
			g.newProfessorProjectile(angle + 0.16),
		}
	default:
		return []entities.ProfessorProjectile{g.newProfessorProjectile(angle)}
	}
}

func (g *Game) nextProfessorShootCooldown() int {
	minTicks := max(10, config.ProfessorShipShootMinTicks-g.difficultyLevel()*4-g.professorVisitLevel()*2)
	maxTicks := max(minTicks+6, config.ProfessorShipShootMaxTicks-g.difficultyLevel()*5-g.professorVisitLevel()*2)
	extraTicks := maxTicks - minTicks
	if extraTicks <= 0 {
		return minTicks
	}
	return minTicks + g.rng.Intn(extraTicks+1)
}

func (g *Game) spawnEnemies() {
	now := time.Now()
	if now.Sub(g.lastEnemySpawn) < time.Duration(g.enemySpawnMs())*time.Millisecond {
		return
	}
	g.lastEnemySpawn = now

	sprite := g.images.Meteors[g.rng.Intn(len(g.images.Meteors))]
	enemyWidth, enemyHeight := sprite.Size()
	minX, maxX := meteorSpawnBounds(enemyWidth)
	x := float64(minX + g.rng.Intn(maxX-minX+1))
	y := -float64(enemyHeight)
	speed := config.EnemySpeed + g.rng.Float64()*85 + float64(g.difficultyLevel()*14)
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

func (g *Game) updateProfessorProjectiles() {
	alive := g.professorProjectiles[:0]
	for projectileIndex := range g.professorProjectiles {
		g.professorProjectiles[projectileIndex].Update()
		if !g.professorProjectiles[projectileIndex].IsOffScreen() {
			alive = append(alive, g.professorProjectiles[projectileIndex])
		}
	}
	g.professorProjectiles = alive
}

func (g *Game) updatePowerUps() {
	alive := g.powerUps[:0]
	for powerUpIndex := range g.powerUps {
		g.powerUps[powerUpIndex].Update()
		if !g.powerUps[powerUpIndex].IsOffScreen() {
			alive = append(alive, g.powerUps[powerUpIndex])
		}
	}
	g.powerUps = alive
}

func (g *Game) updatePowerUpTimers() {
	if g.doubleShotTicks > 0 {
		g.doubleShotTicks--
	}
	if g.shieldTicks > 0 {
		g.shieldTicks--
	}
}

func (g *Game) updateExplosions() {
	alive := g.explosions[:0]
	for explosionIndex := range g.explosions {
		g.explosions[explosionIndex].Update()
		if !g.explosions[explosionIndex].IsDone() {
			alive = append(alive, g.explosions[explosionIndex])
		}
	}
	g.explosions = alive
}

func (g *Game) updateFloatingTexts() {
	alive := g.floatingTexts[:0]
	for textIndex := range g.floatingTexts {
		g.floatingTexts[textIndex].Update()
		if !g.floatingTexts[textIndex].IsDone() {
			alive = append(alive, g.floatingTexts[textIndex])
		}
	}
	g.floatingTexts = alive
}

func (g *Game) updateScreenFlash() {
	if g.screenFlashTicks > 0 {
		g.screenFlashTicks--
	}
}

func (g *Game) updateEnemies() {
	for i := range g.enemies {
		g.enemies[i].Update()
	}
}

func (g *Game) updateWinScreenMotion() {
	g.updateBullets()
	g.updateProfessorProjectiles()
	g.updatePowerUps()
	g.updateEnemies()
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

func (g *Game) updateRewardMessage() {
	if g.rewardMessageTicks > 0 {
		g.rewardMessageTicks--
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
	g.professorProjectiles = make([]entities.ProfessorProjectile, 0, 16)
	g.powerUps = make([]entities.PowerUp, 0, 8)
	g.explosions = make([]entities.Explosion, 0, 24)
	g.floatingTexts = make([]floatingText, 0, 16)
	g.professorShip = nil
	g.professorSpawnTicks = config.ProfessorShipSpawnDelayTicks
	g.professorSpawnTotalTicks = config.ProfessorShipSpawnDelayTicks
	g.professorActiveTicks = 0
	g.professorChargeTicks = 0
	g.pauseTicks = 0
	g.screenFlashTicks = 0
	g.lastEnemySpawn = time.Now()
	g.shootCooldown = 0
	g.professorShootCooldown = 0
	g.score = 0
	g.lives = config.InitialLives
	g.multiplier = 1
	g.paused = false
	g.gameOver = false
	g.won = false
	g.professorDefeats = 0
	g.professorVisits = 0
	g.doubleShotTicks = 0
	g.shieldTicks = 0
	g.penaltyValue = 0
	g.penaltyTicks = 0
	g.rewardMessageTicks = 0
}

func (g *Game) addExplosion(x, y, radius float64) {
	g.explosions = append(g.explosions, entities.NewExplosion(x, y, radius))
}

func (g *Game) addFloatingText(x, y float64, value string, clr color.RGBA) {
	g.floatingTexts = append(g.floatingTexts, newFloatingText(x, y, value, clr))
}

func (g *Game) updateHighScore() {
	if g.score > g.highScore {
		g.highScore = g.score
	}
}

func (g *Game) updateMusic() {
	if g.sounds == nil {
		return
	}
	g.sounds.updateMusic(g.muted)
}

func (g *Game) defeatProfessorShip() {
	g.addExplosion(g.professorShip.CenterX(), g.professorShip.CenterY(), g.professorShip.Radius()*2.4)
	g.score += g.professorShip.Reward * g.multiplier
	if g.multiplier < config.MaxMultiplier {
		g.multiplier++
	}

	g.professorDefeats++
	g.pauseTicks = config.ProfessorDefeatPauseTicks
	g.playSound(soundExplosion)
	if g.professorDefeats >= config.ProfessorDefeatsToWin {
		g.professorShip = nil
		g.professorProjectiles = g.professorProjectiles[:0]
		g.professorChargeTicks = 0
		g.won = true
		g.playSound(soundVictory)
		return
	}

	g.rewardMessageTicks = config.ProfessorRewardMessageTicks
	g.scheduleProfessorRespawn()
}

func (g *Game) difficultyLevel() int {
	return max(0, min(10, g.score/1500+g.professorDefeats+g.professorVisitLevel()))
}

func (g *Game) enemySpawnMs() int {
	return max(320, config.EnemySpawnMs-g.difficultyLevel()*45)
}

func (g *Game) professorProjectileSpeed() float64 {
	return config.ProfessorProjectileSpeed + float64(g.difficultyLevel()*22+g.professorVisitLevel()*18)
}

func (g *Game) professorProjectileSpread() float64 {
	return max(0.14, config.ProfessorProjectileSpread-float64(g.difficultyLevel())*0.045-float64(g.professorVisitLevel())*0.02)
}

func (g *Game) professorVisitLevel() int {
	return min(config.ProfessorShipMaxVisitLevel, max(0, g.professorVisits-1))
}

func (g *Game) professorShipHealth() int {
	return config.ProfessorShipHealth + g.professorVisitLevel()*config.ProfessorShipHealthPerVisit
}

func (g *Game) professorActiveTicksForVisit() int {
	return config.ProfessorShipActiveTicks + g.professorVisitLevel()*70
}

func (g *Game) professorFollowSpeed() float64 {
	return config.ProfessorShipFollowSpeed + float64(g.professorVisitLevel())*config.ProfessorShipFollowPerVisit
}

func (g *Game) professorEntrySpeed() float64 {
	return config.ProfessorShipEntrySpeed + float64(g.professorVisitLevel())*config.ProfessorShipEntryPerVisit
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
