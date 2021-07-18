package main

import (
	// "flag"
	"fmt"
	"image/color"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

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
	Memory [4096]uint8
	// 16 8-bit (one byte) general-purpose variable registers numbered
	V [16]byte

	// call subroutines/functions and return from them
	Stack [16]uint16

	// Stack pointer
	SP uint16

	// Display: 64 x 32
	Display [64][32]byte

	// Keypad is HEX based: 0x0-0xF
	//  1  2  3  C
	//  4  5  6  D
	//  7  8  9  E
	//  A  0  B  F
	Keypad [16]byte

	DelayTimer byte

	SoundTimer byte
}

func main() {

	chip8 := Chip8{}

	chip8.Init()

	data := ReadROM("roms/TETRIS") // Testing purposes

	for i, d := range data {
		chip8.Memory[512+i] = d
	}

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Chip 8")
	if err := ebiten.RunGame(&chip8); err != nil {
		panic(err)
	}

	// file := flag.String("d", "", "`disassemble` a chip 8 program")
	// flag.Parse()

	// if *file != "" {
	// 	chip8Disassebmler := Chip8{}
	// 	chip8Disassebmler.Init()
	// 	fmt.Println(chip8Disassebmler.PC)
	// 	chip8Disassebmler.PC = 512
	// 	program := ReadROM(*file)
	// 	for i := range program {
	// 		chip8Disassebmler.Memory[512+i] = data[i]
	// 	}
	// 	for _, p := range chip8Disassebmler.Memory[chip8Disassebmler.PC : 512+len(program)] {
	// 		chip8Disassebmler.Disassemble(p, chip8Disassebmler.PC)
	// 		chip8Disassebmler.PC += 2
	// 	}
	// }

}

