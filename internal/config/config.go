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

	ProfessorShipFollowSpeed      = 300.0
	ProfessorShipFollowPerVisit   = 34.0
	ProfessorShipEntrySpeed       = 210.0
	ProfessorShipEntryPerVisit    = 24.0
	ProfessorShipVerticalDistance = 235.0
	ProfessorShipSpawnDelayTicks  = 420
	ProfessorShipActiveTicks      = 720
	ProfessorShipRespawnTicks     = 600
	ProfessorShipWarningTicks     = 120
	ProfessorShipChargeTicks      = 34
	ProfessorShipShootMinTicks    = 38
	ProfessorShipShootMaxTicks    = 76
	ProfessorProjectileSpeed      = 270.0
	ProfessorProjectileSpread     = 0.55
	ProfessorShipReward           = 10000
	ProfessorShipHealth           = 24
	ProfessorShipHealthPerVisit   = 10
	ProfessorShipMaxVisitLevel    = 10
	ProfessorDefeatsToWin         = 4

	InitialLives  = 3
	MaxMultiplier = 10

	// Duracao de efeitos em ticks (1 tick = 1 chamada de Update)
	ShootCooldownTicks     = 12
	InvulnerabilityTicks   = 90
	MeteorDamageFlashTicks = 8
	PenaltyMessageTicks    = 60

	ProfessorRewardMessageTicks = 100
	ProfessorDefeatPauseTicks   = 36
	PowerUpSpeed                = 115.0
	PowerUpDurationTicks        = 600
	PowerUpDropChance           = 0.18

	// Recompensas por destruir meteoros
	SmallMeteorReward  = 100
	MediumMeteorReward = 200
	LargeMeteorReward  = 300

	// Penalidades por deixar meteoros escaparem
	SmallMeteorEscapePenalty  = 1000
	MediumMeteorEscapePenalty = 2000
	LargeMeteorEscapePenalty  = 3000
)
