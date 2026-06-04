package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"galago/internal/config"
	"galago/internal/game"
)

func main() {
	// O main deve ficar pequeno: ele configura a janela, cria o estado inicial
	// do jogo e entrega esse estado para o loop da Ebitengine.
	//
	// A Ebitengine chama automaticamente os metodos Update, Draw e Layout do
	// objeto passado para RunGame. Por isso, a maior parte do jogo fica em
	// internal/game, e nao espalhada no main.go.
	ebiten.SetWindowSize(config.ScreenWidth, config.ScreenHeight)
	ebiten.SetWindowTitle("GalaGO - MVP")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	g, err := game.New()
	if err != nil {
		log.Fatal(err)
	}

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
