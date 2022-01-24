// Package main provides primitive CPU functions i.e
// fetching - decoding - executing a program.
package main

import (
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

// A CPU represents the internal of a cpu for the chip-8.
type CPU struct {
	// current opcode in program
	opcode uint16

	// 4KB of memory
	ram Memory

	// qwerty keyboard
	keypad KeyPad

	// executes the current opcode
	pc uint16

	// pointer for stack register
	sp uint16

	// 16 8-bit (one byte) general-purpose variable registers numbered
	v [16]byte

	// register for storing memory addresses
	i uint16

	// store addresses that the interpreter should return from a subroutine
	stack [16]uint16

	// 64x32-pixel monochrome display
	display Display

	// active whenever the delay timer register (DT) is non-zero
	// This timer does nothing more than subtract 1 from the value of DT at a rate of 60Hz. When DT reaches 0, it deactivates.
	delayTimer byte

	// active whenever the sound timer register (ST) is non-zero
	// This timer also decrements at a rate of 60Hz, however, as long as ST's value is greater than zero, the Chip-8 buzzer will sound
	soundTimer byte
}

// NewCPU
func NewCPU() *CPU {

	cpu := &CPU{
		ram:        Memory{},
		pc:         0x200,
		sp:         0,
		v:          [16]byte{},
		i:          0,
		stack:      [16]uint16{},
		delayTimer: 0,
		soundTimer: 0,
	}

	// load character bytes into memory
	for i := 0; i < len(FontSet); i++ {
		cpu.ram.RAM[i] = FontSet[i]
	}

	return cpu
}

// EmulateCycle
func (cpu *CPU) EmulateCycle() {
	cpu.fetchOpcode()
	cpu.executeOpcode()
	cpu.updateTimers()
}

// updateTimers decrements timers unless value is less than 0
func (cpu *CPU) updateTimers() {

	if cpu.delayTimer > 0 {
		cpu.delayTimer--
	}

	if cpu.soundTimer > 0 {
		cpu.soundTimer--
	}
}

// FetchOpcode fetches the current instruction based on the pc.
//
func (cpu *CPU) fetchOpcode() {
	cpu.opcode = (uint16(cpu.ram.RAM[cpu.pc]) << 8) | uint16(cpu.ram.RAM[cpu.pc+1])
}

// ExecuteOpcode decodes the opcode and execute its function.
func (cpu *CPU) executeOpcode() {

	var (
		// A 4-bit value, the lower 4 bits of the high byte of the instruction
		x = (cpu.opcode >> 8) & 0x000F

		// A 4-bit value, the upper 4 bits of the low byte of the instruction
		y = (cpu.opcode >> 4) & 0x000F

		// right-most nibble of byte
		nibble = byte(cpu.opcode & 0x000F)

		// A 4-bit value, the lowest 4 bits of the instruction
		nn = cpu.opcode & 0x00FF

		// A 12-bit value, the lowest 12 bits of the instruction
		nnn = cpu.opcode & 0x0FFF
	)

	// TODO: abstract key presses into methods
	cpu.keypad.Keys[x] = 0

	if ebiten.IsKeyPressed(ebiten.KeyDigit1) {
		cpu.keypad.Keys[x] = 0x0
	}

	if ebiten.IsKeyPressed(ebiten.KeyDigit2) {
		cpu.keypad.Keys[x] = 0x1
	}

	if ebiten.IsKeyPressed(ebiten.KeyDigit3) {
		cpu.keypad.Keys[x] = 0x2
	}

	if ebiten.IsKeyPressed(ebiten.KeyDigit4) {
		cpu.keypad.Keys[x] = 0x3
	}

	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		cpu.keypad.Keys[x] = 0x4
	}

	if ebiten.IsKeyPressed(ebiten.KeyW) {
		cpu.keypad.Keys[x] = 0x5
	}

	if ebiten.IsKeyPressed(ebiten.KeyE) {
		cpu.keypad.Keys[x] = 0x6
	}

	if ebiten.IsKeyPressed(ebiten.KeyR) {
		cpu.keypad.Keys[x] = 0x7
	}

	if ebiten.IsKeyPressed(ebiten.KeyA) {
		cpu.keypad.Keys[x] = 0x8
	}

	if ebiten.IsKeyPressed(ebiten.KeyS) {
		cpu.keypad.Keys[x] = 0x9
	}

	if ebiten.IsKeyPressed(ebiten.KeyD) {
		cpu.keypad.Keys[x] = 0xA
	}

	if ebiten.IsKeyPressed(ebiten.KeyF) {
		cpu.keypad.Keys[x] = 0xB
	}

	if ebiten.IsKeyPressed(ebiten.KeyZ) {
		cpu.keypad.Keys[x] = 0xC
	}

	if ebiten.IsKeyPressed(ebiten.KeyX) {
		cpu.keypad.Keys[x] = 0xD
	}

	if ebiten.IsKeyPressed(ebiten.KeyC) {
		cpu.keypad.Keys[x] = 0xE
	}

	if ebiten.IsKeyPressed(ebiten.KeyV) {
		cpu.keypad.Keys[x] = 0xF
	}

	switch cpu.opcode & 0xF000 {
	case 0x0000:
		switch cpu.opcode & 0x0FFF {
		case 0x00E0:
			cpu.cls()
		case 0x00EE:
			cpu.ret()
		default:
			cpu.unknownOpcode()
		}

	case 0x1000:
		cpu.jump(nnn)
	case 0x2000:
		cpu.call(nnn)
	case 0x3000:
		cpu.skipIf(x, nn)
	case 0x4000:
		cpu.skipIfNot(x, nn)
	case 0x5000:
		cpu.skipIfXY(x, y)
	case 0x6000:
		cpu.loadX(x, nn)
	case 0x7000:
		cpu.addX(x, nn)

	case 0x8000:
		switch cpu.opcode & 0x000F {
		case 0x0000:
			cpu.loadXY(x, y)
		case 0x0001:
			cpu.or(x, y)
		case 0x0002:
			cpu.and(x, y)
		case 0x0003:
			cpu.xor(x, y)
		case 0x0004:
			cpu.addXY(x, y)
		case 0x0005:
			cpu.subXY(x, y)
		case 0x0006:
			cpu.shr(x)
		case 0x0007:
			cpu.subYX(x, y)
		case 0x000E:
			cpu.shl(x)
		default:
			cpu.unknownOpcode()
		}

	case 0x9000:
		cpu.skipIfNotXY(x, y)
	case 0xA000:
		cpu.loadI(nnn)
	case 0xB000:
		cpu.jumpV0(nnn, uint16(cpu.v[0]))
	case 0xC000:
		cpu.loadRnd(x, nn)
	case 0xD000:
		cpu.draw(uint16(nibble), x, y)

	case 0xE000:
		switch cpu.opcode & 0x00FF {
		case 0x009E:
			cpu.skipIfPressed(x)
		case 0x00A1:
			cpu.skipIfNotPressed(x)
		default:
			cpu.unknownOpcode()
		}

	case 0xF000:
		switch cpu.opcode & 0x00FF {
		case 0x0007:
			cpu.loadXDT(x)
		case 0x000A:
			cpu.loadVK(x)
		case 0x0015:
			cpu.loadDTX(x)
		case 0x0018:
			cpu.loadSTX(x)
		case 0x001E:
			cpu.addIX(x)
		case 0x0029:
			cpu.loadF(x)
		case 0x0033:
			cpu.loadBX(x)
		case 0x0055:
			cpu.loadRegIX(x)
		case 0x0065:
			cpu.loadRegX(x)
		default:
			cpu.unknownOpcode()
		}

	default:
		cpu.unknownOpcode()
	}
}

