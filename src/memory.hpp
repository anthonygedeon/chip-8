#if !defined(H_MEMORY_HPP)
#define H_MEMORY_HPP

#include <algorithm>
#include <array>
#include <fstream>
#include <iostream>
#include <string>

extern const int max_mem = 0xFFF;
extern const int min_mem = 0x200;

namespace memory {

std::array<uint8_t, max_mem> ram{0};
std::array<uint8_t, 63*31> vram{0};

constexpr std::array<uint8_t, 110> fonts{
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
    0xF0, 0x80, 0xF0, 0x80, 0x80  // F
};

int load_rom(std::string filename) {
    std::ifstream inf{"roms/" + filename, std::ios::in | std::ios::binary};
    if (!inf) {
        std::cout << "failed to read " << filename << " from disk\n";
        return 1;
    }

    std::for_each(ram.begin(), ram.end(),
                  [idx = min_mem, &inf](int i) mutable {
                      if (!inf.eof()) {
                          auto byte{inf.get()};
                          ram[idx] = byte;
                          ++idx;
                      }
                  });
    return 0;
}

int load_font() {
    for (auto i = 0u; i < fonts.size(); ++i) {
        ram[i+0x50] = fonts[i];
    }
    return 0;
}

};  // namespace memory

#endif
