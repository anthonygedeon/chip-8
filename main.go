package main

import (
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/urfave/cli/v2"
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

	app := &cli.App{
		Name:    "chip8",
		Usage:   "chip8 [command]",
		Version: "0.1.0",

		Flags: []cli.Flag{

			&cli.BoolFlag{
				Name:    "debug",
				Aliases: []string{"d"},
				Usage:   "Debug running program",
			},

			&cli.BoolFlag{
				Name:    "disassemble",
				Aliases: []string{"disasm"},
				Usage:   "Print out disassembly of program",
			},

		},

		Action: func(c *cli.Context) error {
			return nil	
		},

		Commands: []*cli.Command{
			{
				Name:    "run",
				Aliases: []string{"r"},
				Usage:   "run `path/to/rom`",
				Action: func(c *cli.Context) error {

					rom := c.Args().First()
					cpu.ram.LoadProgram(rom)

					ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
					ebiten.SetWindowTitle("CHIP-8")
					if err := ebiten.RunGame(&vm{}); err != nil {
						return err
					}

					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
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
