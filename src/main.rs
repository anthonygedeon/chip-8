extern crate sdl2;

use sdl2::event::Event;
use sdl2::keyboard::Keycode;
use sdl2::pixels;
use sdl2::pixels::{Color, PixelFormatEnum};

use std::fs;
use std::time::Duration;

const MAX_THRESHOLD: usize = 4096;

#[derive(Copy, Clone, Default, Debug)]
pub struct Pixel {
    x: i32,
    y: i32,
    on: i8,
}

#[derive(Debug)]
pub struct MemoryMap {
    ram: [u16; MAX_THRESHOLD],
}

pub struct Display {
    grid: [[Pixel; 64]; 32],
}

pub struct Cpu {
    v: [u8; 16],
    i: u16,

    pc: u16,
    sp: u8,

    delay_timer: u8,
    sound_timer: u8,

    stack: [u16; 15],

    mem: MemoryMap,

    gfx: Display,
}

impl Display {
    pub fn create_grid(&mut self) {
        for (i, row) in self.grid.iter_mut().enumerate() {
            for (j, col) in row.iter_mut().enumerate() {
                col.x = j as i32;
                col.y = i as i32;
                col.on = 0;
            }
        }
    }
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

impl Cpu {
    fn new() -> Self {
        Self {
            stack: [0; 15],
            sound_timer: 0,
            delay_timer: 0,
            pc: 0x200,
            sp: 0,
            v: [0; 16],
            i: 0,
            mem: MemoryMap {
                ram: [0; MAX_THRESHOLD],
            },
            gfx: Display {
                grid: [[Pixel {
                    x: 0i32,
                    y: 0i32,
                    on: 0i8,
                }; 64]; 32],
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
        // println!("{:?}", self.pc);
        let opcode = (self.mem.ram[self.pc as usize] << 8) | (self.mem.ram[(self.pc + 1) as usize]);

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
                let x = (opcode & 0x0F00) >> 8;
                let nn = opcode & 0x00FF;
                self.v[x as usize] = nn as u8;
                println!("LD V{:#x?}, {:#x?}", x, nn);
                self.pc += 2;
            }

            0x7000 => {
                let x = (opcode & 0x0F00) >> 8;
                let nn = opcode & 0x00FF;
                self.v[x as usize] = (x + nn) as u8;
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
                let mut x = self.v[((opcode & 0x0F00) >> 8) as usize];
                let mut y = self.v[((opcode & 0x00FF) >> 4) as usize];
                let nibble = self.mem.ram[self.i as usize];
                println!("Pixels: {:#018b}", nibble);

                /*
                    PSEUDO CODE

                    read n bytes from memory -> start at address stored in I

                    so if reg I is 0x32 we need to put that as the index for the memory array
                    like so -> memory[I] this will get the n byte from memory from whatever address I is

                    what if this was the n byte?

                    memory[I]
                    base 16 -> 0xFD00
                    base 2  -> 1111 1101

                    I think I've been approaching this completely wrong
                    I need to represent the sdl window as bits

                    - The SDL2 window is 640 x 400
                    - Chip 8 is 64x32
                    - A naive solution would be to scale the SDL2 window and center it on the resolution of CHIP-8

                    I need the window to be basically all bits
                    so whatever
                    +-----------------+
                    |00000000000000000|
                    |00000000000000000|
                    |00000000000000000|
                    |00000000000000000|
                    |00000000000000000|
                    |00000000000000000|
                    +-----------------+
                */

                for bytes in nibble..x as u16 {
                    for sprite in bytes..y as u16 {
                        x = x ^ sprite as u8;
                        y = y ^ sprite as u8;

                        self.gfx.grid[x as usize][y as usize].on = 1;
                        if x == 0 {
                            self.v[0xF] = 0;
                        } else if y == 0 {
                            self.v[0xF] = 1;
                        }
                    }
                }

                println!("DRW V{:#x?}, V{:#x?}, {:#x?}", x, y, nibble);
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

fn main() -> Result<(), String> {
    let mut cpu = Cpu::new();

    cpu.gfx.create_grid();

    let sdl_context = sdl2::init()?;
    let video_subsystem = sdl_context.video()?;

    let window = video_subsystem
        .window("Chip 8", 640, 480)
        .position_centered()
        .opengl()
        .build()
        .map_err(|e| e.to_string())?;

    let mut canvas = window.into_canvas().build().map_err(|e| e.to_string())?;
    let texture_creator = canvas.texture_creator();

    let mut texture = texture_creator
        .create_texture_streaming(PixelFormatEnum::RGB24, 256, 256)
        .map_err(|e| e.to_string())?;

    let mut texture2 = texture_creator
        .create_texture_streaming(PixelFormatEnum::RGB24, 256, 256)
        .map_err(|e| e.to_string())?;

    // Create a red-green gradient
    texture.with_lock(None, |buffer: &mut [u8], pitch: usize| {
        for y in 0..256 {
            for x in 0..256 {
                let offset = y * pitch + x * 3;
                buffer[offset] = 255;
                buffer[offset + 1] = 255;
                buffer[offset + 2] = 255;
            }
        }
    })?;

    texture2.with_lock(None, |buffer: &mut [u8], pitch: usize| {
        for y in 0..256 {
            for x in 0..256 {
                let offset = y * pitch + x * 3;
                buffer[offset] = 0;
                buffer[offset + 1] = 0;
                buffer[offset + 2] = 0;
            }
        }
    })?;

    canvas.set_draw_color(Color::RGB(0, 0, 0));
    canvas.set_scale(8.8, 12.0).expect("could not scale window");
    canvas.clear();

    let mut event_pump = sdl_context.event_pump()?;
    let mut rects: Vec<sdl2::rect::Rect> = vec![];
    
    for row in cpu.gfx.grid {
        for col in row {
            let rectangle = sdl2::rect::Rect::new(col.x, col.y, 10, 10);
            if col.on == 1 {
                canvas.copy(&texture, None, Some(rectangle))?;
            } else if col.on == 0 {
                canvas.copy(&texture2, None, Some(rectangle))?;
            }

            rects.push(rectangle);
        }
    }


    'running: loop {
        for event in event_pump.poll_iter() {
            match event {
                Event::Quit { .. }
                | Event::KeyDown {
                    keycode: Some(Keycode::Escape),
                    ..
                } => break 'running,

                _ => {}
            }
        }
        cpu.fetch();
        canvas.present();
        println!("{:?}", rects);

        ::std::thread::sleep(Duration::new(0, 1_000_000_000u32 / 30));
    }

    Ok(())
}
