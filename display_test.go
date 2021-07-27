package main

import "testing"

func TestClear(t *testing.T) {

	display := Display{}

	display.Clear()

	for i := 0; i < len(display.gfx); i++ {
		for j := 0; j < len(display.gfx[i]); j++ {
			if display.gfx[i][j] != 0 {
				t.Errorf("wanted %[1]q, got %[1]q", display.gfx)
			}
		}
	}

}