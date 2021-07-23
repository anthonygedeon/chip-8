// Package main provides primitive CPU functions i.e
// fetching - decoding - executing a program.
package main

import (
	"log"
	"math/rand"
)

// A CPU represents the internal of a cpu for the chip-8.
type CPU struct {
	// current opcode in program
	opcode uint16

	// 4KB of memory
	ram Memory

	// executes the current opcode
	pc uint16

	// pointer for stack register
	sp uint16

	// general purpose register
	v [16]uint16

	// register for storing memory addresses
	i uint16

	// store addresses that the interpreter should return from a subroutine
	stack [16]uint16

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
	return &CPU{
		ram:        Memory{},
		pc:         0x200,
		sp:         0,
		v:          [16]uint16{},
		i:          0,
		stack:      [16]uint16{},
		delayTimer: 0,
		soundTimer: 0,
	}
}

// FetchOpcode fetches the current instruction based on the pc.
//
func (cpu *CPU) FetchOpcode() {
	cpu.opcode = (uint16(cpu.ram.RAM[cpu.pc]) << 8) | uint16(cpu.ram.RAM[cpu.pc+1])
}

// ExecuteOpcode decodes the opcode and execute its function.
func (cpu *CPU) ExecuteOpcode() {

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
			cpu.addXY()
		case 0x0005:
			cpu.subXY()
		case 0x0006:
			cpu.shr()
		case 0x0007:
			cpu.subYX()
		case 0x000E:
			cpu.shl()
		default:
			cpu.unknownOpcode()
		}

	case 0x9000:
		cpu.skipIfNotXY(x, y)
	case 0xA000:
		cpu.loadI(nnn)
	case 0xB000:
		cpu.jumpV0(nnn, cpu.v[0])
	case 0xC000:
		cpu.loadRnd(x, nn)
	case 0xD000:
		cpu.draw(uint16(nibble), x, y)

	case 0xE000:
		switch cpu.opcode & 0x00F0 {
		case 0x009E:
			cpu.skipIfPressed()
		case 0x00A1:
			cpu.skipIfNotPressed()
		default:
			cpu.unknownOpcode()
		}

	case 0xF000:
		switch cpu.opcode & 0x00F0 {
		case 0x0007:
			cpu.loadXDT(x)
		case 0x000A:
			cpu.loadVK()
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
			cpu.loadRegIX()
		case 0x0065:
			cpu.loadRegX()
		default:
			cpu.unknownOpcode()
		}

	default:
		cpu.unknownOpcode()
	}

}

// unknownOpcode
func (cpu *CPU) unknownOpcode() {
	log.Fatalf("unknown opcode [0x%X000]: 0x%X\n", cpu.opcode&0xF000, cpu.opcode)
}

