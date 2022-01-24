package main

import "fmt"

func Disassemble(c *CPU) {
	c.fetchOpcode()
	fmt.Printf("%-4X %#4X\t", c.pc, c.opcode)
	switch c.opcode & 0xF000 {
	case 0x0000:
		switch c.opcode & 0x00F0 {
		case 0x00E0:
			fmt.Println("CLS")
		case 0x00EE:
			fmt.Println("RET")
		default:
			nnn := c.opcode & 0x0FFF
			fmt.Printf("SYS %X\n", nnn)
		}
	case 0x1000:
		c.pc = c.opcode & 0x0FFF
		fmt.Printf("JP %X\n", c.pc)
	case 0x2000:
		nnn := c.opcode & 0x0FFF
		fmt.Printf("CALL %X\n", nnn)
	case 0x3000:
		nn := byte(c.opcode & 0x00FF)
		x := (c.opcode >> 8) & 0x000F
		fmt.Printf("SE V%X, %X\n", x, nn)
	case 0x4000:
		x := (c.opcode >> 8) & 0x000F
		nn := byte(c.opcode & 0x00FF)
		fmt.Printf("SNE V%X, %X\n", x, nn)
	case 0x5000:
		x := (c.opcode >> 8) & 0x000F
		y := (c.opcode >> 4) & 0x000F
		fmt.Printf("SE V%X, V%X\n", x, y)
	case 0x6000:
		v := c.opcode >> 8 & 0x000F
		nn := byte(c.opcode & 0x00FF)
		fmt.Printf("LD V%X, %X\n", v, nn)
	case 0x7000:
		nn := byte(c.opcode & 0x00FF)
		x := (c.opcode >> 8) & 0x000F
		fmt.Printf("ADD V%X, %X\n", x, nn)
	case 0x8000:
		switch c.opcode & 0x000F {
		case 0x0:
			x := (c.opcode >> 8) & 0x000F
			y := (c.opcode >> 4) & 0x000F
			fmt.Printf("LD V%X, V%X\n", x, y)
		case 0x1:
			x := (c.opcode >> 8) & 0x000F
			y := (c.opcode >> 4) & 0x000F
			fmt.Printf("OR V%X, V%X\n", x, y)
		case 0x2:
			x := (c.opcode >> 8) & 0x000F
			y := (c.opcode >> 4) & 0x000F
			fmt.Printf("AND V%X, V%X\n", x, y)
		case 0x3:
			x := (c.opcode >> 8) & 0x000F
			y := (c.opcode >> 4) & 0x000F
			fmt.Printf("XOR V%X, V%X\n", x, y)
		case 0x4:
			x := (c.opcode >> 8) & 0x000F
			y := (c.opcode >> 4) & 0x000F
			fmt.Printf("ADD V%X, V%X\n", x, y)
		case 0x5:
			x := (c.opcode >> 8) & 0x000F
			y := (c.opcode >> 4) & 0x000F
			fmt.Printf("SUB V%X, V%X\n", x, y)
		case 0x6:
			x := (c.opcode >> 8) & 0x000F
			fmt.Printf("SHR V%X\n", x)
		case 0x7:
			x := (c.opcode >> 8) & 0x000F
			y := (c.opcode >> 4) & 0x000F
			fmt.Printf("SUBN V%X, V%X\n", x, y)
		case 0xE:
			x := (c.opcode >> 8) & 0x000F
			fmt.Printf("SHL V%X\n", x)
		default:
			fmt.Printf("UNKN %X\n", c.opcode)
		}
	case 0x9000:
		x := (c.opcode >> 8) & 0x000F
		y := (c.opcode >> 4) & 0x000F
		fmt.Printf("SNE V%X, V%X\n", x, y)
	case 0xA000:
		nnn := c.opcode & 0x0FFF
		fmt.Printf("LD I, %X\n", nnn)
	case 0xB000:
		nnn := c.opcode & 0x0FFF
		fmt.Printf("JP V0, %X\n", nnn)
	case 0xC000:
		nn := byte(c.opcode & 0x00FF)
		x := (c.opcode >> 8) & 0x000F
		fmt.Printf("RND V%X, %X\n", x, nn)
	case 0xD000:
		n := c.opcode & 0x000F
		x := (c.opcode >> 8) & 0x000F
		y := (c.opcode >> 4) & 0x000F
		fmt.Printf("DRW V%X, V%X, %X\n", x, y, n)
	case 0xE000:
		switch c.opcode & 0x00FF {
		case 0x9E:
			x := (c.opcode >> 8) & 0x000F
			fmt.Printf("SKP V%X\n", x)
		case 0xA1:
			x := (c.opcode >> 8) & 0x000F
			fmt.Printf("SKNP V%X\n", x)
		default:
			fmt.Printf("UNKN %X\n", c.opcode)
		}
	case 0xF000:
		switch c.opcode & 0x00FF {
		case 0x07:
			x := (c.opcode >> 8) & 0x000F
			fmt.Printf("LD V%X, DT\n", x)
		case 0x0A:
			x := (c.opcode >> 8) & 0x000F
			fmt.Printf("LD V%X, K\n", x)
		case 0x15:
			x := (c.opcode >> 8) & 0x000F
			fmt.Printf("LD DT, V%X\n", x)
		case 0x18:
			x := (c.opcode >> 8) & 0x000F
			fmt.Printf("LD ST, V%X\n", x)
		case 0x1E:
			x := (c.opcode >> 8) & 0x000F
			fmt.Printf("ADD I, V%X\n", x)
		case 0x29:
			x := (c.opcode >> 8) & 0x000F
			fmt.Printf("LD F, V%X\n", x)
		case 0x33:
			x := (c.opcode >> 8) & 0x000F
			fmt.Printf("LD B, V%X\n", x)
		case 0x55:
			x := (c.opcode >> 8) & 0x000F
			fmt.Printf("LD [I], V%X\n", x)
		case 0x65:
			x := (c.opcode >> 8) & 0x000F
			fmt.Printf("LD V%X, [I]\n", x)
		default:
			fmt.Printf("UNKN %X\n", c.opcode)
		}
	default:
		fmt.Printf("UNKN %X\n", c.opcode)
	}
	c.pc += 2
}
