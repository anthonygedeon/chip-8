extern crate sdl2;

use sdl2::event::Event;
use sdl2::keyboard::Keycode;
use sdl2::pixels::{Color, PixelFormatEnum};

use std::time::Duration;

mod memory;
mod display;
mod cpu;

fn main() -> Result<(), String> {
    let mut cpu = cpu::Cpu::new();

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

    // create white texture
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
    
    // create black texture
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
        cpu.cycle();
        canvas.present();
        
        for (mut i, row) in cpu.display.grid.iter().enumerate() {
            for (mut j, _col) in row.iter().enumerate() {
                let rectangle = sdl2::rect::Rect::new((j) as i32, (i) as i32, 10, 10);

                if cpu.display.grid[i][j] == 1 {
                    canvas.copy(&texture, None, Some(rectangle))?;
                } else if cpu.display.grid[i][j] == 0 {
                    canvas.copy(&texture2, None, Some(rectangle))?;
                }

                j + 10;
                i + 10;
            }
        }

        ::std::thread::sleep(Duration::new(0, 1_000_000_000u32 / 30));
    }

    Ok(())
}