func (chip8 *Chip8) Update() error {

	chip8.Opcode = uint16(chip8.Memory[chip8.PC])<<8 | uint16(chip8.Memory[chip8.PC+1])
	fmt.Printf("opcode 0x%X\n", chip8.Opcode)
	switch chip8.Opcode & 0xF000 {
	case 0x0000:
		switch chip8.Opcode & 0x000F {
		case 0x0000:
			// Clears the screen
			// disp_clear()
			chip8.Display = [64][32]byte{}
			chip8.PC += 2
		case 0x000E:
			// Returns from a subroutine
			chip8.SP--
			chip8.PC = uint16(chip8.Stack[chip8.SP])
			
			chip8.PC += 2

		default:
			panic(fmt.Sprintf("unknown opcode [0x0000]: 0x%X\n", chip8.Opcode))
		}
	case 0x1000: // 1NNN
		// Jumps to address NNN
		nnn := chip8.Opcode & 0x0FFF
		chip8.PC = nnn
	case 0x2000: // 2NNN
		// Calls subroutine at NNN.
		nnn := chip8.Opcode & 0x0FFF
		chip8.Stack[chip8.SP] = uint16(chip8.PC)
		chip8.SP++
		chip8.PC = nnn
	case 0x3000: // 3XNN
		// Skips the next instruction if VX equals NN. (Usually the next instruction is a jump to skip a code block)
		// if(Vx==NN)
		nn := chip8.Opcode & 0x00FF
		x := chip8.Opcode >> 8 & 0x000F

		if uint16(chip8.V[x]) == nn {
			chip8.PC += 4
		} else {
			chip8.PC += 2
		}

	case 0x4000: // 4XNN
		// Skips the next instruction if VX does not equal NN. (Usually the next instruction is a jump to skip a code block)
		// if(Vx!=NN)
		nn := chip8.Opcode & 0x00FF
		x := chip8.Opcode >> 8 & 0x000F
		if uint16(chip8.V[x]) != nn {
			chip8.PC += 4
		} else {
			chip8.PC += 2
		}
	case 0x5000: // 5XY0
		// Skips the next instruction if VX equals VY. (Usually the next instruction is a jump to skip a code block)
		// if(Vx==Vy)
		x := (chip8.Opcode >> 8) & 0x000F
		y := (chip8.Opcode >> 4) & 0x000F

		if chip8.V[x] == chip8.V[y] {
			chip8.PC += 4
		} else {
			chip8.PC += 2
		}

	case 0x6000: // 6XNN
		// Sets VX to NN
		// Vx = NN
		x := chip8.Opcode >> 8 & 0x000F
		nn := byte(chip8.Opcode & 0x00FF)
		chip8.V[x] = nn

		chip8.PC += 2
	case 0x7000: // 7XNN
		// Adds NN to VX. (Carry flag is not changed)
		// Vx += NN
		x := chip8.Opcode >> 8 & 0x000F
		nn := byte(chip8.Opcode & 0x00FF)
		chip8.V[x] += nn
		chip8.PC += 2
	case 0x8000:
		switch chip8.Opcode & 0x000F {
		case 0x0000:
			// Sets VX to the value of VY.
			// Vx=Vy
			x := (chip8.Opcode >> 8) & 0x000F
			y := (chip8.Opcode >> 4) & 0x000F
			chip8.V[x] += byte(chip8.V[y])
			chip8.PC += 2
		case 0x0001:
			// Sets VX to VX or VY. (Bitwise OR operation)
			// Vx=Vx|Vy

			x := (chip8.Opcode >> 8) & 0x000F
			y := (chip8.Opcode >> 4) & 0x000F

			chip8.V[x] = chip8.V[x] | chip8.V[y]

			chip8.PC += 2
		case 0x0002:
			// Sets VX to VX and VY. (Bitwise AND operation)
			// Vx=Vx&Vy
			x := (chip8.Opcode >> 8) & 0x000F
			y := (chip8.Opcode >> 4) & 0x000F

			chip8.V[x] = chip8.V[x] & chip8.V[y]

			chip8.PC += 2
		case 0x0003:
			// Sets VX to VX xor VY.
			// Vx=Vx^Vy
			x := (chip8.Opcode >> 8) & 0x000F
			y := (chip8.Opcode >> 4) & 0x000F

			chip8.V[x] = chip8.V[x] ^ chip8.V[y]

			chip8.PC += 2
		case 0x0004:
			// Adds VY to VX. VF is set to 1 when there's a carry, and to 0 when there is not.
			// Vx += Vy
			x := (chip8.Opcode >> 8) & 0x000F
			y := (chip8.Opcode >> 4) & 0x000F

			if chip8.V[y] > (0xFF-chip8.V[x]) {
				chip8.V[0xF] = 1
			} else {
				chip8.V[0xF] = 0
			}

			chip8.V[x] += chip8.V[y]

			chip8.PC += 2
		case 0x0005:
			// VY is subtracted from VX. VF is set to 0 when there's a borrow, and 1 when there is not.
			// Vx -= Vy

			x := (chip8.Opcode >> 8) & 0x000F
			y := (chip8.Opcode >> 4) & 0x000F

			if chip8.V[x] > chip8.V[y] {
				chip8.V[0xF] = 1
			} else {
				chip8.V[0xF] = 0
			}

			chip8.V[x] -= chip8.V[y]

			chip8.PC += 2

		case 0x0006:
			// Stores the least significant bit of VX in VF and then shifts VX to the right by 1.
			// Vx>>=1

			x := (chip8.Opcode >> 8) & 0x000F

			if (chip8.V[x] & 0x01) == 0x01 {
				chip8.V[0xF] = 1
			} else {
				chip8.V[0xF] = 0
			}

			chip8.V[x] /= 2

			chip8.PC += 2

		case 0x0007:
			// Sets VX to VY minus VX. VF is set to 0 when there's a borrow, and 1 when there is not
			// Vx=Vy-Vx

			x := (chip8.Opcode >> 8) & 0x000F
			y := (chip8.Opcode >> 4) & 0x000F

			if chip8.V[y] > chip8.V[x] {
				chip8.V[0xF] = 1
			} else {
				chip8.V[0xF] = 0
			}

			chip8.V[x] -= chip8.V[y]

			chip8.PC += 2

		case 0x000E:
			// Stores the most significant bit of VX in VF and then shifts VX to the left by 1.
			// Vx<<=1

			x := (chip8.Opcode >> 8) & 0x000F

			if (chip8.V[x] & 0x80) == 0x80 {
				chip8.V[0xF] = 1
			} else {
				chip8.V[0xF] = 0
			}

			chip8.V[x] *= 2

			chip8.PC += 2

		default:
			panic(fmt.Sprintf("unknown opcode [0x8000]: 0x%X\n", chip8.Opcode))
		}
	case 0x9000: // 9XY0
		// Skips the next instruction if VX does not equal VY. (Usually the next instruction is a jump to skip a code block)
		// if(Vx!=Vy)
		x := (chip8.Opcode >> 8) & 0x000F
		y := (chip8.Opcode >> 4) & 0x000F

		if chip8.V[x] != chip8.V[y] {
			chip8.PC += 4
		} else {
			chip8.PC += 2
		}

	case 0xA000: // ANNN
		// Sets I to the address NNN.
		// I = NNN
		nnn := chip8.Opcode & 0x0FFF
		chip8.I = nnn

		chip8.PC += 2
	case 0xB000: // BNNN
		// Jumps to the address NNN plus V0.
		// PC=V0+NNN
		nnn := chip8.Opcode & 0x0FFF
		chip8.PC = uint16(chip8.V[0]) + uint16(nnn)
	case 0xC000: // CXNN
		// Sets VX to the result of a bitwise and operation on a random number (Typically: 0 to 255) and NN.
		// Vx=rand()&NN
		x := (chip8.Opcode >> 8) & 0x000F
		nn := byte(chip8.Opcode & 0x00FF)
		chip8.V[x] = RandomByte() & nn

		chip8.PC += 2
	case 0xD000: // DXYN
		// Draws a sprite at coordinate (VX, VY) that has a width of 8 pixels and a height of N+1 pixels
		// Each row of 8 pixels is read as bit-coded starting from memory location I; I value does not change after the execution of this instruction
		// As described above, VF is set to 1 if any screen pixels are flipped from set to unset when the sprite is drawn, and to 0 if that does not happen
		// draw(Vx,Vy,N)

		n := byte(chip8.Opcode & 0x000F)
		x := chip8.V[(chip8.Opcode>>8)&0x000F]
		y := chip8.V[(chip8.Opcode>>4)&0x000F]
		chip8.V[0xF] = 0

		for posY := 0; byte(posY) < n; posY++ {
			data := chip8.Memory[chip8.I+uint16(posY)]
			for posX := 0; posX < 8; posX++ {
				if (data & (0x80 >> posX)) != 0 {
					if chip8.Display[(int(x) + posX)][(int(y)+posY)] == 1 {
						chip8.V[0xF] = 1
					}
					chip8.Display[(int(x) + posX)][(int(y) + posY)] ^= 1
				}
			}
		}

		chip8.PC += 2

	case 0xE000:
		switch chip8.Opcode & 0x00FF {
		case 0x009E:
			// Skips the next instruction if the key stored in VX is pressed. (Usually the next instruction is a jump to skip a code block)
			// if(key()==Vx)
			x := (chip8.Opcode >> 8) & 0x000F
			if chip8.Keypad[x] == chip8.V[x] {
				chip8.PC += 2
			}

			chip8.PC += 2

		case 0x00A1:
			// Skips the next instruction if the key stored in VX is not pressed. (Usually the next instruction is a jump to skip a code block)
			// if(key()!=Vx)

			x := (chip8.Opcode >> 8) & 0x000F
			if chip8.Keypad[x] != chip8.V[x] {
				chip8.PC += 2
			}

			chip8.PC += 2

		default:
			panic(fmt.Sprintf("unknown opcode [0xE000]: 0x%X\n", chip8.Opcode))
		}
	case 0xF000:
		switch chip8.Opcode & 0x00FF {
		case 0x0007:
			// Sets VX to the value of the delay timer.
			// Vx = get_delay()

			x := (chip8.Opcode>>8)&0x000F
			chip8.V[x] = chip8.DelayTimer

			chip8.PC += 2

		case 0x000A:
			// A key press is awaited, and then stored in VX. (Blocking Operation. All instruction halted until next key event)
			// Vx = get_key()

			x := (chip8.Opcode >> 8) & 0x000F

			chip8.V[x] = chip8.Keypad[x]

		case 0x0015:
			// Sets the delay timer to VX.
			// delay_timer(Vx)
			x := (chip8.Opcode>>8)&0x000F
			chip8.DelayTimer = chip8.V[x]

			chip8.PC += 2
		case 0x0018:
			// Sets the sound timer to VX.
			// sound_timer(Vx)
			x := (chip8.Opcode>>8)&0x000F
			chip8.SoundTimer = chip8.V[x]

			chip8.PC += 2
		case 0x001E:
			// Adds VX to I. VF is not affected
			// I +=Vx
			x := chip8.V[(chip8.Opcode>>8)&0x000F]
			chip8.I += uint16(x)

			chip8.PC += 2
		case 0x0029:
			// Sets I to the location of the sprite for the character in VX. Characters 0-F (in hexadecimal) are represented by a 4x5 font.
			// I=sprite_addr[Vx]
			x := chip8.Opcode >> 8 & 0x000F
			chip8.I = uint16(chip8.V[x]) * uint16(0x5)

			chip8.PC += 2
		case 0x0033:
			// Stores the binary-coded decimal representation of VX, with the most
			// significant of three digits at the address in I, the middle digit at I plus 1, and the
			// least significant digit at I plus 2. (In other words, take the decimal
			// representation of VX, place the hundreds digit in memory at location in I, the
			// tens digit at location I+1, and the ones digit at location I+2.)
			// set_BCD(Vx)
			// *(I+0)=BCD(3)
			// *(I+1)=BCD(2)
			// *(I+2)=BCD(1)
			x := chip8.Opcode >> 8 & 0x000F

			chip8.Memory[chip8.I] = chip8.V[x] / 100
			chip8.Memory[chip8.I+1] = (chip8.V[x] / 10) % 10
			chip8.Memory[chip8.I+2] = chip8.V[x] % 10

			chip8.PC += 2
		case 0x0055:
			// Stores V0 to VX (including VX) in memory starting at address I.
			// The offset from I is increased by 1 for each value written, but I itself is left unmodified
			// reg_dump(Vx,&I)
			x := chip8.Opcode >> 8 & 0x000F

			for i := uint16(0); i <= x; i++ {
				chip8.Memory[i+chip8.I] = chip8.V[i]
			}

			chip8.PC += 2
		case 0x0065:
			// Fills V0 to VX (including VX) with values from memory starting at address I
			// The offset from I is increased by 1 for each value written, but I itself is left unmodified.
			// reg_load(Vx,&I)
			x := chip8.Opcode >> 8 & 0x000F

			for i := 0; uint16(i) <= x; i++ {
				chip8.V[i] = chip8.Memory[i+int(chip8.I)]
			}

			chip8.PC += 2
		default:
			panic(fmt.Sprintf("unknown opcode [0xF000]: 0x%X\n", chip8.Opcode))
		}

	default:
		panic(fmt.Sprintf("unknown opcode: 0x%X\n", chip8.Opcode))
	}

	if chip8.DelayTimer > 0 {
		chip8.DelayTimer--
	}

	if chip8.SoundTimer > 0 {
		chip8.SoundTimer--
	} 

	return nil
}

