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

	ProfessorShipFollowSpeed      = 280.0
	ProfessorShipVerticalDistance = 235.0
	ProfessorShipSpawnDelayTicks  = 420
	ProfessorShipActiveTicks      = 720
	ProfessorShipRespawnTicks     = 600
	ProfessorShipShootMinTicks    = 38
	ProfessorShipShootMaxTicks    = 76
	ProfessorProjectileSpeed      = 240.0
	ProfessorProjectileSpread     = 0.55
	ProfessorShipReward           = 1200
	ProfessorShipHealth           = 10

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
	SmallMeteorEscapePenalty  = 1000
	MediumMeteorEscapePenalty = 2000
	LargeMeteorEscapePenalty  = 3000
)
