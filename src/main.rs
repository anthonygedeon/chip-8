extern crate sdl2;

use sdl2::event::Event;
use sdl2::keyboard::Keycode;
use sdl2::pixels::{Color, PixelFormatEnum};

use std::thread;
use std::time::Duration;

mod keyboard;
mod cpu;
mod display;
mod memory;

const WINDOW_WIDTH: u32 = 800;
const WINDOW_HEIGHT: u32 = 400;
const WINDOW_TITLE: &str = "CHIP-8";

fn main() -> Result<(), String> {
    let mut cpu = cpu::Cpu::new();

    let sdl_context = sdl2::init()?;
    let video_subsystem = sdl_context.video()?;

    let window = video_subsystem
        .window(WINDOW_TITLE, WINDOW_WIDTH, WINDOW_HEIGHT)
        .position_centered()
        .opengl()
        .build()
        .map_err(|e| e.to_string())?;

    let mut canvas = window.into_canvas().build().map_err(|e| e.to_string())?;

    canvas
        .set_scale(WINDOW_WIDTH as f32 / 64.0, WINDOW_HEIGHT as f32 / 32.0)
        .expect("could not scale window");

    let mut event_pump = sdl_context.event_pump()?;
    'running: loop {

        canvas.set_draw_color(Color::RGB(0, 0, 0));
        canvas.clear();

        for event in event_pump.poll_iter() {
            match event {
                Event::Quit { .. }
                | Event::KeyDown {
                    keycode: Some(Keycode::Escape),
                    ..
                } => break 'running,
                Event::KeyDown {
                    keycode: Some(Keycode::Num1),
                    ..
                } => cpu.keyboard.set_keypress(0x1),
                Event::KeyDown {
                    keycode: Some(Keycode::Num2),
                    ..
                } => cpu.keyboard.set_keypress(0x2),
                Event::KeyDown {
                    keycode: Some(Keycode::Num3),
                    ..
                } => cpu.keyboard.set_keypress(0x3),
                Event::KeyDown {
                    keycode: Some(Keycode::Q),
                    ..
                } => cpu.keyboard.set_keypress(0x4),
                Event::KeyDown {
                    keycode: Some(Keycode::W),
                    ..
                } => cpu.keyboard.set_keypress(0x5),
                Event::KeyDown {
                    keycode: Some(Keycode::Num6),
                    ..
                } => (),
                Event::KeyDown {
                    keycode: Some(Keycode::Num7),
                    ..
                } => (),
                Event::KeyDown {
                    keycode: Some(Keycode::Num8),
                    ..
                } => (),
                Event::KeyDown {
                    keycode: Some(Keycode::Num9),
                    ..
                } => (),
                Event::KeyDown {
                    keycode: Some(Keycode::A),
                    ..
                } => (),
                Event::KeyDown {
                    keycode: Some(Keycode::B),
                    ..
                } => (),
                Event::KeyDown {
                    keycode: Some(Keycode::C),
                    ..
                } => (),
                Event::KeyDown {
                    keycode: Some(Keycode::D),
                    ..
                } => (),
                Event::KeyDown {
                    keycode: Some(Keycode::E),
                    ..
                } => (),
                Event::KeyDown {
                    keycode: Some(Keycode::F),
                    ..
                } => (),
                Event::KeyUp {
                    keycode: Some(Keycode::Num1),
                   ..  
                } => cpu.keyboard.set_keypress(0), 

                Event::KeyUp {
                    keycode: Some(Keycode::Num2),
                   ..  
                } => cpu.keyboard.set_keypress(0), 

                Event::KeyUp {
                    keycode: Some(Keycode::Num3),
                   ..  
                } => cpu.keyboard.set_keypress(0), 
                Event::KeyUp {
                    keycode: Some(Keycode::Q),
                   ..  
                } => cpu.keyboard.set_keypress(0), 
                Event::KeyUp {
                    keycode: Some(Keycode::W),
                   ..  
                } => cpu.keyboard.set_keypress(0), 
                _ => {}
            }
        }

        cpu.cycle();
        
        for (i, row) in cpu.display.grid.iter().enumerate() {
            for (j, _) in row.iter().enumerate() {
                let rectangle = sdl2::rect::Rect::new(j as i32, i as i32, 10, 10);

                // draw white cell
                if cpu.display.grid[i][j] == 1 {
                    canvas.set_draw_color(Color::RGB(255, 255, 255));
                    canvas.draw_rect(rectangle)?;
               
                // draw black cell
                } else if cpu.display.grid[i][j] == 0 {
                    canvas.set_draw_color(Color::RGB(0, 0, 0));
                    canvas.draw_rect(rectangle)?;
                }
            }
        }

        canvas.present();

        ::std::thread::sleep(Duration::new(0, 1_000_000_000u32 / 250));
    }

    Ok(())
}
