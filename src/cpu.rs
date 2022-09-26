use crate::display::Display;
use crate::memory::{Memory};

pub struct Instruction {
    x: usize, 

    y: usize, 

    nn: u16, 

    opcode: u16, 

    address: usize, 
}

// CHIP-8 Registers, think of these as variables that can be manipualated by the emulator
#[derive(Default, Debug)]
pub struct Register {
    // the program counter is essentially a pointer that points to the current instruction
    // in the CHIP-8 Memory 
    pc: usize, 

    // the stack pointer points to the word in the stack.
    sp: usize, 

    // special purpose register for delaying the timer
    // when this timer is non-zero i.e { 1.. } then it should be decremented
    dt: u8, 

    // special purpose register for playing sound only when the value is non-zero
    st: u8, 
    
    // stack which is used by the sp register
    stack: [u16; 16], 
    
    // general-purpose registers for v[0] - v[16]
    v: [u8; 16],

    // general-purpose register which is used to store the memory address
    i: usize, 
}

pub struct Cpu {
    pub register: Register, 
    pub memory: Memory, 
    pub display: Display,
}

impl Cpu {

    pub fn new() -> Self {
        let mut cpu = Self {
            register: Register { pc: 0x200, ..Default::default() }, 
            memory: Memory::new(),
            display: Display { grid: [[0; 64]; 32] },
        };
       
        if cpu.memory.load_rom("tests/test_opcode.ch8").is_err() {
           panic!("Could not load the binary to memory.");
        }

        cpu
    }

    fn get_opcode(&self) -> u16 {
        (self.memory.ram[self.register.pc] << 8) | (self.memory.ram[self.register.pc + 1]) 
    }

    fn fetch(&mut self) -> Instruction {
        let opcode = self.get_opcode();
        Instruction {
            opcode, 
            x: ((opcode >> 8) & 0x000F) as usize, 
            y: ((opcode >> 4) & 0x000F) as usize, 
            nn: opcode & 0x00FF, 
            address: (opcode & 0x0FFF) as usize,
        } 
    }

    fn decode(&mut self, instr: Instruction) {
        match instr.opcode & 0xF000 {
            0x0000 => {
                match instr.opcode & 0x00FF {
                    0xE0 => {
                        self.display.clear();
                        println!("CLS");
                        self.register.pc += 2;
                    }

                    0xEE => {
                        self.register.pc = self.register.stack[0xF] as usize;
                        self.register.sp -= 1;
                        println!("RET");
                    }
                    _ => {}
                }
            }

            0x1000 => {
                println!("JP {:#x?}", instr.address);
                self.register.pc = instr.address;
            }

            0x2000 => {
                unimplemented!();
            }

            0x3000 => {
                if self.register.v[instr.x] as u8 == instr.nn as u8 {
                    self.register.pc += 4;
                }
            }

            0x4000 => {
                if self.register.v[instr.x] != self.register.v[instr.y] {
                    self.register.pc += 2;
                }
            }

            0x5000 => {
                unimplemented!();
            }

            0x6000 => {
                self.register.v[instr.x as usize] = instr.nn as u8;
                println!("LD V{:#x?}, {:#x?}", instr.x, instr.nn);
                self.register.pc += 2;
            }

            0x7000 => {
                self.register.v[instr.x as usize] += instr.nn as u8;
                println!("ADD Vx, {:#x?}", instr.nn);
                self.register.pc += 2;
            }

            0x8000 => match instr.opcode & 0x000F {
                0x1 => {}

                0x2 => {}

                0x3 => {}

                0x4 => {}

                0x5 => {}

                0x6 => {}

                0x7 => {}

                0xE => {}

                _ => {}
            },

            0x9000 => {}

            0xA000 => {
                self.register.i = instr.address;
                println!("LD I, {:#x?}", instr.address);
                self.register.pc += 2;
            }

            0xB000 => {}

            0xC000 => {}

            0xD000 => {
                let x = self.register.v[instr.x as usize];
                let y = self.register.v[instr.y as usize];
                let n = instr.opcode & 0x000F;
                
                for height in 0..n {
                    let byte = self.memory.ram[self.register.i + height as usize];
                    for width in 0..=7 {
                        let pixel = (((byte<<width) & 0x80) >> 7) as u8;
                        self.display.set_pos(height as u8 + y, width + x, pixel);
                        if self.display.get_pos(y, x) == 1 {
                            self.register.v[0xF] = 0x01;
                        } else {
                            self.register.v[0xF] = 0x00;
                        }
                    }
                }
                self.register.pc += 2;
            }

            0xE000 => match instr.opcode {
                0xE09E => {}

                0xE0A1 => {}

                _ => {}
            },

            0xF000 => match instr.opcode {
                0xF007 => {}

                0xF00A => {}

                0xF015 => {}

                0xF018 => {}

                0xF01E => {}

                0xF029 => {}

                0xF033 => {}

                0xF055 => {}

                0xF065 => {}

                _ => {}
            },

            _ => {}
        }    
    }

    pub fn cycle(&mut self) {
        let instruction = self.fetch();
        self.decode(instruction) 
    }
}
