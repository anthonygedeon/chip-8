use std::fs::{self, File};
use std::io;

/// The amount of memory that the CHIP-8 can hold
const MAX_RAM: usize = 4096;

/// CHIP-8 programs start at location 0x200
const END_OF_RESERVED: usize = 0x200;

#[derive(Debug)]
pub struct Memory {
    /// CHIP-8 internal memory capped at 4096 bytes.
    pub ram: [u16; MAX_RAM],
}

impl Memory {

    pub fn new() -> Self {
        Self { ram: [0; MAX_RAM] }     
    }

    /// Load the `rom` into memory.
    pub fn load_rom(&mut self, rom: &str) -> io::Result<()> {
        let bytes = self.read_rom(rom)?;
        for (i, opcode) in bytes.iter().enumerate() {
            self.ram[END_OF_RESERVED + i] = *opcode as u16;
        }

        Ok(())
    }
    
    /// Read the entire rom. 
    fn read_rom(&self, rom: &str) -> io::Result<Vec<u8>> {
        let mut f = fs::read(rom)?;
        Ok(f) 
    }
    
}

#[cfg(test)]
mod test {
   use super::*;
    

}