// unknownOpcode interrupts cpu cycle when opcode is unidentifiable
func (cpu *CPU) unknownOpcode() {
	log.Fatalf("unknown opcode [%#X]: 0x%X\n", cpu.opcode&0xF000, cpu.opcode)
}

// cls 00E0 - CLS
// Clear the display.
func (cpu *CPU) cls() {
	cpu.display.Clear()
	cpu.pc += 2
}

// ret return from a subroutine.
// The interpreter sets the program counter to the address at the top of the stack, then subtracts 1 from the stack pointer.
func (cpu *CPU) ret() {
	cpu.sp--
	cpu.pc = cpu.stack[cpu.sp]
	cpu.pc += 2
}

// call subroutine at nnn.
// The interpreter increments the stack pointer, then puts the current PC on the top of the stack. The PC is then set to nnn.
func (cpu *CPU) call(address uint16) {
	cpu.stack[cpu.sp] = cpu.pc
	cpu.sp++
	cpu.pc = address
}

// jump to location nnn.
// The interpreter sets the program counter to nnn.
func (cpu *CPU) jump(address uint16) {
	cpu.pc = address
}

// jumpV0 Bnnn - JP V0, addr
// Jump to location nnn + V0.
// The program counter is set to nnn plus the value of V0.
func (cpu *CPU) jumpV0(address uint16, v0 uint16) {
	cpu.pc = address + v0
}

