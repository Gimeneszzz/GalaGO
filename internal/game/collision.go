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
		for enemyIndex := range g.enemies {
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
			g.lives--
			g.multiplier = 1
			g.player.TakeDamage()
			if g.lives <= 0 {
				g.lives = 0
				g.gameOver = true
			}
			break
		}
	}

	g.bullets = keepBullets(g.bullets, removedBullets)
	g.enemies = keepEnemies(g.enemies, removedEnemies)
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
