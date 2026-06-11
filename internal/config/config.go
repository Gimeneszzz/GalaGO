package config

// Constantes centrais do projeto, como dimensoes da tela e regras de gameplay.

const (
	// A resolucao logica permanece constante mesmo quando a janela real e
	// redimensionada.
	ScreenWidth  = 800
	ScreenHeight = 600

	// A Ebitengine chama Update em uma taxa fixa por padrao: 60 fps
	// Usamos esse delta fixo para escrever velocidades em pixels por segundo.
	DeltaTime = 1.0 / 60.0

	PlayerSpeed     = 420.0
	BulletSpeed     = 560.0
	EnemySpeed      = 75.0
	EnemySpawnMs    = 750
	SpawnSideMargin = 20

	InitialLives  = 3
	MaxMultiplier = 10

	// Duracao de efeitos em ticks (1 tick = 1 chamada de Update)
	ShootCooldownTicks     = 12
	InvulnerabilityTicks   = 90
	MeteorDamageFlashTicks = 8
	PenaltyMessageTicks    = 60

	// Recompensas por destruir meteoros
	SmallMeteorReward  = 100
	MediumMeteorReward = 200
	LargeMeteorReward  = 300

	// Penalidades por deixar meteoros escaparem
	SmallMeteorEscapePenalty  = 100
	MediumMeteorEscapePenalty = 200
	LargeMeteorEscapePenalty  = 300
)
