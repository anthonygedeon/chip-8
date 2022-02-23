use std::fs::{self, File};
use std::io;

/// The amount of memory that the CHIP-8 can hold
pub const MAX_THRESHOLD: usize = 4096;

#[derive(Debug)]
pub struct Memory {
    pub ram: [u16; MAX_THRESHOLD],
}

impl Memory {

    /// Load the `binary` into memory  
    pub fn load_binary(&mut self, binary: &str) -> io::Result<()> {
        let bytes = self.read_binary(binary)?;
        let offset = 512;
        for (i, opcode) in bytes.iter().enumerate() {
            self.ram[offset + i] = *opcode as u16;
        }

        Ok(())
    }
    
    /// Read the entire binary of a file
    fn read_binary(&self, binary: &str) -> io::Result<Vec<u8>> {
        let mut f = fs::read(binary)?;
        Ok(f) 
    }
    
}

