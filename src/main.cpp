#include <array>
#include <cstdint>
#include <fstream>
#include <iostream>
#include <string>

const int MAX_MEM = 0xFFF;
const int MIN_MEM = 0x200;

class MemoryMap {
   public:
    std::array<uint8_t, MAX_MEM> ram;

    int load_file(std::string filename) {
        std::ifstream inf{"roms/" + filename, std::ios::in | std::ios::binary};
        if (!inf) {
            std::cout << "failed to read " << filename << " from disk\n";
            return 1;
        }

        int i = MIN_MEM;
        while (!inf.eof()) {
            uint8_t mem_byte{};
            // TODO: consider using ::get() to preserve bytes instead of
            // operator>> removing whitespace chars
            inf >> std::noskipws >> mem_byte;
            this->ram[i] = mem_byte;
            i++;
        }
        return 0;
    }
};

class CPU {
   public:
    uint16_t pc;
    uint8_t sp;

    std::array<uint8_t, 0xF> v_register;
    std::array<uint16_t, 16> stack_register;

    uint16_t i_register;

    uint8_t delay_timer_register;
    uint8_t sound_timer_register;

    uint16_t fetch_opcode(std::array<uint8_t, MAX_MEM> ram) {
        return (ram[this->pc] << 8) | (ram[this->pc + 1]);
    }
};

int main() {
    MemoryMap m_map = {0};
    CPU cpu = {MIN_MEM, 0, 0, 0, 0, 0, 0};

    m_map.load_file("IBMLOGO");

    for (;;) {
        uint16_t opcode = cpu.fetch_opcode(m_map.ram);

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
                        cpu.pc += 2;
                        break;
                    }

                    case 0xEE:
                        cpu.pc = cpu.stack_register[0xF];
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
                cpu.v_register[x] = nn;
                cpu.pc += 2;
                break;
            }
            case 0x7000: {
                std::cout << "LD V[" << std::hex << +x << "], " << std::hex
                          << +nn << "\n";
                cpu.v_register[x] += nn;
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
                cpu.i_register = nnn;
                cpu.pc += 2;
                break;
            }
            case 0xB000: {}
            case 0xC000: {}
            case 0xD000: {
                std::cout << "DRW Vx, Vy, nibble\n";
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
