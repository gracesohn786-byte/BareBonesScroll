package main

import (
	"fmt"
	_ "image/png" // underscore is for making the code compilable from the error of unused imports
	"io"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	_ "github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type scrollDemo struct {
	text            string
	drawFace        font.Face
	font            font.Face
	player          *ebiten.Image
	background      *ebiten.Image
	backgroundXView int
	Stage           int
}

func (demo *scrollDemo) Update() error { // appropriate in the same interface to use all pointers or all none pointers
	if demo.Stage == 1 {
		change := inpututil.IsKeyJustPressed(ebiten.KeySpace)
		if change {
			demo.Stage = 0
		}
	} else {
		backgroundWidth := demo.background.Bounds().Dx()
		maxX := backgroundWidth * 2 // 3 copies of image
		demo.backgroundXView -= 4
		demo.backgroundXView %= maxX // when remainder=0, backgroundXView moves back to the starting point
		//inpututil.IsKeyJustPressed(ebiten.KeyLeft)
	}
	return nil
}

func (demo *scrollDemo) Draw(screen *ebiten.Image) {
	if demo.Stage == 1 {
		drawFace := text.NewGoXFace(demo.font)
		textOpts := &text.DrawOptions{
			DrawImageOptions: ebiten.DrawImageOptions{},
			LayoutOptions:    text.LayoutOptions{},
		}
		textOpts.GeoM.Reset() // GeoM is field of DrawImageOptions, which is a subfield of textOpts, so textOpts can access it directly
		textOpts.GeoM.Translate(350, 450)
		textOpts.ColorScale.ScaleWithColor(colornames.Red)
		text.Draw(screen, demo.text, drawFace, textOpts)
	} else {
		drawOps := ebiten.DrawImageOptions{}
		const repeat = 3
		backgroundWidth := demo.background.Bounds().Dx()
		for count := 0; count < repeat; count += 1 {
			drawOps.GeoM.Reset()
			drawOps.GeoM.Translate(float64(backgroundWidth*count), // Images being translated 3 times (3 copies)
				float64(-1000)) // Translate(x,y float64) start from y=-1000 to draw everything from mid-point to below
			drawOps.GeoM.Translate(float64(demo.backgroundXView), 0)
			screen.DrawImage(demo.background, &drawOps)
		}
	}
}

func (s scrollDemo) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	//	ebiten.SetWindowSize(1000, 1000)
	ebiten.SetFullscreen(true)
	ebiten.SetWindowTitle("Scroller Example")
	//New image from file returns image as image.Image (_) and ebiten.Image
	backgroundPict, _, err := ebitenutil.NewImageFromFile("background.png") // 1st var is ebiten image, 2nd var is std go image, 3rd is go's error object (any struct that implements it is an error)
	if err != nil {
		fmt.Println("Unable to load background image:", err)
	}

	face := LoadFont("Square-Black.ttf", 55)
	demo := scrollDemo{
		player:     nil,
		background: backgroundPict,
		Stage:      1,
		drawFace:   face,
		// backgroundXView is zero (default int val) to start off
	}
	err = ebiten.RunGame(&demo)
	if err != nil {
		fmt.Println("Failed to run game", err)
	}
}

func LoadFont(fontFile string, size float64) font.Face {
	fileHandle, err := os.Open(fontFile)
	if err != nil {
		log.Fatal(err)
	}
	fontData, err := io.ReadAll(fileHandle)
	if err != nil {
		log.Fatal(err)
	}
	ttFont, err := opentype.Parse(fontData)
	if err != nil {
		log.Fatal(err)
	}
	fontFace, err := opentype.NewFace(ttFont, &opentype.FaceOptions{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	return fontFace
}
