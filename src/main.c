#include <SDL_audio.h>
#include <SDL_events.h>
#include <SDL_render.h>
#include <SDL_surface.h>
#include <SDL_video.h>
#include <errno.h>
#include <SDL.h>
#include <stdint.h>
#include <stdio.h>
#include <stdbool.h>
#include <stdlib.h>

#define WINDOW_WIDTH 640
#define WINDOW_HEIGHT 480

#define FB_ROWS 32
#define FB_COLS 64

#define ARRAY_LENGTH(x) (sizeof(x) / sizeof(x[0]))

#define MAX_STACK_SIZE 100
#define MEMORY_MAX 4096
#define VARIABLE_MAX 16

#define DISPLAY_SIZE 64 * 32

typedef uint16_t addr_t;

typedef enum {
	KEY_UP   = 0,
	KEY_DOWN = 1 
} key_event;

typedef struct {
	uint8_t ram[MEMORY_MAX];
} memory_t;

memory_t memory_new() {
	return (memory_t){ .ram = { 0 } };
}

const uint8_t FONT_SET[80] = {
    0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
    0x20, 0x60, 0x20, 0x20, 0x70, // 1
    0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
    0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
    0x90, 0x90, 0xF0, 0x10, 0x10, // 4
    0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
    0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
    0xF0, 0x10, 0x20, 0x40, 0x40, // 7
    0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
    0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
    0xF0, 0x90, 0xF0, 0x90, 0x90, // A
    0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
    0xF0, 0x80, 0x80, 0x80, 0xF0, // C
    0xE0, 0x90, 0x90, 0x90, 0xE0, // D
    0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
    0xF0, 0x80, 0xF0, 0x80, 0x80, // F
};

void beep(void *userdata, Uint8 *stream, int len) {
	Uint8 wave[] = {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255};
	int wrap = 0;
	for (int i = 0; i < len; i++) {
		stream[i] = wave[wrap];
		wrap = (wrap + 1) % ARRAY_LENGTH(wave);
	}
}

SDL_AudioSpec *audio_spec_new() {
	SDL_AudioSpec *spec = malloc(sizeof(SDL_AudioSpec));
	spec->channels = 1;
	spec->freq = 44100;
	spec->samples = 1024;
	spec->format = AUDIO_U8;
	spec->callback = *beep;
	return spec;
}

void memory_write(memory_t *memory, size_t size, const uint8_t *data) {
	if (memory == NULL) {
		printf("Error: invalid pointer");
	}
	
	int offset = 0x200;
	for (size_t i = 0; i < size;  i++) {
		memory->ram[offset] = data[i];
		offset++;
	}

	int font_offset = 0x050;
	for (size_t i = 0; i < 80;  i++) {
		memory->ram[font_offset] = FONT_SET[i];
		offset++;
	}

}

typedef struct {
	uint8_t fb[FB_ROWS][FB_COLS];
} display_t;

display_t display_new() {
	return (display_t){ .fb = { 0 }};
}

typedef struct {
	uint16_t pc;
	uint8_t  sp;

	uint16_t i;
	
	uint8_t delay_timer;
	uint8_t sound_timer;

	uint8_t v[VARIABLE_MAX];

	uint16_t stack[MAX_STACK_SIZE];
} reg_t;

reg_t register_new() {
	return (reg_t){
		.pc = (uint16_t)0x200,
		.sp = 0,
		.i  = 0,
		.delay_timer = 0,
		.sound_timer = 0,
		.v = { 0 },
		.stack = { 0 },
	};
}

typedef struct {
	uint16_t keypad;
} keyboard_t;

keyboard_t keyboard_new() {
	return (keyboard_t){ .keypad = 0x0000 };
}

bool keyboard_isset(keyboard_t keyboard, uint8_t key) {
	return (keyboard.keypad >> key) & 0x1U;
}

void keyboard_setkey(keyboard_t *keyboard, uint8_t key) {
	keyboard->keypad |= (0x1U << key);	
}

