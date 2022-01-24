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
                let end = start + bytes.len();
                let bytes: Vec<u16> = bytes.into_iter().map(|opcode| opcode as u16).collect();
                let bytes = bytes.as_slice();
                self.ram[start..end].copy_from_slice(&bytes);
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
        let pc = self.mem.ram[((self.pc & 0xFFFF) >> 8) as usize];

        println!("{:?}", pc);
   } 

}

fn main() {
    let mut cpu = Cpu::new();
    cpu.fetch();
    println!("{:?}", cpu.mem.ram);
}
