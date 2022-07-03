#if !defined(H_REGISTERS_HPP)
#define H_REGISTERS_HPP

#include <array>
#include <cstdint>

struct Register {
    std::array<uint8_t, 0xF> v;
    std::array<uint16_t, 16> s;

    uint16_t i;

    uint8_t delay_timer;
    uint8_t sound_timer;

};

#endif