void keyboard_unsetkey(keyboard_t *keyboard, uint8_t key) {
	keyboard->keypad &= ~(0x1U << key);	
}

typedef struct {
	SDL_AudioDeviceID id;
	bool is_playing;
} sound_t;

typedef struct {
	reg_t reg;
	keyboard_t keyboard;
	memory_t memory;
	display_t display;
	sound_t sound;
} cpu_t;



sound_t sound_new() {
	return (sound_t){ .id = 0, .is_playing = true };
}

cpu_t cpu_new() {
	return (cpu_t){ 
		.reg     = register_new(),
		.memory  = memory_new(),
		.display = display_new(),
		.keyboard = keyboard_new(),
		.sound = sound_new()
	};
}

uint16_t cpu_fetch(cpu_t *cpu) {
	uint16_t opcode = (cpu->memory.ram[cpu->reg.pc] << 8) | (cpu->memory.ram[cpu->reg.pc + 1]);
	cpu->reg.pc += 2;
	return opcode;
}

void cpu_decode(cpu_t *cpu, uint16_t opcode) {
	printf("%.4X %.2X %-6.2X", cpu->reg.pc, opcode >> 8, opcode & 0x00FF);
	
	uint16_t addr = opcode & 0x0FFF;
	uint8_t nibble = opcode & 0x000F;
	uint8_t x = (opcode >> 8) & 0x000F;
	uint8_t y = (opcode >> 4) & 0x000F;
	uint8_t byte = opcode & 0x00FF;

	
	switch (opcode & 0xF000) {
		case 0x0000:
			switch(opcode & 0x00FF) {
				case 0xE0:
					printf("CLS\n");
					memset(cpu->display.fb, 0, sizeof(cpu->display.fb[0][0]) * DISPLAY_SIZE);
					break;
				case 0xEE:
					printf("RET\n");
					cpu->reg.pc = cpu->reg.stack[--cpu->reg.sp];
					break;
				default:
					printf("UNKP\n");
					break;
			}
			break;
		case 0x1000:
			printf("JP 0x%.3X\n", addr);
			cpu->reg.pc = addr;
			break;
		case 0x2000:
			printf("CALL 0x%.3X\n", addr);
			cpu->reg.stack[cpu->reg.sp++] = cpu->reg.pc;
			cpu->reg.pc = addr;
			break;
		case 0x3000:
			printf("SE V[0x%X], %X\n", x, byte);
			if (cpu->reg.v[x] == byte) {
				cpu->reg.pc += 2;
			}
			break;
		case 0x4000:
			printf("SNE V[0x%X], %X\n", x, byte);
			if (cpu->reg.v[x] != byte) {
				cpu->reg.pc += 2;
			} 			
			break;
		case 0x5000:
			printf("SE V[0x%X], V[0x%X]\n", x, y);
			if (cpu->reg.v[x] == cpu->reg.v[y]) {
				cpu->reg.pc += 2;
			} 			
			break;
		case 0x6000:
			printf("LD V[%X], 0x%.3X\n", x, byte);
			cpu->reg.v[x] = byte;
			break;
		case 0x7000:
			printf("ADD V[%X], 0x%.3X\n", x, byte);
			cpu->reg.v[x] += byte;
			break;
		case 0x8000:
			switch(opcode & 0x000F) {
				case 0x0:
					printf("LD V[%X], V[%X]\n", x, y);
					cpu->reg.v[x] = cpu->reg.v[y];
					break;
				case 0x1:
					printf("OR V[%X], V[%X]\n", x, y);
					cpu->reg.v[x] |= cpu->reg.v[y];
					break;
				case 0x2:
					printf("AND V[%X], V[%X]\n", x, y);
					cpu->reg.v[x] &= cpu->reg.v[y];
					break;
				case 0x3:
					printf("XOR V[%X], V[%X]\n", x, y);
					cpu->reg.v[x] ^= cpu->reg.v[y];
					break;
				case 0x4:
					printf("ADD V[%X], V[%X]\n", x, y);
					int overflow = (int)cpu->reg.v[x] + (int)cpu->reg.v[y];
					if (overflow >= UINT8_MAX) {
						cpu->reg.v[x] += cpu->reg.v[y];
						cpu->reg.v[0xF] = 1;
					} else {
						cpu->reg.v[x] += cpu->reg.v[y];
						cpu->reg.v[0xF] = 0;
					}
					break;
				case 0x5:
					printf("SUB V[%X], V[%X]\n", x, y);
					if (cpu->reg.v[x] >= cpu->reg.v[y]) {
						cpu->reg.v[x] -= cpu->reg.v[y];
						cpu->reg.v[0xF] = 1;
					} else {
						cpu->reg.v[x] -= cpu->reg.v[y];
						cpu->reg.v[0xF] = 0;
					}
					break;
				case 0x6: {
					printf("SHR V[%X], { V[%X] }\n", x, y);
					// QUIRK: cpu->reg.v[x] = cpu->reg.v[y];
					uint8_t bit = cpu->reg.v[x] & 0x01;
                    cpu->reg.v[x] >>= 1;
                    cpu->reg.v[0xF] = bit;
					break;
				}
				case 0x7:
					printf("SUB V[%X], V[%X]\n", y, x);
					if (cpu->reg.v[y] >= cpu->reg.v[x]) {
						cpu->reg.v[x] = cpu->reg.v[y] - cpu->reg.v[x];
						cpu->reg.v[0xF] = 1;
					} else {
						cpu->reg.v[x] = cpu->reg.v[y] - cpu->reg.v[x];
						cpu->reg.v[0xF] = 0;
					}
					break;
				case 0xE:
					printf("SHL V[%X], { V[%X] }\n", x, y);
					// QUIRK: cpu->reg.v[x] = cpu->reg.v[y];
					uint8_t bit = (cpu->reg.v[x] & 0x80) >> 7;
                    cpu->reg.v[x] <<= 1;
                    cpu->reg.v[0xF] = bit;
					break;
			}
			break;
		case 0x9000:
			printf("SNE V[%X], V[%X]\n", x, y);
			if (cpu->reg.v[x] != cpu->reg.v[y]) {
				cpu->reg.pc += 2;
			} 			
			break;
		case 0xA000:
			printf("LD I, 0x%.3X\n", addr);
			cpu->reg.i = addr;
			break;
		case 0xB000:
			cpu->reg.pc = addr + cpu->reg.v[0];
			break;
		case 0xC000:
			cpu->reg.v[x] = (rand() % 0xFF) & byte;
			break;
		case 0xD000:
			printf("DRW V[0x%X], V[0x%X], 0x%X\n", x, y, nibble);
			uint8_t sprite_x = cpu->reg.v[x] % 64;
			uint8_t sprite_y = cpu->reg.v[y] % 32;

			for (int height = 0; height < nibble; height++) {
				uint8_t sprite_row = cpu->memory.ram[cpu->reg.i + height];
				
				for (int width = 0; width <= 7; width++) {
					uint8_t pixel = (((sprite_row<<width) & 0x80) >> 7);
						
					cpu->display.fb[height + sprite_y][width + sprite_x] ^= pixel;
					if (cpu->display.fb[sprite_y][sprite_x] == 1) {
						cpu->reg.v[0xF]	= 0x1;
					} else { 
						cpu->reg.v[0xF]	= 0x0;
					}
				}
			}

			break;

		case 0xE000:
			switch (opcode & 0x00FF) {
				case 0x9E:
					printf("SKP V[%X]\n", x);
					if (keyboard_isset(cpu->keyboard, cpu->reg.v[x])) {
						cpu->reg.pc += 2;
					} 					
					break;
				case 0xA1:
					printf("SKNP V[%X]\n", x);
					if (!keyboard_isset(cpu->keyboard, cpu->reg.v[x])) {
						cpu->reg.pc += 2;
					} 
					break;
				default:
					printf("UNKP\n");
			}
			break;
		case 0xF000:
			switch(opcode & 0x00FF) {
				case 0x07:
					printf("LD V[%x], DT = %X\n", x, cpu->reg.delay_timer);
					cpu->reg.v[x] = cpu->reg.delay_timer;
					break;
				case 0x15:
					printf("LD DT, V[%x]\n", x);
					cpu->reg.delay_timer = cpu->reg.v[x];
					break;
				case 0x18:
					printf("LD ST, V[%x]\n", x);
					cpu->reg.sound_timer = cpu->reg.v[x];
					break;
				case 0x0A:
					printf("LD V[%X], K\n", x);

					bool key_down = false;
					for (uint8_t i = 0; i <= 0xF; i++) {
						if (keyboard_isset(cpu->keyboard, i)) {
							key_down = true;
						} else if (key_down) {
							cpu->reg.v[x] = i;
							cpu->reg.pc += 2;
							key_down = false;
						}
					}

					cpu->reg.pc -= 2;
					break;

				case 0x29:
					printf("LD F, V[%X]\n", x);
					cpu->reg.i = cpu->reg.v[x];
					break;
				case 0x1E:
					printf("ADD I, V[%X]\n", x);
					cpu->reg.i += cpu->reg.v[x];
					break;
				case 0x33:
					printf("LD B, V[%X]\n", x);
					cpu->memory.ram[cpu->reg.i] = (cpu->reg.v[x] / 100) % 10;
					cpu->memory.ram[cpu->reg.i + 1]= (cpu->reg.v[x] / 10) % 10;
					cpu->memory.ram[cpu->reg.i + 2] = (cpu->reg.v[x] / 1) % 10;
					break;
				case 0x55:
					printf("LD [%X], V[%X]\n", cpu->reg.i, x);
					for (int i = 0; i <= x; i++) {
						cpu->memory.ram[cpu->reg.i + i] = cpu->reg.v[i];
					}
					break;
				case 0x65:
					printf("LD [%X], V[%X]\n", cpu->reg.i, x);
					for (int i = 0; i <= x; i++) {
						cpu->reg.v[i] = cpu->memory.ram[cpu->reg.i + i];
					}
					break;
				default:
					printf("UNKP\n");
					break;
			}
			break;
		default:
			printf("UNKP\n");
			break;
		
	}

	if (cpu->reg.delay_timer > 0) {
		cpu->reg.delay_timer--;
	}

	if (cpu->reg.sound_timer > 0) {
		SDL_PauseAudioDevice(cpu->sound.id, 0);
		cpu->reg.sound_timer--;
	} else {
		SDL_PauseAudioDevice(cpu->sound.id, 1);
	}

}

