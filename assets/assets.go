package assets

import (
	"embed"
	"errors"
	"fmt"
	"image"
	"image/color"
	stddraw "image/draw"
	_ "image/png"
	"io/fs"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	xdraw "golang.org/x/image/draw"
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
var ProfessorShipSprite = loadProfessorShipImage("professornave.png")

var MeteorSprites = mustLoadImages("meteors/*.png")
var StarsSprites = mustLoadImages("stars/*.png")
var PlanetsSprites = mustLoadPlanetImages("planets/*.png")
var UIFontData = mustLoadBytes("fontui.ttf")

func mustLoadImage(name string) *ebiten.Image {
	img := mustLoadDecodedImage(name)
	return ebiten.NewImageFromImage(img)
}

func loadOptionalImage(name string) *ebiten.Image {
	img := loadOptionalDecodedImage(name)
	if img == nil {
		return nil
	}
	return ebiten.NewImageFromImage(img)
}

func loadProfessorShipImage(name string) *ebiten.Image {
	img := loadOptionalDecodedImage(name)
	if img == nil {
		return nil
	}

	cleaned := trimTransparentBounds(removeWhiteBackground(img))
	scaled := scaleToMaxSize(cleaned, 120, 96)
	return ebiten.NewImageFromImage(scaled)
}

func mustLoadDecodedImage(name string) image.Image {
	f, err := files.Open(name)
	if err != nil {
		panic(fmt.Errorf("nao foi possivel abrir a sprite %q: %w", name, err))
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		panic(fmt.Errorf("nao foi possivel decodificar a sprite %q: %w", name, err))
	}

	return img
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

func mustLoadPlanetImages(path string) []*ebiten.Image {
	matches, err := fs.Glob(files, path)
	if err != nil {
		panic(err)
	}

	professorFace := loadOptionalDecodedImage("professor.png")
	if professorFace != nil {
		professorFace = trimTransparentBounds(removeWhiteBackground(professorFace))
	}
	images := make([]*ebiten.Image, len(matches))
	for i, match := range matches {
		planet := mustLoadDecodedImage(match)
		if professorFace != nil {
			planet = drawFaceOnPlanet(planet, professorFace)
		}
		images[i] = ebiten.NewImageFromImage(planet)
	}

	return images
}

func loadOptionalDecodedImage(name string) image.Image {
	f, err := files.Open(name)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}
		panic(fmt.Errorf("nao foi possivel abrir o asset opcional %q: %w", name, err))
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		panic(fmt.Errorf("nao foi possivel decodificar o asset opcional %q: %w", name, err))
	}

	return img
}

func drawFaceOnPlanet(planet image.Image, face image.Image) image.Image {
	planetBounds := planet.Bounds()
	width := planetBounds.Dx()
	height := planetBounds.Dy()
	output := image.NewNRGBA(image.Rect(0, 0, width, height))
	stddraw.Draw(output, output.Bounds(), planet, planetBounds.Min, stddraw.Src)

	diameter := int(float64(min(width, height)) * 0.96)
	if diameter < 1 {
		return output
	}

	centerX := width / 2
	centerY := height / 2
	startX := centerX - diameter/2
	startY := centerY - diameter/2
	radius := float64(diameter) / 2
	faceBounds := squareCenterCrop(face.Bounds())

	for y := 0; y < diameter; y++ {
		for x := 0; x < diameter; x++ {
			dx := float64(x) + 0.5 - radius
			dy := float64(y) + 0.5 - radius
			distance := math.Sqrt(dx*dx + dy*dy)
			if distance > radius {
				continue
			}

			targetX := startX + x
			targetY := startY + y
			if targetX < 0 || targetY < 0 || targetX >= width || targetY >= height {
				continue
			}

			sourceX := faceBounds.Min.X + x*faceBounds.Dx()/diameter
			sourceY := faceBounds.Min.Y + y*faceBounds.Dy()/diameter
			edgeAlpha := 1.0
			if edge := radius - distance; edge < 2 {
				edgeAlpha = edge / 2
			}
			blendPixel(output, targetX, targetY, face.At(sourceX, sourceY), edgeAlpha)
		}
	}

	return output
}

