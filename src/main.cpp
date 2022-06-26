#include <array>
#include <cstdint>
#include <fstream>
#include <iostream>

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
    
    void fetch_opcode() {
        
    }

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

    uint16_t opcode = (m_map.ram[cpu.pc] << 8) | (m_map.ram[cpu.pc+1]);
    std::cout << cpu.pc << "\n";
    std::cout << opcode << "\n";
    //std::cout << "[";
    //for (int i{0}; i < m_map.ram.size(); i++) {
        //std::cout << +m_map.ram[i] << " ";
    //}
    //std::cout << "]\n";

    

    return 0;
}