long fsize(FILE *file) {
	fseek(file, 0, SEEK_END);
	long size = ftell(file);
	fseek(file, 0, SEEK_SET);
	return size;
}

uint8_t keyboard_to_hex(SDL_Event type) {
	switch (type.key.keysym.sym) {
		case SDLK_1:
			return 0x1;
		case SDLK_2: 
			return 0x2;
		case SDLK_3:
			return 0x3;
		case SDLK_4: 
			return 0xC;
		case SDLK_q: 
			return 0x4;
		case SDLK_w: 
			return 0x5;
		case SDLK_e:
			return 0x6;
		case SDLK_r: 
			return 0xD;
		case SDLK_a:
			return 0x7;
		case SDLK_s: 
			return 0x8;
		case SDLK_d:
			return 0x9;
		case SDLK_f: 
			return 0xE;
		case SDLK_z:
			return 0xA;
		case SDLK_x: 
			return 0x0;
		case SDLK_c:
			return 0xB;
		case SDLK_v: 
			return 0xF;
		default:
			return -1;
	}
}

int load_rom(const char *rom, cpu_t *cpu) {
	if (rom == NULL) {
		errno = ENOENT;
		return EXIT_FAILURE;
	}

	FILE *fp = fopen(rom, "r");
	if (fp == NULL) {
		return EXIT_FAILURE;
	}

	long size = fsize(fp);

	uint8_t bytes[size];
	const size_t count = fread(bytes, sizeof(bytes[0]), size, fp);

	if (count == size) {
		memory_write(&cpu->memory, ARRAY_LENGTH(bytes), bytes);
	} else {
		if (feof(fp)) {
			printf("Error reading %s: unexpected eof\n", rom);
		}
	}

	fclose(fp);

	return EXIT_SUCCESS;
}

