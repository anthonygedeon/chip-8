package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 320
	screenHeight = 240
)

var chip8FontSet = [80]byte{
	0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
	0x20, 0x60, 0x20, 0x20, 0x70, // 1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
	0x90, 0x90, 0xF0, 0x10, 0x10, // 4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
	0xF0, 0x10, 0x20, 0x40, 0x40, // 7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
	0xF0, 0x90, 0xF0, 0x90, 0x90, // A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
	0xF0, 0x80, 0x80, 0x80, 0xF0, // C
	0xE0, 0x90, 0x90, 0x90, 0xE0, // D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
	0xF0, 0x80, 0xF0, 0x80, 0x80, // F
}

type Pixel struct {
	X, Y byte
}

func (p Pixel) Draw() {

}

type Game struct{}

// Chip8
type Chip8 struct {
	Opcode uint16

	// point at locations in memory
	I uint16

	// points at the current instruction in memory
	PC uint16

	// The 4096 bytes of memory.
	//
	// Memory Map:
	// +---------------+= 0xFFF (4095) End of Chip-8 RAM
	// |               |
	// |               |
	// |               |
	// |               |
	// |               |
	// | 0x200 to 0xFFF|
	// |     Chip-8    |
	// | Program / Data|
	// |     Space     |
	// |               |
	// |               |
	// |               |
	// +- - - - - - - -+= 0x600 (1536) Start of ETI 660 Chip-8 programs
	// |               |
	// |               |
	// |               |
	// +---------------+= 0x200 (512) Start of most Chip-8 programs
	// | 0x000 to 0x1FF|
	// | Reserved for  |
	// |  interpreter  |
	// +---------------+= 0x000 (0) Start of Chip-8 RAM
	Memory [4096]byte
	// 16 8-bit (one byte) general-purpose variable registers numbered
	V [16]byte

	// call subroutines/functions and return from them
	Stack [16]byte

	// Stack pointer
	SP uint16

	// Display: 64 x 32 pixels
	Display [64 * 32]byte

	// HEX Keypad
	Keypad [16]byte

	DelayTimer byte
}

