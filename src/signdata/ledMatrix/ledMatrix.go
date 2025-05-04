package ledMatrix

import (
	"image/color"

	"github.com/tfk1410/go-rpi-rgb-led-matrix"
)

func WriteToMatrix() error {
	config := &rgbmatrix.DefaultConfig
	config.Rows = 32
	config.Cols = 64
	config.Parallel = 1
	config.ChainLength = 1
	config.Brightness = 50
	config.HardwareMapping = "adafruit-hat"
	config.ShowRefreshRate = false
	config.InverseColors = false
	config.DisableHardwarePulsing = false

	m, err := rgbmatrix.NewRGBLedMatrix(config)
	if err != nil {
		return err
	}

	c := rgbmatrix.NewCanvas(m)
	defer c.Close()

	bounds := c.Bounds()
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			c.Set(x, y, color.RGBA{255, 0, 0, 255})
			c.Render()
		}
	}

	return nil
}
