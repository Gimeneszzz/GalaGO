package assets

import (
	"embed"
	"fmt"
	"image"
	_ "image/png"
	"io/fs"

	"github.com/hajimehoshi/ebiten/v2"
)

// Mantemos este arquivo separado de internal/assets porque ele foi fornecido
// junto com as sprites e é mais fácil de usar go:embed aqui do que tentar mover
// os arquivos para internal/assets e usar go:embed lá, o que exigiria
// um caminho mais complexo para acessar

// O pacote internal/assets transforma esses recursos em uma
// estrutura mais conveniente para o jogo.

//go:embed *.png *.ttf meteors/*.png stars/*.png planets/*.png
var files embed.FS

var PlayerSprite = mustLoadImage("player.png")
var GopherPlayerSprite = mustLoadImage("go_player.png")
var LaserSprite = mustLoadImage("laser.png")

var MeteorSprites = mustLoadImages("meteors/*.png")
var StarsSprites = mustLoadImages("stars/*.png")
var PlanetsSprites = mustLoadImages("planets/*.png")
var UIFontData = mustLoadBytes("fontui.ttf")

func mustLoadImage(name string) *ebiten.Image {
	f, err := files.Open(name)
	if err != nil {
		panic(fmt.Errorf("nao foi possivel abrir a sprite %q: %w", name, err))
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		panic(fmt.Errorf("nao foi possivel decodificar a sprite %q: %w", name, err))
	}

	return ebiten.NewImageFromImage(img)
}

func mustLoadImages(path string) []*ebiten.Image {
	matches, err := fs.Glob(files, path)
	if err != nil {
		panic(err)
	}

	images := make([]*ebiten.Image, len(matches))
	for i, match := range matches {
		images[i] = mustLoadImage(match)
	}

	return images
}

func mustLoadBytes(name string) []byte {
	data, err := files.ReadFile(name)
	if err != nil {
		panic(fmt.Errorf("nao foi possivel abrir o asset %q: %w", name, err))
	}
	return data
}
