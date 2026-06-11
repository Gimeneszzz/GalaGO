package config

// Constantes centrais do projeto, como dimensoes da tela e regras de gameplay.

const (
	// 960x540 usa a proporcao horizontal 16:9. A resolucao logica permanece
	// constante mesmo quando a janela real e redimensionada.
	ScreenWidth  = 960
	ScreenHeight = 540

	// A Ebitengine chama Update em uma taxa fixa por padrao: 60 fps
	// Usamos esse delta fixo para escrever velocidades em pixels por segundo.
	DeltaTime = 1.0 / 60.0

	PlayerSpeed  = 420.0
	BulletSpeed  = 560.0
	EnemySpeed   = 75.0
	EnemySpawnMs = 750

	InitialLives              = 3
	MaxMultiplier             = 10
	ShootCooldownTicks        = 12
	InvulnerabilityTicks      = 90
	MeteorDamageFlashTicks    = 8
	PenaltyMessageTicks       = 60
	SmallMeteorReward         = 100
	MediumMeteorReward        = 200
	LargeMeteorReward         = 300
	SmallMeteorEscapePenalty  = 1000
	MediumMeteorEscapePenalty = 2000
	LargeMeteorEscapePenalty  = 3000
)
