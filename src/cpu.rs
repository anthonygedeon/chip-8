use std::time::{SystemTime, UNIX_EPOCH};

use crate::keyboard::Keyboard;
use crate::display::Display;
use crate::memory::{Memory, FONT_SET};

pub struct Instruction {
    x: usize, 

    y: usize, 

    nn: u8, 

    opcode: u16, 

    nnn: usize, 
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
    pub keyboard: Keyboard
}

fn rand() -> u8 {
    let sys_time = SystemTime::now();
    let time = sys_time.duration_since(UNIX_EPOCH).unwrap().as_secs();
    let value = ((1103515245 * time + 12345) % 2_u64.pow(31)) as u8;
    value
}

impl Cpu {

    pub fn new() -> Self {
        let mut cpu = Self {
            register: Register { pc: 0x200, ..Default::default() }, 
            memory: Memory::new(),
            display: Display { grid: [[0; 64]; 32] },
            keyboard: Keyboard { key: 0 }, 
        };
       
        if cpu.memory.load_rom("res/chip8-test-suite.ch8").is_err() {
           panic!("Could not load the binary to memory.");
        }

        if cpu.memory.load_font(FONT_SET).is_err() {
           panic!("Could not load the font to memory.");
        }
        
        // 1 - IBM LOGO
        // 2 - Corax89's opcode test
        // 3 - Flags test
        // 4 - Quirks test
        // 5 - Keypad test
        
        // DEBUG PURPOSES
        cpu.memory.ram[0x1FF] = 4;

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
            nn: (opcode & 0x00FF) as u8, 
            nnn: (opcode & 0x0FFF) as usize,
        } 
    }

    fn decode(&mut self, instr: Instruction) {
        println!("PC: {:X}\nSP: {}\nI: {}\nSTK: {:?}\nV: {:?}", self.register.pc, self.register.sp, self.register.i, self.register.stack, self.register.v);
        println!("Keypress: {}", self.keyboard.key);

        if self.register.dt > 0 {
            self.register.dt -= 1;
        }

        if self.register.st > 0 {
            self.register.st -= 1;
            // TODO: Make beep noise
        }

        match instr.opcode & 0xF000 {
            0x0000 => {
                match instr.opcode & 0x00FF {
                    0xE0 => {
                        println!("CLS");
                        self.display.clear();
                        self.register.pc += 2;
                    }

                    0xEE => {
                        println!("RET");
                        self.register.sp -= 1;
                        self.register.pc = self.register.stack[self.register.sp] as usize;
                        self.register.pc += 2;
                    }
                    _ => unreachable!()
                }
            }

            0x1000 => {
                println!("0x{:X} JP nnn={:#x?}", instr.opcode, instr.nnn);
                self.register.pc = instr.nnn;
            }

            0x2000 => {
                println!("0x{:X?} CALL {}", instr.opcode, instr.nnn);
                self.register.stack[self.register.sp] = self.register.pc as u16;
                self.register.sp += 1;
                self.register.pc = instr.nnn;
            }
    
            // SKIP
            0x3000 => {
                println!("0x{:X?} SE V[{}], {}", instr.opcode, self.register.v[instr.x], instr.nn);
                if self.register.v[instr.x] == instr.nn {
                    self.register.pc += 4;
                } else {
                    self.register.pc += 2;
                }
            }

            0x4000 => {
                println!("0x{:X?} SNE V[{}], {}", instr.opcode, self.register.v[instr.x], instr.nn);
                if self.register.v[instr.x] != instr.nn {
                    self.register.pc += 4;
                } else {
                    self.register.pc += 2;
                }
            }

            0x5000 => {
                println!("0x{:X?} SE {}, {}", instr.opcode, self.register.v[instr.x], self.register.v[instr.y]);
                 if self.register.v[instr.x] == self.register.v[instr.y] {
                    self.register.pc += 4;
                } else {
                    self.register.pc += 2;
                }
            }    

            0x9000 => {
                println!("0x{:X?} SNE {}, {}", instr.opcode, self.register.v[instr.x], self.register.v[instr.y]);
                 if self.register.v[instr.x] != self.register.v[instr.y] {
                    self.register.pc += 4;
                } else {
                    self.register.pc += 2;
                }
            }
            
            0x6000 => {
                println!("0x{:X} LD V[{:#x?}], {:#x?}", instr.opcode, instr.x, instr.nn);
                self.register.v[instr.x] = instr.nn;
                self.register.pc += 2;
            }

            0x7000 => {
                println!("0x{:X?} ADD Vx, {:#x?}", instr.opcode, instr.nn);
                self.register.v[instr.x] = self.register.v[instr.x].wrapping_add(instr.nn as u8);
                self.register.pc += 2;
            }

            0x8000 => match instr.opcode & 0x000F {
                0x0 => {
                    println!("0x{:X?} LD V[{}], V[{}]", instr.opcode, self.register.v[instr.x], self.register.v[instr.y]);
                    self.register.v[instr.x] = self.register.v[instr.y];
                    self.register.pc += 2;
                }

                0x1 => {
                    println!("0x{:X?}, OR V[{}], V[{}]", instr.opcode, self.register.v[instr.x], self.register.v[instr.y]);
                    self.register.v[instr.x] |= self.register.v[instr.y];
                    self.register.pc += 2;
                }

                0x2 => {
                    println!("0x{:X?} AND V[{}], V[{}]", instr.opcode, self.register.v[instr.x], self.register.v[instr.y]);
                    self.register.v[instr.x] &= self.register.v[instr.y];
                    self.register.pc += 2;
                }

                0x3 => {
                    println!("0x{:X?} XOR V[{}], V[{}]", instr.opcode, self.register.v[instr.x], self.register.v[instr.y]);
                    self.register.v[instr.x] ^= self.register.v[instr.y];
                    self.register.pc += 2;
                }
                
                0x4 => {
                    println!("0x{:X?} ADD V[{}], V[{}]", instr.opcode, self.register.v[instr.x], self.register.v[instr.y]);
                    if self.register.v[instr.x].checked_add(self.register.v[instr.y]) == None {
                        self.register.v[instr.x] = self.register.v[instr.x].wrapping_add(self.register.v[instr.y]);
                        self.register.v[0xF] = 1;
                    } else {
                        self.register.v[instr.x] = self.register.v[instr.x].wrapping_add(self.register.v[instr.y]);
                        self.register.v[0xF] = 0;
                    }
                    self.register.pc += 2;
                }

                0x5 => {
                    println!("0x{:X?} SUB V[{}], V[{}]", instr.opcode, self.register.v[instr.x], self.register.v[instr.y]);

                    self.register.v[instr.x] = self.register.v[instr.x].wrapping_sub(self.register.v[instr.y]);

                    if self.register.v[instr.x] > self.register.v[instr.y] {
                        self.register.v[0xF] = 1;
                    } else if self.register.v[instr.x].checked_sub(self.register.v[instr.y]) == None {
                        self.register.v[0xF] = 0;
                    }


                    self.register.pc += 2;
                }

                0x6 => {
                    println!("0x{:X?} SHR V[{}] {{, V[{}]}}", instr.opcode, self.register.v[instr.x], self.register.v[instr.y]);
                    let tmp = self.register.v[instr.x] & 0x01;
                    self.register.v[instr.x] >>= 1;
                    self.register.v[0xF] = tmp;
                    self.register.pc += 2;
                }

                0x7 => {
                    println!("0x{:X?} SUBN V[{}], V[{}]", instr.opcode, self.register.v[instr.x], self.register.v[instr.y]);

                    self.register.v[instr.x] = self.register.v[instr.y].wrapping_sub(self.register.v[instr.x]);

                    if self.register.v[instr.x] < self.register.v[instr.y] {
                        self.register.v[0xF] = 1;
                    } else if self.register.v[instr.y].checked_sub(self.register.v[instr.x]) == None {
                        self.register.v[0xF] = 0;
                    }

                    self.register.pc += 2;
                }

                0xE => {
                    println!("0x{:X?} SHL V[{}], {{, V[{}]}}", instr.opcode, self.register.v[instr.x], self.register.v[instr.y]);
                    let bit = (self.register.v[instr.x] & 0x80) >> 7;
                    self.register.v[instr.x] <<= 1;
                    self.register.v[0xF] = bit;
                    self.register.pc += 2;
                }

                _ => unreachable!()
            },


            0xA000 => {
                println!("0x{:X?} LD I, {:#x?}", instr.opcode, instr.nnn);
                self.register.i = instr.nnn;
                self.register.pc += 2;
            }

            0xB000 => {
                println!("0x{:X?} JP {}, {}", instr.opcode, self.register.v[0], instr.nnn);
                self.register.pc = instr.nnn as usize + self.register.v[0] as usize;
            }

            0xC000 => {
                self.register.v[instr.x] = rand() & instr.nn;
                self.register.pc += 2;
            }

            0xD000 => {
                let x = self.register.v[instr.x as usize];
                let y = self.register.v[instr.y as usize];
                let n = instr.opcode & 0x000F;
                println!("0x{:X?} DRW V[{}], V[{}], {}", instr.opcode, x, y, n);
                
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

            0xE000 => match instr.opcode & 0x00FF {
                0x9E => {
                    println!("SKP Vx={}", self.register.v[instr.x]);
                    if self.keyboard.key == self.register.v[instr.x] {
                        self.register.pc += 4;
                    } else {
                        self.register.pc += 2;
                    }
               }

                0xA1 => {
                    println!("SKNP Vx={}", self.register.v[instr.x]);
                    if self.keyboard.key != self.register.v[instr.x] {
                        self.register.pc += 4;
                    } else {
                        self.register.pc += 2;
                    }
                }

                _ => unreachable!()
            },

            0xF000 => match instr.opcode & 0x00FF {
                0x07 => {
                    println!("0x{:X?} LD V[{}], DT={}", instr.opcode, self.register.v[instr.x], self.register.dt);
                    self.register.v[instr.x] = self.register.dt;
                    self.register.pc += 2;
                }

                0x0A => {
                        println!("LD Vx={}, K={}", self.register.v[instr.x], self.keyboard.key);
                        if self.keyboard.is_pressed() {
                            self.register.v[instr.x] = self.keyboard.key;
                            self.register.pc += 2;
                        }                 
                }

                0x15 => {
                    println!("0x{:X?} LD DT={}, V[{}]", instr.opcode, self.register.dt, self.register.v[instr.x]);
                    self.register.dt = self.register.v[instr.x];
                    self.register.pc += 2;
                }

                0x18 => {
                    println!("0x{:X?} LD ST={}, V[{}]", instr.opcode, self.register.st, self.register.v[instr.x]);
                    self.register.st = self.register.v[instr.x];
                    self.register.pc += 2;
                }

                0x1E => {
                    println!("0x{:X?} ADD I={}, V[{}]", instr.opcode, self.register.i, self.register.v[instr.x]);
                    self.register.i += self.register.v[instr.x] as usize;
                    self.register.pc += 2;
                }

                0x29 => {
                    println!("0x{:X?} LD F={}, V[{}]", instr.opcode, self.memory.ram[0x50], self.register.v[instr.x]);
                    self.register.i = self.memory.ram[0x50] as usize; 
                    self.register.pc += 2;
                }

                0x33 => {
                    let v_byte = self.register.v[instr.x];
                    println!("0x{:X?} LD B={}, V[{}]", instr.opcode, self.memory.ram[self.register.i], v_byte);
                    self.memory.ram[self.register.i]   = (v_byte / 100 % 10) as u16;
                    self.memory.ram[self.register.i+1] = (v_byte / 10 % 10) as u16;
                    self.memory.ram[self.register.i+2] = (v_byte / 1 % 10) as u16;
                    self.register.pc += 2;
                }

                0x55 => {
                    println!("0x{:X?} LD I=[{}], V[{}]", instr.opcode, self.register.i, self.register.v[instr.x]);
                    for i in 0..=instr.x {
                        self.memory.ram[self.register.i+i] = self.register.v[i] as u16;
                    }
                    self.register.pc += 2;
                }

                0x65 => {
                    println!("0x{:X?} LD V[{}], I=[{}]", instr.opcode, self.register.v[instr.x], self.register.i);
                    for i in 0..=instr.x {
                        self.register.v[i] = self.memory.ram[self.register.i+i] as u8;
                    }
                    self.register.pc += 2;
                }

                _ => unreachable!()
            },

            _ => unreachable!()
        }    
    }

    pub fn cycle(&mut self) {
        let instruction = self.fetch();
        self.decode(instruction) 
    }
}
