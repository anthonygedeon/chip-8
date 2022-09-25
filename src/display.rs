pub struct Display {
    pub grid: [[u8; 64]; 32],
}

impl Display {
    pub fn clear(&mut self) {
        self.grid = [[0; 64]; 32];
    }
    
    pub fn get_pos(&self, y: u8, x: u8) -> u8 {
        self.grid[y as usize][x as usize] 
    }

    pub fn set_pos(&mut self, y: u8, x: u8, bit: u8) {
        self.grid[y as usize][x as usize] ^= bit
    }
}