// loadRnd Cxkk - RND Vx, byte
// Set Vx = random byte AND kk.
// The interpreter generates a random number from 0 to 255, which is then ANDed with the value kk. The results are stored in Vx. See instruction 8xy2 for more information on AND.
func (cpu *CPU) loadRnd(x uint16, nn uint16) {
	cpu.v[x] = byte(rand.Intn(255)) & byte(nn)
	cpu.pc += 2
}

// loadI Annn - LD I, addr
// Set I = nnn.
// The value of register I is set to nnn.
func (cpu *CPU) loadI(address uint16) {
	cpu.i = address
	cpu.pc += 2
}

// skipIfNotXY 9xy0 - SNE Vx, Vy
// Skip next instruction if Vx != Vy.
// The values of Vx and Vy are compared, and if they are not equal, the program counter is increased by 2.
func (cpu *CPU) skipIfNotXY(x uint16, y uint16) {
	if cpu.v[x] != cpu.v[y] {
		cpu.pc += 4
	} else {
		cpu.pc += 2
	}
}

// draw Dxyn - DRW Vx, Vy, nibble (n)
// Display n-byte sprite starting at memory location I at (Vx, Vy), set VF = collision.
func (cpu *CPU) draw(n uint16, x uint16, y uint16) {
	x = uint16(cpu.v[x])
	y = uint16(cpu.v[y])
	cpu.v[0xF] = 0

	for row := uint16(0); row < n; row++ {
		pixel := cpu.ram.RAM[cpu.i+row]
		for col := uint16(0); col < 8; col++ {
			var (
				xIdx = x + col
				yIdx = y + row
			)

			if pixel&(0x80>>col) != 0 {
				if cpu.display.gfx[xIdx][yIdx] == 1 {
					cpu.v[0xF] = 1
				}
				cpu.display.gfx[xIdx][yIdx] ^= 1
			}

		}
	}
	cpu.pc += 2
}

// skipIf Skip next instruction if Vx = nn.
// The interpreter compares register Vx to nn, and if they are equal, increments the program counter by 2.
func (cpu *CPU) skipIf(x uint16, nn uint16) {
	if uint16(cpu.v[x]) == nn {
		cpu.pc += 4
	} else {
		cpu.pc += 2
	}
}

// skipIfNot Skip next instruction if Vx != nn.
// The interpreter compares register Vx to nn, and if they are not equal, increments the program counter by 2.
func (cpu *CPU) skipIfNot(x uint16, nn uint16) {
	if uint16(cpu.v[x]) != nn {
		cpu.pc += 4
	} else {
		cpu.pc += 2
	}
}

// skipIfXY 5xy0 - SE Vx, Vy
// Skip next instruction if Vx = Vy.
// The interpreter compares register Vx to register Vy, and if they are equal, increments the program counter by 2.
func (cpu *CPU) skipIfXY(x uint16, y uint16) {
	if cpu.v[x] == cpu.v[y] {
		cpu.pc += 4
	} else {
		cpu.pc += 2
	}
}

// loadX 6xkk - LD Vx, byte
// Set Vx = kk.
// The interpreter puts the value kk into register Vx.
func (cpu *CPU) loadX(x uint16, nn uint16) {
	cpu.v[x] = byte(nn)
	cpu.pc += 2
}

// addX 7xkk - ADD Vx, byte
// Set Vx = Vx + kk.
// Adds the value kk to the value of register Vx, then stores the result in Vx.
func (cpu *CPU) addX(x uint16, nn uint16) {
	cpu.v[x] += byte(nn)
	cpu.pc += 2
}

