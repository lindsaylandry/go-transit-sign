package ledMatrix

import (
	"github.com/tfk1410/go-rpi-rgb-led-matrix"
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

type LedMatrix struct {

}

func NewLedMatrix() (*LedMatrix, error) {
	led := LedMatrix{}

	led.Rows = 32
	led.Cols = 64

	return &led	
}

func(*l LedMatrix) PrintMatrix()  error {
	if err := rpio.Open(); err != nil {
		return err
	}
	defer rpio.Close()
}

func(*l LedMatrix) PrintTitle(title string) error {

}
