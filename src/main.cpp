#include <array>
#include <cstdint>
#include <fstream>
#include <iostream>

#define MAX_MEM 4096
#define RESERVED 512

class MemoryMap {
   public:
    std::array<uint8_t, MAX_MEM> ram;

    int read_file(std::string filename) {
        std::ifstream inf{"roms/" + filename, std::ios::in};
        if (!inf) {
            std::cout << "failed to read " << filename << " from disk\n";
            return 1;
        }

        int i = RESERVED;
        while (inf) {
            uint8_t m_byte;
            inf >> m_byte;
            this->ram[i] = m_byte;
            i++;
        }
        inf.close();
        return 0;
    }
};

int main() {
    MemoryMap m_map = {0};

    m_map.read_file("IBMLOGO");
    std::cout << "[";
    for (int i{0}; i < m_map.ram.size(); i++) {
        std::cout << +m_map.ram[i] << " ";
    }
    std::cout << "]\n";
    return 0;
}
