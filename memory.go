package main

import (
	"io"
	"log"
	"os"
)

// FontSet sprites representing the hexadecimal digits 0 through F.
// These sprites are 5 bytes long, or 8x5 pixels.
// The data should be stored in the interpreter area of Chip-8 memory (0x000 to 0x1FF).
// Below is a listing of each character's bytes, in binary and hexadecimal
var FontSet = [80]byte{
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

// start of programs
const offset = 512

// A Memory contains the ram of the program
type Memory struct {
	// 4KB of ram
	RAM [4096]byte
}

// Read
func (mem *Memory) Read(file io.Reader) ([]byte, error) {
	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

// LoadProgram
func (mem *Memory) LoadProgram(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("failed to load program into memory: %q", err)
	}

	defer file.Close()

	fInfo, err := file.Stat()
	if int(fInfo.Size()) >= len(mem.RAM) {
		log.Fatalf("file is too large: %q", err)
	}

	bytes, err := mem.Read(file)
	if err != nil {
		log.Fatalf("failed to read bytes from file: %q", err)
	}

	for i, data := range bytes {
		mem.RAM[i+offset] = byte(data)
	}
}
