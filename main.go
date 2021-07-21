package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 320
	screenHeight = 240
)

var cpu = NewCPU()

type vm struct {
	//cpu     *CPU
	//display Display
	//ram     *Memory
}

func main() {
	cpu.ram.LoadProgram("test_roms/test_opcode.ch8")

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("CHIP-8")
	if err := ebiten.RunGame(&vm{}); err != nil {
		log.Fatal(err)
	}
}

func (vm *vm) Update() error {
	cpu.FetchOpcode()
	cpu.ExecuteOpcode()
	fmt.Printf("0x%X\n", cpu.opcode)
	return nil
}

func (vm *vm) Draw(screen *ebiten.Image) {

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

func (vm *vm) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
