package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"galago/internal/config"
	"galago/internal/entities"
)

// handleCollisions concentra as regras de interacao entre objetos.
// colisoes: tiro/meteoro e jogador/meteoro.
func (g *Game) handleCollisions() {
	removedBullets := make([]bool, len(g.bullets))
	removedEnemies := make([]bool, len(g.enemies))
	removedProfessorProjectiles := make([]bool, len(g.professorProjectiles))
	spawnedFragments := make([]entities.Enemy, 0, 4)

	if g.professorShip != nil {
		professorReady := g.professorShip.IsReady(g.player)
		for bulletIndex := range g.bullets {
			if removedBullets[bulletIndex] {
				continue
			}

			bullet := g.bullets[bulletIndex]
			if !circlesOverlap(
				bullet.CenterX(), bullet.CenterY(), bullet.Radius(),
				g.professorShip.CenterX(), g.professorShip.CenterY(), g.professorShip.Radius(),
			) {
				continue
			}

			removedBullets[bulletIndex] = true
			if !professorReady {
				g.addFloatingText(g.professorShip.CenterX(), g.professorShip.CenterY()-28, "Escudo", color.RGBA{R: 120, G: 210, B: 255, A: 255})
				g.playSound(soundDamage)
				continue
			}
			if g.professorShip.TakeDamage() {
				g.addFloatingText(g.professorShip.CenterX(), g.professorShip.CenterY()-26, "-1", color.RGBA{R: 255, G: 95, B: 120, A: 255})
				g.defeatProfessorShip()
				break
			}
			g.addFloatingText(g.professorShip.CenterX(), g.professorShip.CenterY()-26, "-1", color.RGBA{R: 255, G: 95, B: 120, A: 255})
			g.playSound(soundDamage)
		}
	}

	for bulletIndex := range g.bullets {
		for enemyIndex := range g.enemies {
			if removedBullets[bulletIndex] || removedEnemies[enemyIndex] {
				continue
			}

			bullet := g.bullets[bulletIndex]
			enemy := &g.enemies[enemyIndex]
			if !circlesOverlap(
				bullet.CenterX(), bullet.CenterY(), bullet.Radius(),
				enemy.CenterX(), enemy.CenterY(), enemy.Radius(),
			) {
				continue
			}

			removedBullets[bulletIndex] = true
			if enemy.TakeDamage() {
				removedEnemies[enemyIndex] = true
				g.score += enemy.Reward * g.multiplier
				if g.multiplier < config.MaxMultiplier {
					g.multiplier++
				}
				g.addExplosion(enemy.CenterX(), enemy.CenterY(), enemy.Radius())
				spawnedFragments = append(spawnedFragments, g.maybeSpawnMeteorFragments(*enemy)...)
				g.maybeDropPowerUp(enemy.CenterX(), enemy.CenterY())
				g.playSound(soundExplosion)
			}
		}
	}

	removedPowerUps := make([]bool, len(g.powerUps))
	for powerUpIndex := range g.powerUps {
		powerUp := g.powerUps[powerUpIndex]
		if !circlesOverlap(
			g.player.CenterX(), g.player.CenterY(), g.player.Radius(),
			powerUp.CenterX(), powerUp.CenterY(), powerUp.Radius(),
		) {
			continue
		}

		removedPowerUps[powerUpIndex] = true
		g.collectPowerUp(powerUp.Kind)
	}

	if !g.player.IsInvulnerable() {
		if g.professorShip != nil && circlesOverlap(
			g.player.CenterX(), g.player.CenterY(), g.player.Radius(),
			g.professorShip.CenterX(), g.professorShip.CenterY(), g.professorShip.Radius(),
		) {
			g.damagePlayer()
		}

		for projectileIndex := range g.professorProjectiles {
			if g.player.IsInvulnerable() {
				break
			}

			projectile := g.professorProjectiles[projectileIndex]
			if !circlesOverlap(
				g.player.CenterX(), g.player.CenterY(), g.player.Radius(),
				projectile.CenterX(), projectile.CenterY(), projectile.Radius(),
			) {
				continue
			}

			removedProfessorProjectiles[projectileIndex] = true
			g.damagePlayer()
		}

		for enemyIndex := range g.enemies {
			if g.player.IsInvulnerable() {
				break
			}
			if removedEnemies[enemyIndex] {
				continue
			}

			enemy := g.enemies[enemyIndex]
			if !circlesOverlap(
				g.player.CenterX(), g.player.CenterY(), g.player.Radius(),
				enemy.CenterX(), enemy.CenterY(), enemy.Radius(),
			) {
				continue
			}

			removedEnemies[enemyIndex] = true
			g.damagePlayer()
		}
	}

	g.bullets = keepBullets(g.bullets, removedBullets)
	g.enemies = keepEnemies(g.enemies, removedEnemies)
	g.enemies = append(g.enemies, spawnedFragments...)
	g.professorProjectiles = keepProfessorProjectiles(g.professorProjectiles, removedProfessorProjectiles)
	g.powerUps = keepPowerUps(g.powerUps, removedPowerUps)
}

