package game

import (
	"galago/internal/config"
	"galago/internal/entities"
)

// handleCollisions concentra as regras de interacao entre objetos.
// colisoes: tiro/meteoro e jogador/meteoro.
func (g *Game) handleCollisions() {
	removedBullets := make([]bool, len(g.bullets))
	removedEnemies := make([]bool, len(g.enemies))
	removedProfessorProjectiles := make([]bool, len(g.professorProjectiles))

	if g.professorShip != nil {
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
			if g.professorShip.TakeDamage() {
				g.score += g.professorShip.Reward * g.multiplier
				if g.multiplier < config.MaxMultiplier {
					g.multiplier++
				}
				g.scheduleProfessorRespawn()
				break
			}
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
			}
		}
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
	g.professorProjectiles = keepProfessorProjectiles(g.professorProjectiles, removedProfessorProjectiles)
}

func (g *Game) damagePlayer() {
	g.lives--
	g.multiplier = 1
	g.player.TakeDamage()
	if g.lives <= 0 {
		g.lives = 0
		g.gameOver = true
	}
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
