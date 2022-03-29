#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#define ARRAY_LEN(x) (sizeof(x) / sizeof(x[0]))

#define RESERVED_OFFSET 512

typedef struct memory_map memory_map;

struct memory_map {
        
        const char* filename;

        unsigned char ram[4095];
};

memory_map memory_map_new(const char* rom_path) {
       memory_map m_map = { rom_path, 0 }; 
       return m_map;
}

void memory_map_write_rom(memory_map* memory, int index, unsigned char opcode) {
        memory->ram[index+RESERVED_OFFSET] = opcode;
}

void memory_map_read_rom(memory_map memory) {
    FILE* file_ptr = fopen(memory.filename, "r");

    if (file_ptr == NULL) {
        printf("file can't be opened");
        exit(-1);
    }

    unsigned char opcode;
    for (int i = 0; !feof(file_ptr); i++) {
        opcode = fgetc(file_ptr);
        memory_map_write_rom(&m_map, i, opcode);
    }

    fclose(file_ptr);
}

void main(void) {
    memory_map mem = memory_map_new("roms/IBMLOGO");

    printf("[");

    for (int i = 0; i < ARRAY_LEN(mem.ram); i++) {
        printf("%x\n", mem.ram[i]);
    }
    printf(" ]");
}
