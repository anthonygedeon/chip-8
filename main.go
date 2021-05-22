package main

import (
	"fmt"
	"io/ioutil"

	"github.com/hajimehoshi/ebiten/v2"
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

// Chip8
type Chip8 struct {
	Opcode uint16
	I      uint16
	PC     uint16
	SP     uint16

	Memory [4096]byte
	V      [16]byte
	Stack  [16]byte

	Display [64][32]byte
	isDrawing bool

	DelayTimer byte


	Sound func()
}

func main() {

	// ebiten.SetWindowSize(640, 480)
	// ebiten.SetWindowTitle("Chip 8")
	// game := &Game{}
	// if err := ebiten.RunGame(game); err != nil {
	// 	panic(err)
	// }

	chip8 := &Chip8{}

	chip8.Init()

	data := ReadROM("roms/IBMLogo.ch8")

	for i := range data {
		chip8.Memory[i+512] = data[i]
	}

	for {

		chip8.Opcode = chip8.Memory[chip8.PC]<<8 | chip8.Memory[chip8.PC+1]

		switch chip8.Opcode & 0xF000 {
		case 0x0000:
			switch chip8.Opcode & 0x000F {
			case 0xE0:
				// Clears the screen
				// disp_clear()
			case 0xEE:
				// Returns from a subroutine
				// return;
			default:
				panic(fmt.Sprintf("unknown opcode [0x0000]: %x", chip8.Opcode))
			}
		case 0x1000: // 1NNN
			// Jumps to address NNN
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
		case 0x7000: // 7XNN
			// Adds NN to VX. (Carry flag is not changed)
			// Vx += NN
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
				panic(fmt.Sprintf("unknown opcode [0x8000]: %x", chip8.Opcode))
			}
		case 0x9000: // 9XY0
			// Skips the next instruction if VX does not equal VY. (Usually the next instruction is a jump to skip a code block)
			// if(Vx!=Vy)
		case 0xA000: // ANNN
			// Sets I to the address NNN.
			// I = NNN
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
				panic(fmt.Sprintf("unknown opcode [0xE000]: %x", chip8.Opcode))
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
				panic(fmt.Sprintf("unknown opcode [0xF000]: %x", chip8.Opcode))
			}

		default:
			panic(fmt.Sprintf("unknown opcode: %x", chip8.Opcode))
		}

	}

}

func (c Chip8) Init() {
	
}

func (c Chip8) LoadROM() {

}

func ReadROM(filename string) []byte {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return bytes
}


func ClearDisplay() {
	
}

func DrawDisplay() {

}

func RetSub() {

}

func Jump() {
	
}

func CallSub() {

}

func SkipIfEqual() {
	
}

func SkipIfNotEqual() {

}

type Game struct{}

func (g *Game) Update() error {
	// update the logical state
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// render the screen
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	// return the game size
	return 320, 240
}
