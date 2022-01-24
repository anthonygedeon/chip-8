use std::fs;

const MAX_THRESHOLD: usize = 4096;

#[derive(Debug)]
pub struct MemoryMap {
    ram: [u16; MAX_THRESHOLD] 
}

impl MemoryMap {

    fn load_rom(&mut self) {
        // hardcoded path for temporary testing
        match fs::read("roms/IBMLOGO") {
            Ok(bytes) => {
                let start = 512;
                let bytes: Vec<u16> = bytes
                    .chunks_exact(2)
                    .into_iter()
                    .map(|op| u16::from_ne_bytes([op[0], op[1]]))
                    .collect();

                let bytes = bytes.as_slice();
                for (i, opcode) in bytes.iter().enumerate() {
                    self.ram[i+start] = *opcode;
                }

            }

            Err(_) => {
                println!("fail")
            }
        };
    }

}

pub struct Cpu {
    v: u8,  
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
            v: 0,
            i: 0, 
            mem: MemoryMap{ram: [0; MAX_THRESHOLD]}, 
        }
   }

   fn fetch(&mut self) {
        self.mem.load_rom();
        let opcode = self.mem.ram[self.pc as usize];
        let pc = (opcode &0xFFFF) >> 8;
        println!("{:x?}", pc);
        
        match opcode {

        }
   } 

}

fn main() {
    let mut cpu = Cpu::new();
    cpu.fetch();
}