func (g *Chip8) Draw(screen *ebiten.Image) {

	var pixelImage = ebiten.NewImage(64, 32)

	for row := 0; row < len(g.Display); row++ {
		for col := 0; col < len(g.Display[row]); col++ {
			var currentColor color.Color

			option := &ebiten.DrawImageOptions{}
			option.GeoM.Translate(float64(row), float64(col))
			option.GeoM.Scale(float64(screen.Bounds().Dx())/64, float64(screen.Bounds().Dy())/32)
			if g.Display[row][col] == 1 {
				currentColor = color.White
			} else {
				currentColor = color.Black
			}

			pixelImage.Fill(currentColor)
			screen.DrawImage(pixelImage, option)
		}
	}
}

func (g *Chip8) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func RandomByte() byte {
	rand.Seed(time.Now().UnixNano())

	return byte(rand.Intn(0xFF))
}

func ReadROM(filename string) []byte {

	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	return bytes
}

func (c *Chip8) Init() {
	c.Opcode = 0
	c.PC = 0x200
	c.I = 0
	c.SP = 0
	c.DelayTimer = 0
	c.SoundTimer = 0

	// Clear Register
	for x := 0; x < 16; x++ {
		c.V[x] = 0
	}

	// clears the display
	c.Display = [64][32]byte{}

	// clear stack
	for i := 0; i < 16; i++ {
		c.Stack[i] = 0
	}

	// clear memory
	for i := 0; i < 4096; i++ {
		c.Memory[i] = 0
	}

	// load font set
	for i := 0; i < 80; i++ {
		c.Memory[i] = chip8FontSet[i]
	}
}

