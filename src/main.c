#include <stdio.h>

#define ARRAY_LEN(x) (sizeof(x) / sizeof(x[0]))

typedef struct memory_map memory_map;

struct memory_map {
        unsigned char ram[4095];
};

void memory_map_read_rom(memory_map m_map) {
        
}

void memory_map_write_rom(memory_map m_map) {

}

void main(void) {

    memory_map m_map = {};

}
