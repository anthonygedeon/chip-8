#include <array>
#include <cstdint>
#include <fstream>
#include <iostream>
#include <string>

#define MAX_MEM 0xFFF
#define MIN_MEM 0x200

class MemoryMap {
   public:
    std::array<uint8_t, MAX_MEM> ram;

    int load_file(std::string filename) {
        std::ifstream inf{"roms/" + filename, std::ios::in};
        if (!inf) {
            std::cout << "failed to read " << filename << " from disk\n";
            return 1;
        }

        int i = MIN_MEM;
        while (inf) {
            uint8_t mem_byte{};
            inf >> mem_byte;
            this->ram[i] = mem_byte;
            i++;
        }
        return 0;
    }
};

class CPU {
    void fetch_opcode() {}

   public:
    uint16_t pc;
    uint8_t sp;

    std::array<uint8_t, 0xF> v_register;
    std::array<uint16_t, 16> stack_register;

    uint16_t i_register;

    uint8_t delay_timer_register;
    uint8_t sound_timer_register;
};

int main() {
    MemoryMap m_map = { 0 };
    CPU cpu = { MIN_MEM, 0, 0, 0, 0, 0, 0 };

    m_map.load_file("IBMLOGO");

    for (;;) {
        uint16_t opcode = (m_map.ram[cpu.pc] << 8) | (m_map.ram[cpu.pc + 1]);

        uint8_t x = opcode & 0x0F00;
        uint8_t y = opcode & 0x00F0;
        uint8_t n = opcode & 0x000F;
        uint8_t nn = opcode & 0x00FF;
        uint16_t nnn = opcode & 0x0FFF;

        std::cout << std::hex << opcode << "\n";

        switch (opcode & 0xF000) {
            case 0x0000: {
                switch (opcode & 0x00FF) {
                    case 0xE0: {
                        std::cout << "cls\n";
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
            case 0x6000: {
                std::cout << "LD V[" << std::hex <<+x << "], " << std::hex << +nn << "\n";
                cpu.v_register[x] = nn;
                cpu.pc += 2;
                break;
            }
            case 0x7000: {
                std::cout << "LD V[" << std::hex <<+x << "], " << std::hex << +nn << "\n";
                cpu.v_register[x] += nn;
                cpu.pc += 2;
                break;
            }
            case 0xA000: {
                std::cout << "LD I, " << std::hex << +nnn << "\n";
                cpu.i_register = nnn;
                cpu.pc += 2;
                break;
            }
            case 0xD000: {
                std::cout << "DRW Vx, Vy, nibble\n";
                cpu.pc += 2;
                break;
            }
        }
    }

    return 0;
}
