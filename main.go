package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 320
	screenHeight = 240
)

var cpu = NewCPU()

type vm struct {
	cpu     CPU
	display Display
	ram     Memory
}

func main() {

	cpu.ram.LoadProgram("roms/KALEID")

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("CHIP-8")
	if err := ebiten.RunGame(&vm{}); err != nil {
		log.Fatal(err)
	}
}

func (vm *vm) Update() error {
	cpu.EmulateCycle()
	return nil
}

func (vm *vm) Draw(screen *ebiten.Image) {
	vm.display.DrawSprite(screen)
}

func (vm *vm) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