func (c Chip8) Disassemble(program byte, pc uint16) {

	c.Opcode = (uint16(c.Memory[pc]) << 8) | uint16(c.Memory[pc+1])

	fmt.Printf("%-4X %4X\t", c.PC, c.Opcode)
	switch c.Opcode & 0xF000 {
	case 0x0000:
		switch c.Opcode & 0x00F0 {
		case 0x00E0:
			fmt.Println("CLS")
		case 0x00EE:
			fmt.Println("RET")
		default:
			nnn := c.Opcode & 0x0FFF
			fmt.Printf("SYS %X\n", nnn)
		}
	case 0x1000:
		c.PC = c.Opcode & 0x0FFF
		fmt.Printf("JP %X\n", c.PC)
	case 0x2000:
		nnn := c.Opcode & 0x0FFF
		fmt.Printf("CALL %X\n", nnn)
	case 0x3000:
		nn := byte(c.Opcode & 0x00FF)
		x := (c.Opcode >> 8) & 0x000F
		fmt.Printf("SE V%X, %X\n", x, nn)
	case 0x4000:
		x := (c.Opcode >> 8) & 0x000F
		nn := byte(c.Opcode & 0x00FF)
		fmt.Printf("SNE V%X, %X\n", x, nn)
	case 0x5000:
		x := (c.Opcode >> 8) & 0x000F
		y := (c.Opcode >> 4) & 0x000F
		fmt.Printf("SE V%X, V%X\n", x, y)
	case 0x6000:
		v := c.Opcode >> 8 & 0x000F
		nn := byte(c.Opcode & 0x00FF)
		fmt.Printf("LD V%X, %X\n", v, nn)
	case 0x7000:
		nn := byte(c.Opcode & 0x00FF)
		x := (c.Opcode >> 8) & 0x000F
		fmt.Printf("ADD V%X, %X\n", x, nn)
	case 0x8000:
		switch c.Opcode & 0x000F {
		case 0x0:
			x := (c.Opcode >> 8) & 0x000F
			y := (c.Opcode >> 4) & 0x000F
			fmt.Printf("LD V%X, V%X\n", x, y)
		case 0x1:
			x := (c.Opcode >> 8) & 0x000F
			y := (c.Opcode >> 4) & 0x000F
			fmt.Printf("OR V%X, V%X\n", x, y)
		case 0x2:
			x := (c.Opcode >> 8) & 0x000F
			y := (c.Opcode >> 4) & 0x000F
			fmt.Printf("AND V%X, V%X\n", x, y)
		case 0x3:
			x := (c.Opcode >> 8) & 0x000F
			y := (c.Opcode >> 4) & 0x000F
			fmt.Printf("XOR V%X, V%X\n", x, y)
		case 0x4:
			x := (c.Opcode >> 8) & 0x000F
			y := (c.Opcode >> 4) & 0x000F
			fmt.Printf("ADD V%X, V%X\n", x, y)
		case 0x5:
			x := (c.Opcode >> 8) & 0x000F
			y := (c.Opcode >> 4) & 0x000F
			fmt.Printf("SUB V%X, V%X\n", x, y)
		case 0x6:
			x := (c.Opcode >> 8) & 0x000F
			fmt.Printf("SHR V%X\n", x)
		case 0x7:
			x := (c.Opcode >> 8) & 0x000F
			y := (c.Opcode >> 4) & 0x000F
			fmt.Printf("SUBN V%X, V%X\n", x, y)
		case 0xE:
			x := (c.Opcode >> 8) & 0x000F
			fmt.Printf("SHL V%X\n", x)
		default:
			fmt.Printf("UNKN %X\n", c.Opcode)
		}
	case 0x9000:
		x := (c.Opcode >> 8) & 0x000F
		y := (c.Opcode >> 4) & 0x000F
		fmt.Printf("SNE V%X, V%X\n", x, y)
	case 0xA000:
		nnn := c.Opcode & 0x0FFF
		fmt.Printf("LD I, %X\n", nnn)
	case 0xB000:
		nnn := c.Opcode & 0x0FFF
		fmt.Printf("JP V0, %X\n", nnn)
	case 0xC000:
		nn := byte(c.Opcode & 0x00FF)
		x := (c.Opcode >> 8) & 0x000F
		fmt.Printf("RND V%X, %X\n", x, nn)
	case 0xD000:
		n := c.Opcode & 0x000F
		x := (c.Opcode >> 8) & 0x000F
		y := (c.Opcode >> 4) & 0x000F
		fmt.Printf("DRW V%X, V%X, %X\n", x, y, n)
	case 0xE000:
		switch c.Opcode & 0x00FF {
		case 0x9E:
			x := (c.Opcode >> 8) & 0x000F
			fmt.Printf("SKP V%X\n", x)
		case 0xA1:
			x := (c.Opcode >> 8) & 0x000F
			fmt.Printf("SKNP V%X\n", x)
		default:
			fmt.Printf("UNKN %X\n", c.Opcode)
		}
	case 0xF000:
		switch c.Opcode & 0x00FF {
		case 0x07:
			x := (c.Opcode >> 8) & 0x000F
			fmt.Printf("LD V%X, DT\n", x)
		case 0x0A:
			x := (c.Opcode >> 8) & 0x000F
			fmt.Printf("LD V%X, K\n", x)
		case 0x15:
			x := (c.Opcode >> 8) & 0x000F
			fmt.Printf("LD DT, V%X\n", x)
		case 0x18:
			x := (c.Opcode >> 8) & 0x000F
			fmt.Printf("LD ST, V%X\n", x)
		case 0x1E:
			x := (c.Opcode >> 8) & 0x000F
			fmt.Printf("ADD I, V%X\n", x)
		case 0x29:
			x := (c.Opcode >> 8) & 0x000F
			fmt.Printf("LD F, V%X\n", x)
		case 0x33:
			x := (c.Opcode >> 8) & 0x000F
			fmt.Printf("LD B, V%X\n", x)
		case 0x55:
			x := (c.Opcode >> 8) & 0x000F
			fmt.Printf("LD [I], V%X\n", x)
		case 0x65:
			x := (c.Opcode >> 8) & 0x000F
			fmt.Printf("LD V%X, [I]\n", x)
		default:
			fmt.Printf("UNKN %X\n", c.Opcode)
		}
	default:
		fmt.Printf("UNKN %X\n", c.Opcode)
	}
}
