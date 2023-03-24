package ledstrip

const (
	logPkg string = "ledstrip"
)

// RGBPixel has one byte each color.
type RGBPixel struct {
	Red   uint8
	Green uint8
	Blue  uint8
}

func PlaceInFront(leds []RGBPixel, led RGBPixel) []RGBPixel {
	newLeds := append([]RGBPixel{led}, leds...)
	leds = newLeds[:len(leds)]
	return leds
}

func PlaceAtBack(leds []RGBPixel, led RGBPixel) []RGBPixel {
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
