use crate::display::Display;
use crate::memory::{Memory, MAX_THRESHOLD};

pub struct Cpu {
    v: [u8; 16],
    i: u16,

    pc: u16,
    sp: u8,

    delay_timer: u8,
    sound_timer: u8,

    stack: [u16; 15],

    memory: Memory, 

    pub display: Display,
}

impl Cpu {

    pub fn new() -> Self {
        Self {
            stack: [0; 15],
            sound_timer: 0,
            delay_timer: 0,
            pc: 0x200,
            sp: 0,
            v: [0; 16],
            i: 0,
            memory: Memory {
                ram: [0; MAX_THRESHOLD],
            },
            display: Display { grid: [[0; 64]; 32] },
        }
    }
    

    pub fn fetch(&mut self) {
        self.memory.load_rom();
        /*
           - we need to get the first byte of the opcode since every byte is an instr.
           - in memory we put the pc to point at location 0x200 since that's where programs start

           pc = 0x200;

           ram[pc] => this grabs 0x00E0

           what we want is 0x0 since that is the first byte and a seperate instr.

           ## Bit Manipulation

           0x00E0 >> 8 we need to shift it 8 bits to the right to get rid of the right most byte
        */
        // println!("{:?}", self.pc);
        let opcode = (self.memory.ram[self.pc as usize] << 8) | (self.memory.ram[(self.pc + 1) as usize]);

        match opcode & 0xF000 {
            0x0000 => {
                match opcode & 0x00F0 {
                    0xE0 => {
                        self.display.clear();
                        println!("CLS");
                        self.pc += 2;
                    }

                    //0xEE => {
                    //    println!("RET");
                    //    self.pc += 2;
                    //}
                    //
                    _ => {}
                }
            }

            0x1000 => {
                let nnn = opcode & 0x0FFF;
                println!("JP {:#x?}", nnn);
                self.pc = nnn;
            }

            // 0x2000 => {
            //      let nnn = opcode & 0x0FFF;
            //      self.sp += 1;
            //      self.stack[self.sp as usize] = self.pc;
            //      self.pc = nnn;
            //      println!("CALL {:x}", nnn);
            // }

            // 0x3000 => {
            //     println!("SE Vx, byte");
            //     self.pc += 2;
            // }

            // 0x4000 => {
            //     println!("SNE Vx, byte");
            //     self.pc += 2;
            // }

            // 0x5000 => {
            //     println!("SE Vx, Vy");
            //     self.pc += 2;
            // }
            0x6000 => {
                let x = (opcode >> 8) & 0x000F;
                let nn = opcode & 0x00FF;
                self.v[x as usize] = nn as u8;
                println!("LD V{:#x?}, {:#x?}", x, nn);
                self.pc += 2;
            }

            0x7000 => {
                let x = (opcode >> 8) & 0x000F;
                let nn = opcode & 0x00FF;
                self.v[x as usize] += nn as u8;
                println!("ADD Vx, {:#x?}", nn);
                self.pc += 2;
            }

            0x8000 => match opcode & 0x000F {
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
                let nnn = opcode & 0x0FFF;
                self.i = nnn;
                println!("LD I, {:#x?}", nnn);
                self.pc += 2;
            }

            0xB000 => {}

            0xC000 => {}

            0xD000 => {
                let x = self.v[((opcode >> 8) & 0x000F) as usize];
                let y = self.v[((opcode >> 4) & 0x000F) as usize];
                let n = opcode & 0x000F;

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

            0xE000 => match opcode {
                0xE09E => {}

                0xE0A1 => {}

                _ => {}
            },

            0xF000 => match opcode {
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
}
