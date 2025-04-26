package ledMatrix

import (
	"github.com/mcuadros/go-rpi-rgb-led-matrix"
)

type LedMatrixConfig struct {
	Rows int
	Cols int
	Parallel int
	Chain int
	Brightness int
	HardwareMapping string
	ShowRefresh bool
	inverse_colors bool
	disableHardwarePulsing bool
}

func NewLedMatrix() (*LedMatrix, error) {
	led := LedMatrix{}

	led.Rows = 32
	led.Cols = 64

	return led	
}

func(*l LedMatrix) PrintMatrix()  error {
	if err := rpio.Open(); err != nil {
		return err
	}
	defer rpio.Close()
}
