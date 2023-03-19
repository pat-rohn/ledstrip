package ledstrip

import (
	log "github.com/sirupsen/logrus"
)

const (
	logPkg string = "ledstrip"
)

// RGBPixel has one byte each color.
type RGBPixel struct {
	Red   uint8
	Green uint8
	Blue  uint8
}

// RenderLEDs translates RGBPixels into SPI message and transfers the message
func (conn *ConnectionSPI) RenderLEDs(pixels []RGBPixel) {
	logFields := log.Fields{"package": logPkg, "conn": "SPI", "func": "RenderLEDs"}
	log.WithFields(logFields).Infof("RenderLEDs with len %v", len(pixels))

	var translatedRGBs []uint8
	for _, pixel := range pixels {
		colorData := conn.GetColorData(pixel)

		for _, c := range colorData {
			translatedRGBs = append(translatedRGBs, c)
		}
	}

	log.Tracef("%08b", translatedRGBs)
	conn.transfer(translatedRGBs)
}

func PlaceInFront(leds []RGBPixel, led RGBPixel) []RGBPixel {
	newLeds := append([]RGBPixel{led}, leds...)
	leds = newLeds[:len(leds)]
	return leds
}

func PlaceInBack(leds []RGBPixel, led RGBPixel) []RGBPixel {
	newLeds := append(leds, []RGBPixel{led}...)
	leds = newLeds[1 : len(leds)+1]
	return leds
}

func Inverse(leds []RGBPixel) []RGBPixel {
	var newLeds []RGBPixel
	for _, led := range leds {
		newLeds = append([]RGBPixel{led}, newLeds...)
	}
	return newLeds
}

func Pmod(a, b int) int {
	m := a % b
	if a < 0 && b < 0 {
		m -= b
	}
	if a < 0 && b > 0 {
		m += b
	}
	return m
}
