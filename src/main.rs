use std::fs;

const MAX_THRESHOLD: usize = 4096;

#[derive(Debug)]
pub struct MemoryMap {
    ram: [u16; MAX_THRESHOLD],
}

impl MemoryMap {
    fn load_rom(&mut self) {
        // hardcoded path for temporary testing
        match fs::read("roms/IBMLOGO") {
            Ok(bytes) => {
                let start = 512;

                let bytes = bytes.as_slice();
                for (i, opcode) in bytes.iter().enumerate() {
                    self.ram[i + start] = *opcode as u16;
                }
            }
            Err(_) => {
                println!("fail")
            }
        };
    }
}

pub struct Cpu {
    v: [u8; 15],
    i: u8,

    pc: u16,
    sp: u8,

    delay_timer: u8,
    sound_timer: u8,

    stack: [u16; 15],

    mem: MemoryMap,
}

impl Cpu {
    fn new() -> Self {
        Self {
            stack: [0; 15],
            sound_timer: 0,
            delay_timer: 0,
            pc: 0x200,
            sp: 0,
            v: [0; 15],
            i: 0,
            mem: MemoryMap {
                ram: [0; MAX_THRESHOLD],
            },
        }
    }

    fn fetch(&mut self) {
        self.mem.load_rom();

        /*
           - we need to get the first byte of the opcode since every byte is an instr.
           - in memory we put the pc to point at location 0x200 since that's where programs start

           pc = 0x200;

           ram[pc] => this grabs 0x00E0

           what we want is 0x0 since that is the first byte and a seperate instr.

           ## Bit Manipulation

           0x00E0 >> 8 we need to shift it 8 bits to the right to get rid of the right most byte
        */
        loop {
            // println!("{:?}", self.pc);
            let opcode =
                (self.mem.ram[self.pc as usize] << 8) | (self.mem.ram[(self.pc + 1) as usize]);

            // println!("opcode: {:x?} pc: {:?}", opcode&0xF000,  self.pc);
            //println!("{:x?}", opcode);
            match opcode & 0xF000 {
                0x0000 => {
                    match opcode & 0x00F0 {
                        0xE0 => {
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
                    println!("JP {}", nnn);
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
                    let x = (opcode & 0x0F00) >> 8;
                    let nn = opcode & 0x00FF;
                    self.v[x as usize] = nn as u8;
                    println!("LD V{}, {}", x, nn);
                    self.pc += 2;
                }

                0x7000 => {
                    let x = (opcode & 0x0F00) >> 8;
                    let nn = opcode & 0x00FF;
                    self.v[x as usize] = (x + nn) as u8;
                    println!("ADD Vx, byte");
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
                    self.i = nnn as u8;
                    println!("LD I, {}", nnn);
                    self.pc += 2;
                }

                0xB000 => {}

                0xC000 => {}

                0xD000 => {
                    println!("DRW Vx, Vy, nibble");
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
}

fn main() {
    let mut cpu = Cpu::new();
    cpu.fetch();
}