// cls
func (cpu *CPU) cls() {
	cpu.display.gfx = [64][32]byte{}
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

// jumpV0
func (cpu *CPU) jumpV0(address uint16, v0 uint16) {
	cpu.pc = address + v0
}

// loadRnd
func (cpu *CPU) loadRnd(x uint16, nn uint16) {
	cpu.v[x] = uint16(rand.Intn(255)) & nn
}

// loadI
func (cpu *CPU) loadI(address uint16) {
	cpu.i = address
	cpu.pc += 2
}

// skipIfNotXY
func (cpu *CPU) skipIfNotXY(x uint16, y uint16) {
	if cpu.v[x] != cpu.v[y] {
		cpu.pc += 4
	} else {
		cpu.pc += 2
	}
}

// draw
func (cpu *CPU) draw(n uint16, x uint16, y uint16) {
	x = cpu.v[x]
	y = cpu.v[y]
	cpu.v[0xF] = 0

	for posY := 0; uint16(posY) < n; posY++ {
		data := cpu.ram.RAM[cpu.i+uint16(posY)]
		for posX := 0; posX < 8; posX++ {
			if (data & (0x80 >> posX)) != 0 {
				if cpu.display.gfx[(int(x) + posX)][(int(y)+posY)] == 1 {
					cpu.v[0xF] = 1
				}
				cpu.display.gfx[(int(x) + posX)][(int(y) + posY)] ^= 1
			}
		}
	}

	cpu.pc += 2
}

// skipIf
func (cpu *CPU) skipIf(x uint16, nn uint16) {
	if cpu.v[x] == nn {
		cpu.pc += 4
	} else {
		cpu.pc += 2
	}
}

// skipIfNot
func (cpu *CPU) skipIfNot(x uint16, nn uint16) {
	if cpu.v[x] != nn {
		cpu.pc += 4
	} else {
		cpu.pc += 2
	}
}

// skipIfXY
func (cpu *CPU) skipIfXY(x uint16, y uint16) {
	if cpu.v[x] == cpu.v[y] {
		cpu.pc += 4
	} else {
		cpu.pc += 2
	}
}

// loadX
func (cpu *CPU) loadX(x uint16, nn uint16) {
	cpu.v[x] = nn
	cpu.pc += 2
}

// addX
func (cpu *CPU) addX(x uint16, nn uint16) {
	cpu.v[x] += nn
	cpu.pc += 2
}

// and
func (cpu *CPU) and(x uint16, y uint16) {
	cpu.v[x] &= cpu.v[y]
	cpu.pc += 2
}

// xor
func (cpu *CPU) xor(x uint16, y uint16) {
	cpu.v[x] ^= cpu.v[y]
	cpu.pc += 2
}

// or
func (cpu *CPU) or(x uint16, y uint16) {
	cpu.v[x] |= cpu.v[y]
	cpu.pc += 2
}

// loadXY
func (cpu *CPU) loadXY(x uint16, y uint16) {
	cpu.v[x] = cpu.v[y]
	cpu.pc += 2
}

// addXY
func (cpu *CPU) addXY() {
	
	cpu.pc += 2
}

// subXY
func (cpu *CPU) subXY() {
	cpu.pc += 2
}

// subYX
func (cpu *CPU) subYX() {
	cpu.pc += 2
}

// shr
func (cpu *CPU) shr() {
	cpu.pc += 2
}	

// addYX
func (cpu *CPU) addYX() {
}

// shl
func (cpu *CPU) shl() {
	cpu.pc += 2
}

// skipIfPressed
func (cpu *CPU) skipIfPressed() {

}

// skipIfNotPressed
func (cpu *CPU) skipIfNotPressed() {

}

// loadXDT
func (cpu *CPU) loadXDT(x uint16) {
	cpu.v[x] = uint16(cpu.delayTimer)
	cpu.pc += 2
}

// loadDTX
func (cpu *CPU) loadDTX(x uint16) {
	cpu.delayTimer = byte(cpu.v[x])
	cpu.pc += 2
}

// loadSTX
func (cpu *CPU) loadSTX(x uint16) {
	cpu.soundTimer = byte(cpu.v[x])
	cpu.pc += 2
}

// loadIX
func (cpu *CPU) loadIX() {

}

// addIX
func (cpu *CPU) addIX(x uint16) {
	cpu.i += cpu.v[x]
	cpu.pc += 2
}

// loadF
func (cpu *CPU) loadF(x uint16) {
	cpu.i = cpu.v[x]
	cpu.pc += 2
}

// loadBX
func (cpu *CPU) loadBX(x uint16) {
	cpu.ram.RAM[cpu.i] = cpu.v[x] / 100          // get hundreds digit
	cpu.ram.RAM[cpu.i+1] = (cpu.v[x] % 100) / 10 // get tens digit
	cpu.ram.RAM[cpu.i+2] = cpu.v[x] % 10         // get ones digit
	cpu.pc += 2
}

// loadRegIX
func (cpu *CPU) loadRegIX() {
	for _, v := range cpu.v {
		cpu.ram.RAM[cpu.i] = v
	}

	cpu.pc += 2
}

// loadRegx
func (cpu *CPU) loadRegX() {
	for i := 0; i < len(cpu.v); i++ {
		cpu.v[cpu.i] = cpu.ram.RAM[i]
	}

	cpu.pc += 2
}

// loadVK
func (cpu *CPU) loadVK() {

}
