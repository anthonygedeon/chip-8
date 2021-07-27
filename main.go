package main

import (
	"fmt"
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
	display Display
}

func main() {

	app := &cli.App{
		Name:    "chip8",
		Usage:   "chip8 [command] [flag]...",
		Version: "0.1.0",

		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "disassemble",
				Aliases: []string{"d"},
				Usage:   "Print out disassembly of `rom`",
			},
		},

		Action: func(c *cli.Context) error {
			if c.Bool("disassemble") {
				rom := c.Args()

				if rom.Len() <= 0 {
					return fmt.Errorf("choose a rom to disassemble")
				}

				file, err := os.Open(rom.First())
				if err != nil {
					return err
				}

				bytes, err := cpu.ram.Read(file)
				if err != nil {
					return err
				}

				cpu.ram.LoadProgram(rom.First())
				for i := 0; i < len(bytes); i++ {
					Disassemble(cpu)
				}
			} else {
				fmt.Printf("Usage: chip8 [command] [flag]...\nTry 'chip8 --help' for more information.\n")
			}

			return nil
		},

		Commands: []*cli.Command{
			{
				Name:  "run",
				Usage: "run `path/to/rom`",
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