// and 8xy2 - AND Vx, Vy
// Set Vx = Vx AND Vy.
// Performs a bitwise AND on the values of Vx and Vy, then stores the result in Vx. A bitwise AND compares the corrseponding bits from two values, and if both bits are 1, then the same bit in the result is also 1. Otherwise, it is 0.
func (cpu *CPU) and(x uint16, y uint16) {
	cpu.v[x] &= cpu.v[y]
	cpu.pc += 2
}

// xor 8xy3 - XOR Vx, Vy
// Set Vx = Vx XOR Vy.
// Performs a bitwise exclusive OR on the values of Vx and Vy, then stores the result in Vx. An exclusive OR compares the corrseponding bits from two values, and if the bits are not both the same, then the corresponding bit in the result is set to 1. Otherwise, it is 0.
func (cpu *CPU) xor(x uint16, y uint16) {
	cpu.v[x] ^= cpu.v[y]
	cpu.pc += 2
}

// or 8xy1 - OR Vx, Vy
// Set Vx = Vx OR Vy.
// Performs a bitwise OR on the values of Vx and Vy, then stores the result in Vx. A bitwise OR compares the corrseponding bits from two values, and if either bit is 1, then the same bit in the result is also 1. Otherwise, it is 0.
func (cpu *CPU) or(x uint16, y uint16) {
	cpu.v[x] |= cpu.v[y]
	cpu.pc += 2
}

// loadXY 8xy0 - LD Vx, Vy
// Set Vx = Vy.
// Stores the value of register Vy in register Vx.
func (cpu *CPU) loadXY(x uint16, y uint16) {
	cpu.v[x] = cpu.v[y]
	cpu.pc += 2
}

// addXY Set Vx = Vx + Vy, set VF = carry.
// The values of Vx and Vy are added together. If the result is greater than 8 bits (i.e., > 255,) VF is set to 1, otherwise 0. Only the lowest 8 bits of the result are kept, and stored in Vx.
func (cpu *CPU) addXY(x uint16, y uint16) {

	if (cpu.v[x] + cpu.v[y]) > 255 {
		cpu.v[0xF] = 1
	} else {
		cpu.v[0xF] = 0
	}
	cpu.v[x] += cpu.v[y]
	cpu.pc += 2
}

// subXY 8xy5 - SUB Vx, Vy
// Set Vx = Vx - Vy, set VF = NOT borrow.
// If Vx > Vy, then VF is set to 1, otherwise 0. Then Vy is subtracted from Vx, and the results stored in Vx.
func (cpu *CPU) subXY(x uint16, y uint16) {
	if cpu.v[x] > cpu.v[y] {
		cpu.v[0xF] = 1
	} else {
		cpu.v[0xF] = 0
	}
	cpu.v[x] -= cpu.v[y]
	cpu.pc += 2
}

// subYX 8xy7 - SUBN Vx, Vy
// Set Vx = Vy - Vx, set VF = NOT borrow.
// If Vy > Vx, then VF is set to 1, otherwise 0. Then Vx is subtracted from Vy, and the results stored in Vx.
func (cpu *CPU) subYX(x uint16, y uint16) {
	if cpu.v[y] > cpu.v[x] {
		cpu.v[0xF] = 1
	} else {
		cpu.v[0xF] = 0
	}
	cpu.v[x] -= cpu.v[y]
	cpu.pc += 2
}

// shr 8xy6 - SHR Vx {, Vy}
// Set Vx = Vx SHR 1.
// If the least-significant bit of Vx is 1, then VF is set to 1, otherwise 0. Then Vx is divided by 2.
func (cpu *CPU) shr(x uint16) {
	if (cpu.v[x] & 0x1) == 1 {
		cpu.v[0xF] = 1
	} else {
		cpu.v[0xF] = 0
	}
	cpu.v[x] /= 2
	cpu.pc += 2
}

// shl 8xyE - SHL Vx {, Vy}
// Set Vx = Vx SHL 1.
// If the most-significant bit of Vx is 1, then VF is set to 1, otherwise to 0. Then Vx is multiplied by 2.
func (cpu *CPU) shl(x uint16) {
	if (cpu.v[x] & 0x80) == 1 {
		cpu.v[0xF] = 1
	} else {
		cpu.v[0xF] = 0
	}
	cpu.v[x] *= 2
	cpu.pc += 2
}