func main() {

	// ebiten.SetWindowSize(640, 480)
	// ebiten.SetWindowTitle("Chip 8")
	// game := &Game{}
	// if err := ebiten.RunGame(game); err != nil {
	// 	panic(err)
	// }

	chip8 := Chip8{}

	chip8.Init()

	data := ReadROM("roms/IBMLOGO.ch8")

	for i := range data {
		chip8.Memory[512+i] = data[i]
	}

	for {

		chip8.Opcode = (uint16(chip8.Memory[chip8.PC]) << 8) | uint16(chip8.Memory[chip8.PC+1])

		switch chip8.Opcode & 0xF000 {
		case 0x0000:
			switch chip8.Opcode & 0x000F {
			case 0x0000:
				// Clears the screen
				// disp_clear()
				for i := 0; i < len(chip8.Display); i++ {
					chip8.Display[i] = 0
				}
				chip8.PC += 2
			case 0x000E:
				// Returns from a subroutine
				// return;
			default:
				panic(fmt.Sprintf("unknown opcode [0x0000]: 0x%X\n", chip8.Opcode))
			}
		case 0x1000: // 1NNN
			// Jumps to address NNN
			chip8.PC = chip8.Opcode & 0x0FFF
		case 0x2000: // 2NNN
			// Calls subroutine at NNN.
		case 0x3000: // 3XNN
			// Skips the next instruction if VX equals NN. (Usually the next instruction is a jump to skip a code block)
			// if(Vx==NN)
		case 0x4000: // 4XNN
			// Skips the next instruction if VX does not equal NN. (Usually the next instruction is a jump to skip a code block)
			// if(Vx!=NN)
		case 0x5000: // 5XY0
			// Skips the next instruction if VX equals VY. (Usually the next instruction is a jump to skip a code block)
			// if(Vx==Vy)
		case 0x6000: // 6XNN
			// Sets VX to NN
			// Vx = NN
			v := chip8.Opcode >> 8 & 0x000F
			nn := byte(chip8.Opcode & 0x00FF)
			chip8.V[v] = nn
			chip8.PC += 2
		case 0x7000: // 7XNN
			// Adds NN to VX. (Carry flag is not changed)
			// Vx += NN
			v := chip8.Opcode >> 8 & 0x000F
			nn := byte(chip8.Opcode & 0x00FF)
			chip8.V[v] += nn
			chip8.PC += 2
		case 0x8000:
			switch chip8.Opcode & 0x000F {
			case 0x0:
				// Sets VX to the value of VY.
				// Vx=Vy
			case 0x1:
				// Sets VX to VX or VY. (Bitwise OR operation)
				// Vx=Vx|Vy
			case 0x2:
				// Sets VX to VX and VY. (Bitwise AND operation)
				// Vx=Vx&Vy
			case 0x3:
				// Sets VX to VX xor VY.
				// Vx=Vx^Vy
			case 0x4:
				// Adds VY to VX. VF is set to 1 when there's a carry, and to 0 when there is not.
				// Vx += Vy
			case 0x5:
				// VY is subtracted from VX. VF is set to 0 when there's a borrow, and 1 when there is not.
				// Vx -= Vy
			case 0x6:
				// Stores the least significant bit of VX in VF and then shifts VX to the right by 1.
				// Vx>>=1
			case 0x7:
				// Sets VX to VY minus VX. VF is set to 0 when there's a borrow, and 1 when there is not
				// Vx=Vy-Vx
			case 0xE:
				// Stores the most significant bit of VX in VF and then shifts VX to the left by 1.
				// Vx<<=1
			default:
				panic(fmt.Sprintf("unknown opcode [0x8000]: 0x%X\n", chip8.Opcode))
			}
		case 0x9000: // 9XY0
			// Skips the next instruction if VX does not equal VY. (Usually the next instruction is a jump to skip a code block)
			// if(Vx!=Vy)
		case 0xA000: // ANNN
			// Sets I to the address NNN.
			// I = NNN
			nnn := chip8.Opcode & 0x0FFF
			chip8.I = nnn
			chip8.PC += 2
		case 0xB000: // BNNN
			// Jumps to the address NNN plus V0.
			// PC=V0+NNN
		case 0xC000: // CXNN
			// Sets VX to the result of a bitwise and operation on a random number (Typically: 0 to 255) and NN.
			// Vx=rand()&NN
		case 0xD000: // DXYN
			// Draws a sprite at coordinate (VX, VY) that has a width of 8 pixels and a height of N+1 pixels
			// Each row of 8 pixels is read as bit-coded starting from memory location I; I value does not change after the execution of this instruction
			// As described above, VF is set to 1 if any screen pixels are flipped from set to unset when the sprite is drawn, and to 0 if that does not happen
			// draw(Vx,Vy,N)
		case 0xE000:
			switch chip8.Opcode & 0x00FF {
			case 0x9E:
				// Skips the next instruction if the key stored in VX is pressed. (Usually the next instruction is a jump to skip a code block)
				// if(key()==Vx)
			case 0xA1:
				// Skips the next instruction if the key stored in VX is not pressed. (Usually the next instruction is a jump to skip a code block)
				// if(key()!=Vx)
			default:
				panic(fmt.Sprintf("unknown opcode [0xE000]: 0x%X\n", chip8.Opcode))
			}
		case 0xF000:
			switch chip8.Opcode & 0x00FF {
			case 0x07:
				// Sets VX to the value of the delay timer.
				// Vx = get_delay()
			case 0x0A:
				// A key press is awaited, and then stored in VX. (Blocking Operation. All instruction halted until next key event)
				// Vx = get_key()
			case 0x15:
				// Sets the delay timer to VX.
				// delay_timer(Vx)
			case 0x18:
				// Sets the sound timer to VX.
				// sound_timer(Vx)
			case 0x1E:
				// Adds VX to I. VF is not affected
				// I +=Vx
			case 0x29:
				// Sets I to the location of the sprite for the character in VX. Characters 0-F (in hexadecimal) are represented by a 4x5 font.
				// I=sprite_addr[Vx]
			case 0x33:
				// Stores the binary-coded decimal representation of VX, with the most
				// significant of three digits at the address in I, the middle digit at I plus 1, and the
				// least significant digit at I plus 2. (In other words, take the decimal
				// representation of VX, place the hundreds digit in memory at location in I, the
				// tens digit at location I+1, and the ones digit at location I+2.)
				// set_BCD(Vx)
				// *(I+0)=BCD(3)
				// *(I+1)=BCD(2)
				// *(I+2)=BCD(1)
			case 0x55:
				// Stores V0 to VX (including VX) in memory starting at address I.
				// The offset from I is increased by 1 for each value written, but I itself is left unmodified
				// reg_dump(Vx,&I)
			case 0x65:
				// Fills V0 to VX (including VX) with values from memory starting at address I
				// The offset from I is increased by 1 for each value written, but I itself is left unmodified.
				// reg_load(Vx,&I)
			default:
				panic(fmt.Sprintf("unknown opcode [0xF000]: 0x%X\n", chip8.Opcode))
			}

		default:
			panic(fmt.Sprintf("unknown opcode: 0x%X\n", chip8.Opcode))
		}

	}
}