func (g *Game) damagePlayer() {
	if g.shieldTicks > 0 {
		g.shieldTicks = 0
		g.player.TakeDamage()
		g.addExplosion(g.player.CenterX(), g.player.CenterY(), g.player.Radius()*1.3)
		g.screenFlashTicks = 10
		g.playSound(soundDamage)
		return
	}

	g.lives--
	g.multiplier = 1
	g.player.TakeDamage()
	g.addExplosion(g.player.CenterX(), g.player.CenterY(), g.player.Radius()*1.6)
	g.screenFlashTicks = 14
	if g.lives <= 0 {
		g.lives = 0
		g.gameOver = true
		g.playSound(soundDefeat)
		return
	}
	g.playSound(soundDamage)
}

func (g *Game) maybeDropPowerUp(x, y float64) {
	if g.rng.Float64() > config.PowerUpDropChance {
		return
	}

	kind := entities.PowerUpDoubleShot
	if g.rng.Intn(2) == 0 {
		kind = entities.PowerUpShield
	}
	g.powerUps = append(g.powerUps, entities.NewPowerUp(x, y, kind))
}

func (g *Game) maybeSpawnMeteorFragments(enemy entities.Enemy) []entities.Enemy {
	if enemy.MaxHealth < 3 || len(g.images.Meteors) == 0 || g.rng.Float64() > 0.45 {
		return nil
	}

	fragments := make([]entities.Enemy, 0, 2)
	for fragmentIndex := 0; fragmentIndex < 2; fragmentIndex++ {
		sprite := g.randomSmallMeteorSprite()
		offset := float64(fragmentIndex*2-1) * enemy.Radius() * 0.35
		speed := enemy.Speed + 55 + g.rng.Float64()*45
		rotationSpeed := (g.rng.Float64()*2 - 1) * 2.2
		fragments = append(fragments, entities.NewEnemy(enemy.CenterX()+offset, enemy.CenterY(), speed, rotationSpeed, sprite))
	}
	return fragments
}

func (g *Game) randomSmallMeteorSprite() *ebiten.Image {
	smallSprites := make([]*ebiten.Image, 0, len(g.images.Meteors))
	for _, sprite := range g.images.Meteors {
		width, height := sprite.Size()
		if max(width, height) <= 60 {
			smallSprites = append(smallSprites, sprite)
		}
	}
	if len(smallSprites) == 0 {
		return g.images.Meteors[g.rng.Intn(len(g.images.Meteors))]
	}
	return smallSprites[g.rng.Intn(len(smallSprites))]
}

func (g *Game) collectPowerUp(kind entities.PowerUpKind) {
	switch kind {
	case entities.PowerUpDoubleShot:
		g.specialShotTicks = config.PowerUpDurationTicks
		g.addFloatingText(g.player.CenterX(), g.player.Y-24, "Tiro Máximo!", color.RGBA{R: 80, G: 210, B: 255, A: 255})
	case entities.PowerUpShield:
		g.shieldTicks = config.PowerUpDurationTicks
		g.addFloatingText(g.player.CenterX(), g.player.Y-24, "Escudo!", color.RGBA{R: 100, G: 245, B: 145, A: 255})
	}
	g.playSound(soundAppear)
}

func circlesOverlap(ax, ay, ar, bx, by, br float64) bool {
	dx := ax - bx
	dy := ay - by
	radius := ar + br
	return dx*dx+dy*dy <= radius*radius
}

func keepBullets(items []entities.Bullet, removed []bool) []entities.Bullet {
	remaining := items[:0]
	for i, item := range items {
		if !removed[i] {
			remaining = append(remaining, item)
		}
	}
	return remaining
}

func keepEnemies(items []entities.Enemy, removed []bool) []entities.Enemy {
	remaining := items[:0]
	for i, item := range items {
		if !removed[i] {
			remaining = append(remaining, item)
		}
	}
	return remaining
}

func keepProfessorProjectiles(items []entities.ProfessorProjectile, removed []bool) []entities.ProfessorProjectile {
	remaining := items[:0]
	for i, item := range items {
		if !removed[i] {
			remaining = append(remaining, item)
		}
	}
	return remaining
}

func keepPowerUps(items []entities.PowerUp, removed []bool) []entities.PowerUp {
	remaining := items[:0]
	for i, item := range items {
		if !removed[i] {
			remaining = append(remaining, item)
		}
	}
	return remaining
}