// skipIfPressed Ex9E - SKP Vx
// Skip next instruction if key with the value of Vx is pressed.
// Checks the keyboard, and if the key corresponding to the value of Vx is currently in the down position, PC is increased by 2.
func (cpu *CPU) skipIfPressed(x uint16) {
	if cpu.keypad.Keys[byte(x)] == cpu.v[x] {
		cpu.pc += 2
	}
	cpu.pc += 2
}

// skipIfNotPressed ExA1 - SKNP Vx
// Skip next instruction if key with the value of Vx is not pressed.
// Checks the keyboard, and if the key corresponding to the value of Vx is currently in the up position, PC is increased by 2.
func (cpu *CPU) skipIfNotPressed(x uint16) {
	if cpu.keypad.Keys[byte(x)] != cpu.v[x] {
		cpu.pc += 2
	}
	cpu.pc += 2
}

// loadXDT Fx07 - LD Vx, DT
// Set Vx = delay timer value.
//The value of DT is placed into Vx.
func (cpu *CPU) loadXDT(x uint16) {
	cpu.v[x] = cpu.delayTimer
	cpu.pc += 2
}

// loadDTX Fx15 - LD DT, Vx
// Set delay timer = Vx.
// DT is set equal to the value of Vx.
func (cpu *CPU) loadDTX(x uint16) {
	cpu.delayTimer = byte(cpu.v[x])
	cpu.pc += 2
}

// loadSTX Fx18 - LD ST, Vx
// Set sound timer = Vx.
// ST is set equal to the value of Vx.
func (cpu *CPU) loadSTX(x uint16) {
	cpu.soundTimer = byte(cpu.v[x])
	cpu.pc += 2
}

// addIX Fx1E - ADD I, Vx
// Set I = I + Vx.
// The values of I and Vx are added, and the results are stored in I.
func (cpu *CPU) addIX(x uint16) {
	cpu.i += uint16(cpu.v[x])
	cpu.pc += 2
}

// loadF Fx29 - LD F, Vx
// Set I = location of sprite for digit Vx.
// The value of I is set to the location for the hexadecimal sprite corresponding to the value of Vx. See section 2.4, Display, for more information on the Chip-8 hexadecimal font.
func (cpu *CPU) loadF(x uint16) {
	cpu.i = uint16(cpu.v[x]) * 5
	cpu.pc += 2
}

// loadBX Fx33 - LD B, Vx
// Store BCD representation of Vx in memory locations I, I+1, and I+2.
// The interpreter takes the decimal value of Vx, and places the hundreds digit in memory at location in I, the tens digit at location I+1, and the ones digit at location I+2.
func (cpu *CPU) loadBX(x uint16) {
	cpu.ram.RAM[cpu.i] = byte(cpu.v[x] / 100)          // get hundreds digit
	cpu.ram.RAM[cpu.i+1] = byte((cpu.v[x] % 100) / 10) // get tens digit
	cpu.ram.RAM[cpu.i+2] = byte(cpu.v[x] % 10)         // get ones digit
	cpu.pc += 2
}

// loadRegIX Fx55 - LD [I], Vx
// Store registers V0 through Vx in memory starting at location I.
// The interpreter copies the values of registers V0 through Vx into memory, starting at the address in I.
func (cpu *CPU) loadRegIX(x uint16) {
	for i := uint16(0); i <= x; i++ {
		cpu.ram.RAM[cpu.i+i] = cpu.v[i]
	}
	cpu.pc += 2
}

// loadRegx Fx65 - LD Vx, [I]
// Read registers V0 through Vx from memory starting at location I.
// The interpreter reads values from memory starting at location I into registers V0 through Vx.
func (cpu *CPU) loadRegX(x uint16) {
	for i := uint16(0); i <= x; i++ {
		cpu.v[i] = cpu.ram.RAM[cpu.i+i]
	}
	cpu.pc += 2
}

// loadVK Fx0A - LD Vx, K
// Wait for a key press, store the value of the key in Vx.
// All execution stops until a key is pressed, then the value of that key is stored in Vx.
func (cpu *CPU) loadVK(x uint16) {
	cpu.v[x] = cpu.keypad.Keys[x]
	cpu.pc += 2
}
