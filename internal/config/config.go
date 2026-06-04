package config

// Constantes centrais do projeto, como dimensoes da tela e velocidades de objetos. 

const (
	ScreenWidth  = 480
	ScreenHeight = 640

	// A Ebitengine chama Update em uma taxa fixa por padrao: 60 fps
	// Usamos esse delta fixo para escrever velocidades em pixels por segundo 
	DeltaTime = 1.0 / 60.0

	PlayerSpeed  = 280.0
	BulletSpeed  = 430.0
	EnemySpeed   = 55.0
	EnemySpawnMs = 900
)
