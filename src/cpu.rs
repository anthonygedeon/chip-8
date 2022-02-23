use crate::display::Display;
use crate::memory::{Memory, MAX_THRESHOLD};

pub struct Instruction {
    opcode: u16, 
    x: u16, 
    y: u16, 
    address: usize, 
    nn: u16, 
}

pub struct Cpu {
    v: [u8; 16],
    i: usize,

    pc: usize,
    sp: u8,

    delay_timer: u8,
    sound_timer: u8,

    stack: [u16; 15],

    pub memory: Memory, 

    pub display: Display,
}

impl Cpu {

    pub fn new() -> Self {
        let mut cpu = Self {
            stack: [0; 15],
            sound_timer: 0,
            delay_timer: 0,
            pc: 0x200,
            sp: 0,
            v: [0; 16],
            i: 0,
            memory: Memory { ram: [0; MAX_THRESHOLD], },
            display: Display { grid: [[0; 64]; 32] },
        };
       
        if cpu.memory.load_binary("roms/IBMLOGO").is_err() {
           panic!("Could not load the binary to memory.");
        }

        cpu
    }

    fn get_opcode(&self) -> u16 {
        (self.memory.ram[self.pc] << 8) | (self.memory.ram[(self.pc + 1)])
    }

    fn fetch(&mut self) -> Instruction {

        let opcode = self.get_opcode();
        
        let x = (opcode >> 8) & 0x000F;
        let y = (opcode >> 4) & 0x000F;
        let nn = opcode & 0x00FF;
        let address = (opcode & 0x0FFF) as usize;

        Instruction {
            opcode, 
            x, 
            y, 
            nn, 
            address, 
        } 
    }

    fn decode(&mut self, instr: Instruction) {
        match instr.opcode & 0xF000 {
            0x0000 => {
                match instr.opcode & 0x00F0 {
                    0xE0 => {
                        self.display.clear();
                        println!("CLS");
                        self.pc += 2;
                    }

                    0xEE => {
                        unimplemented!();
                    }
                    _ => {}
                }
            }

            0x1000 => {
                println!("JP {:#x?}", instr.address);
                self.pc = instr.address;
            }

            0x2000 => {
                unimplemented!();
            }

            0x3000 => {
                unimplemented!();
            }

            0x4000 => {
                unimplemented!();
            }

            0x5000 => {
                unimplemented!();
            }

            0x6000 => {
                self.v[instr.x as usize] = instr.nn as u8;
                println!("LD V{:#x?}, {:#x?}", instr.x, instr.nn);
                self.pc += 2;
            }

            0x7000 => {
                self.v[instr.x as usize] += instr.nn as u8;
                println!("ADD Vx, {:#x?}", instr.nn);
                self.pc += 2;
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
                self.i = instr.address;
                println!("LD I, {:#x?}", instr.address);
                self.pc += 2;
            }

            0xB000 => {}

            0xC000 => {}

            0xD000 => {
                let x = self.v[((instr.opcode >> 8) & 0x000F) as usize];
                let y = self.v[((instr.opcode >> 4) & 0x000F) as usize];
                let n = instr.opcode & 0x000F;

                for height in 0..n {
                    for width in 0..8 {
                        let byte = self.memory.ram[self.i as usize];
                        
                        let x_pos = x;
                        let y_pos = y;

                        self.display.set_pos(height as u8 + y_pos, width as u8 + x_pos, 1);
                        if self.display.get_pos(y, x) == 0 {
                            self.v[0xF] = 1;
                        } else {
                            self.v[0xF] = 0;
                        }


                    }
                }

                self.pc += 2;
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
