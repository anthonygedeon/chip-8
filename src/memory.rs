use std::fs;

/// The amount of memory that the CHIP-8 can hold
pub const MAX_THRESHOLD: usize = 4096;

#[derive(Debug)]
pub struct Memory {
    pub ram: [u16; MAX_THRESHOLD],
}

impl Memory {
    pub fn load_rom(&mut self) {
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
