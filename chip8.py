import pygame

pygame.init()


WIDTH = 500 
HEIGHT = 300
screen = pygame.display.set_mode([WIDTH, HEIGHT])

class CPU:
    def __init__(self):
        self.program_counter = 0x000

    def INST_CLR(self):
        pass

    def INST_JMP(self):
        pass

    def set_VX(self):
        pass

    def set_reg_I(self):
        pass

    def INST_DRW(self):
        pass
