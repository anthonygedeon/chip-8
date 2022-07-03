#include "memory.hpp"
#include "registers.hpp"

#include <array>
#include <cstdint>
#include <iostream>

struct CPU {
    uint16_t pc;
    uint8_t sp;
        
    Register reg;

    uint8_t display[31][63];

    uint16_t fetch_opcode(std::array<uint8_t, max_mem> ram) {
        return (ram[this->pc] << 8) | (ram[this->pc + 1]);
    }

    CPU() {
        this->pc = min_mem;
        this->sp = 0;
        this->reg.i = 0;
        this->reg.delay_timer = 0;
        this->reg.sound_timer = 0;
    }
};

int main() {
    auto cpu = CPU();
    memory::load_rom("IBMLOGO");
    memory::load_font();

    for (;;) {
        uint16_t opcode = cpu.fetch_opcode(memory::ram);

        uint8_t x = (opcode & 0x0F00) >> 8;
        uint8_t y = (opcode & 0x00F0) >> 4;
        uint8_t n = opcode & 0x000F;
        uint8_t nn = opcode & 0x00FF;
        uint16_t nnn = opcode & 0x0FFF;

        switch (opcode & 0xF000) {
            case 0x0000: {
                switch (opcode & 0x00FF) {
                    case 0xE0: {
                        std::cout << "CLS\n";
                        memory::vram.fill({0});
                        cpu.pc += 2;
                        break;
                    }

                    case 0xEE:
                        cpu.pc = cpu.reg.s[0xF];
                        cpu.sp--;
                        std::cout << "RET\n";
                        break;
                }
                break;
            }
            case 0x1000: {
                std::cout << "JP " << std::hex << +nnn << "\n";
                cpu.pc = nnn;
                break;
            }
            case 0x2000: {}
            case 0x3000: {}
            case 0x4000: {}
            case 0x5000: {}
            case 0x6000: {
                std::cout << "LD V[" << std::hex << +x << "], " << std::hex
                          << +nn << "\n";
                cpu.reg.v[x] = nn;
                cpu.pc += 2;
                break;
            }
            case 0x7000: {
                std::cout << "LD V[" << std::hex << +x << "], " << std::hex
                          << +nn << "\n";
                cpu.reg.v[x] += nn;
                cpu.pc += 2;
                break;
            }
            case 0x8000: {
                switch(opcode & 0x000F) {
                        case 0x0: {} 
                        case 0x1: {}
                        case 0x2: {}
                        case 0x3: {}
                        case 0x4: {}
                        case 0x5: {}
                        case 0x6: {}
                        case 0x7: {}
                        case 0xE: {}
                }
            }
            case 0x9000: {}
            case 0xA000: {
                std::cout << "LD I, " << std::hex << +nnn << "\n";
                cpu.reg.i = nnn;
                cpu.pc += 2;
                break;
            }
            case 0xB000: {}
            case 0xC000: {}
            case 0xD000: {
                std::cout << "DRW Vx, Vy, nibble\n";
                
                uint8_t addr  = cpu.reg.i;
                uint8_t x_pos = cpu.reg.v[x];
                uint8_t y_pos = cpu.reg.v[y];
                
                //cpu.v_register[0xF] = 1;
                //cpu.v_register[0xF] = 0;

                cpu.pc += 2;
                break;
            }
            case 0xE000: {
                switch(opcode & 0x00FF) {
                        case 0x9E: {}
                        case 0xA1: {}
                }
            }
            case 0xF000: {
                switch(opcode & 0x00FF) {
                    case 0x07: {}
                    case 0x0A: {}
                    case 0x15: {}
                    case 0x18: {}
                    case 0x1E: {}
                    case 0x29: {}
                    case 0x33: {}
                    case 0x55: {}
                    case 0x65: {}
                }
            }
        }
    }

    return 0;
}