func removeWhiteBackground(img image.Image) *image.NRGBA {
	bounds := img.Bounds()
	output := image.NewNRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))

	for y := 0; y < bounds.Dy(); y++ {
		for x := 0; x < bounds.Dx(); x++ {
			pixel := color.NRGBAModel.Convert(img.At(bounds.Min.X+x, bounds.Min.Y+y)).(color.NRGBA)
			if isNearlyWhite(pixel) {
				pixel.A = 0
			}
			output.SetNRGBA(x, y, pixel)
		}
	}

	return output
}

func isNearlyWhite(pixel color.NRGBA) bool {
	if pixel.A == 0 {
		return true
	}
	minChannel := min(pixel.R, min(pixel.G, pixel.B))
	maxChannel := max(pixel.R, max(pixel.G, pixel.B))
	return minChannel >= 238 && maxChannel-minChannel <= 18
}

func trimTransparentBounds(img *image.NRGBA) *image.NRGBA {
	bounds := img.Bounds()
	minX := bounds.Max.X
	minY := bounds.Max.Y
	maxX := bounds.Min.X
	maxY := bounds.Min.Y

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if img.NRGBAAt(x, y).A <= 8 {
				continue
			}
			minX = min(minX, x)
			minY = min(minY, y)
			maxX = max(maxX, x)
			maxY = max(maxY, y)
		}
	}

	if minX > maxX || minY > maxY {
		return img
	}

	padding := 4
	minX = max(bounds.Min.X, minX-padding)
	minY = max(bounds.Min.Y, minY-padding)
	maxX = min(bounds.Max.X-1, maxX+padding)
	maxY = min(bounds.Max.Y-1, maxY+padding)

	croppedBounds := image.Rect(0, 0, maxX-minX+1, maxY-minY+1)
	cropped := image.NewNRGBA(croppedBounds)
	stddraw.Draw(cropped, croppedBounds, img, image.Point{X: minX, Y: minY}, stddraw.Src)
	return cropped
}

func scaleToMaxSize(img image.Image, maxWidth, maxHeight int) image.Image {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	if width <= maxWidth && height <= maxHeight {
		return img
	}

	scale := min(float64(maxWidth)/float64(width), float64(maxHeight)/float64(height))
	scaledWidth := max(1, int(math.Round(float64(width)*scale)))
	scaledHeight := max(1, int(math.Round(float64(height)*scale)))
	output := image.NewNRGBA(image.Rect(0, 0, scaledWidth, scaledHeight))
	xdraw.CatmullRom.Scale(output, output.Bounds(), img, bounds, xdraw.Over, nil)
	return output
}

func squareCenterCrop(bounds image.Rectangle) image.Rectangle {
	size := min(bounds.Dx(), bounds.Dy())
	left := bounds.Min.X + (bounds.Dx()-size)/2
	top := bounds.Min.Y + (bounds.Dy()-size)/2
	return image.Rect(left, top, left+size, top+size)
}

func blendPixel(img *image.NRGBA, x, y int, source color.Color, opacity float64) {
	sourceColor := color.NRGBAModel.Convert(source).(color.NRGBA)
	if sourceColor.A == 0 || opacity <= 0 {
		return
	}

	destination := img.NRGBAAt(x, y)
	sourceAlpha := (float64(sourceColor.A) / 255) * opacity
	destinationAlpha := float64(destination.A) / 255
	outputAlpha := sourceAlpha + destinationAlpha*(1-sourceAlpha)
	if outputAlpha == 0 {
		img.SetNRGBA(x, y, color.NRGBA{})
		return
	}

	red := (float64(sourceColor.R)*sourceAlpha + float64(destination.R)*destinationAlpha*(1-sourceAlpha)) / outputAlpha
	green := (float64(sourceColor.G)*sourceAlpha + float64(destination.G)*destinationAlpha*(1-sourceAlpha)) / outputAlpha
	blue := (float64(sourceColor.B)*sourceAlpha + float64(destination.B)*destinationAlpha*(1-sourceAlpha)) / outputAlpha

	img.SetNRGBA(x, y, color.NRGBA{
		R: uint8(math.Round(red)),
		G: uint8(math.Round(green)),
		B: uint8(math.Round(blue)),
		A: uint8(math.Round(outputAlpha * 255)),
	})
}

func mustLoadBytes(name string) []byte {
	data, err := files.ReadFile(name)
	if err != nil {
		panic(fmt.Errorf("nao foi possivel abrir o asset %q: %w", name, err))
	}
	return data
}
