import pygame

pygame.init()


WIDTH = 500 
HEIGHT = 300
screen = pygame.display.set_mode([WIDTH, HEIGHT])

class CPU:
    def __init__(self):
        self.pc = 0x000

        self.I = 0

    def disassemble(self, rom_file):
        with open("roms/" + rom_file, "r") as rom:
            print(rom)


    def INST_CLR(self):
        pass

    def INST_JMP(self):
        pass

    def set_VX(self):
        pass

    def set_I(self):
        pass

    def INST_DRW(self):
        pass


processor = CPU()

print(processor.disassemble("IBM_Logo.ch8"))
