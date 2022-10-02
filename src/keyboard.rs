pub struct Keyboard {
    pub key: u8,
}

impl Keyboard {
    pub fn set_keypress(&mut self, keycode: u8) {
        self.key = keycode;
    }

    pub  fn is_pressed(&self) -> bool {
        self.key != 0 
    }
}
