package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Display struct {
	gfx [64][32]byte
}

func (d *Display) Clear() {
	d.gfx = [64][32]byte{}
}

func (d *Display) DrawSprite(screen *ebiten.Image) {
	var pixelImage = ebiten.NewImage(64, 32)

	for row := 0; row < len(cpu.display.gfx); row++ {
		for col := 0; col < len(cpu.display.gfx[row]); col++ {
			var currentColor color.Color

			option := &ebiten.DrawImageOptions{}
			option.GeoM.Translate(float64(row), float64(col))
			option.GeoM.Scale(float64(screen.Bounds().Dx())/64, float64(screen.Bounds().Dy())/32)
			if cpu.display.gfx[row][col] == 1 {
				currentColor = color.White
			} else {
				currentColor = color.Black
			}

			pixelImage.Fill(currentColor)
			screen.DrawImage(pixelImage, option)
		}
	}
}
