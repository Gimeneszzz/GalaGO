package game

import (
	"testing"

	"galago/internal/config"
	"galago/internal/entities"
)

func TestEscapedEnemyResetsMultiplier(t *testing.T) {
	g := &Game{
		score:      5000,
		multiplier: 7,
		enemies: []entities.Enemy{
			{
				Y:             config.ScreenHeight + 1,
				EscapePenalty: config.MediumMeteorEscapePenalty,
			},
		},
	}

	g.removeEscapedEnemies()

	if g.multiplier != 1 {
		t.Fatalf("multiplicador = %d, esperado 1", g.multiplier)
	}
	if g.score != 3000 {
		t.Fatalf("pontuacao = %d, esperado 3000", g.score)
	}
	if len(g.enemies) != 0 {
		t.Fatalf("meteorito que escapou deveria ter sido removido")
	}
}

func TestEscapedEnemyCanMakeScoreNegative(t *testing.T) {
	g := &Game{
		score:      500,
		multiplier: 3,
		enemies: []entities.Enemy{
			{
				Y:             config.ScreenHeight + 1,
				EscapePenalty: config.SmallMeteorEscapePenalty,
			},
		},
	}

	g.removeEscapedEnemies()

	if g.score != -500 {
		t.Fatalf("pontuacao = %d, esperado -500", g.score)
	}
}

func TestNegativeScoreDoesNotMakeDifficultyNegative(t *testing.T) {
	g := &Game{
		score: -30000,
	}

	if got := g.difficultyLevel(); got != 0 {
		t.Fatalf("difficultyLevel() = %d, esperado 0", got)
	}

	if got := g.enemySpawnMs(); got != config.EnemySpawnMs {
		t.Fatalf("enemySpawnMs() = %d, esperado %d", got, config.EnemySpawnMs)
	}
}

func TestMeteorSpawnBoundsRespectSideMargins(t *testing.T) {
	const meteorWidth = 120

	minX, maxX := meteorSpawnBounds(meteorWidth)

	if minX != config.SpawnSideMargin {
		t.Fatalf("limite esquerdo = %d, esperado %d", minX, config.SpawnSideMargin)
	}

	wantMaxX := config.ScreenWidth - config.SpawnSideMargin - meteorWidth
	if maxX != wantMaxX {
		t.Fatalf("limite direito = %d, esperado %d", maxX, wantMaxX)
	}

	if maxX+meteorWidth > config.ScreenWidth-config.SpawnSideMargin {
		t.Fatalf("meteoro ultrapassa a margem direita")
	}
}