int main(int argc, char *argv[]) {
	if (argc < 2) {
		printf("Usage: pixel8 ./ROM_FILE.c8\n");
		return EXIT_FAILURE;
	}

	cpu_t cpu = cpu_new();

	int error = load_rom(argv[1], &cpu);
	if (error == EXIT_FAILURE) {
		printf("Failed to load file");
	}

	if (SDL_Init(SDL_INIT_EVERYTHING) != 0) {
        printf("SDL_Init Error: %s\n", SDL_GetError());
        return 1;
    }

    SDL_Window* window = SDL_CreateWindow("CHIP 8", SDL_WINDOWPOS_CENTERED, SDL_WINDOWPOS_CENTERED, WINDOW_WIDTH, WINDOW_HEIGHT, 0);
    if (window == NULL) {
        printf("SDL_CreateWindow Error: %s\n", SDL_GetError());
        return 1;
    }


    SDL_Renderer* renderer = SDL_CreateRenderer(window, -1, SDL_RENDERER_ACCELERATED);
    if (renderer == NULL) {
        printf("SDL_CreateRenderer Error: %s\n", SDL_GetError());
        return 1;
    }

	SDL_RenderSetScale(renderer, WINDOW_WIDTH / FB_COLS, WINDOW_HEIGHT / FB_ROWS);

	SDL_AudioSpec *spec = audio_spec_new();

	cpu.sound.id = SDL_OpenAudioDevice(NULL, 0, spec, NULL, SDL_AUDIO_ALLOW_ANY_CHANGE);

	
	bool is_running = true;
	
	while (is_running) {
		SDL_Event event;
		while (SDL_PollEvent(&event)) {
			if (event.type == SDL_QUIT) {
				is_running = false;
			}

			uint8_t key = keyboard_to_hex(event);

			if (event.type == SDL_KEYUP) {
				keyboard_unsetkey(&cpu.keyboard, key);
			}

			if (event.type == SDL_KEYDOWN) {
				keyboard_setkey(&cpu.keyboard, key);
			}
		}

		uint16_t opcode = cpu_fetch(&cpu);

		cpu_decode(&cpu, opcode);

		SDL_RenderClear(renderer);
		
		for (int i = 0; i < FB_ROWS; i++) {
			for (int j = 0; j < FB_COLS; j++) {
				SDL_Rect rect = { .x = j , .y = i, .h = 10, .w = 10};

				if (cpu.display.fb[i][j] == 1) {
					SDL_SetRenderDrawColor(renderer, 255, 255, 255, 255);
					SDL_RenderFillRect(renderer, &rect);
				} else if (cpu.display.fb[i][j] == 0) {
					SDL_SetRenderDrawColor(renderer, 0, 0, 0, 255);
					SDL_RenderFillRect(renderer, &rect);
				}
			}
		}

		SDL_SetRenderDrawColor(renderer, 0, 0, 0, 255);

		SDL_RenderPresent(renderer);
	}

	free(spec);
	SDL_CloseAudioDevice(cpu.sound.id);
    SDL_DestroyRenderer(renderer);
    SDL_DestroyWindow(window);
    SDL_Quit();

	return EXIT_SUCCESS;
}
