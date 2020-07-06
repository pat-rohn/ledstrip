package ledstrip

import (
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	logPkg string = "ledstrip"
)

//RGBPixel has one byte each color.
type RGBPixel struct {
	Red   uint8
	Green uint8
	Blue  uint8
}

//CreateWorms creats a slice of RGBPixels
func CreateWorms() []RGBPixel {
	logFields := log.Fields{"package": logPkg, "func": "Test"}
	log.WithFields(logFields).Traceln("Test")

	var leds []RGBPixel
	colorValues := [10]uint8{uint8(0), uint8(2), uint8(4), uint8(4), uint8(8), uint8(8), uint8(16), uint8(32), uint8(32), uint8(64)}

	for i := 0; i < 10; i++ {
		rVal := colorValues[pmod(i, 10)]
		gVal := colorValues[0]
		bVal := colorValues[0]
		leds = append(leds, RGBPixel{
			Red:   rVal,
			Green: gVal,
			Blue:  bVal,
		})
	}
	for i := 0; i < 10; i++ {
		rVal := colorValues[0]
		gVal := colorValues[pmod(i, 10)]
		bVal := colorValues[0]
		leds = append(leds, RGBPixel{
			Red:   rVal,
			Green: gVal,
			Blue:  bVal,
		})

	}
	for i := 0; i < 10; i++ {
		rVal := colorValues[0]
		gVal := colorValues[0]
		bVal := colorValues[pmod(i, 10)]
		leds = append(leds, RGBPixel{
			Red:   rVal,
			Green: gVal,
			Blue:  bVal,
		})

	}
	return leds
}

//CreateTest creats a slice of RGBPixels
func CreateTest() []RGBPixel {
	logFields := log.Fields{"package": logPkg, "func": "Test1"}
	log.WithFields(logFields).Traceln("Test")

	var leds []RGBPixel
	colorValues := [10]uint8{uint8(0), uint8(4), uint8(16), uint8(32), uint8(64), uint8(32), uint8(16), uint8(8), uint8(4), uint8(0)}

	for i := 0; i < 10; i++ {
		rVal := colorValues[pmod(i, 10)/2]
		gVal := colorValues[pmod(i, 10)]
		bVal := colorValues[0]
		leds = append(leds, RGBPixel{
			Red:   rVal,
			Green: gVal,
			Blue:  bVal,
		})
	}
	for i := 0; i < 10; i++ {
		rVal := colorValues[0]
		gVal := colorValues[pmod(i, 10)/2]
		bVal := colorValues[pmod(i, 10)]
		leds = append(leds, RGBPixel{
			Red:   rVal,
			Green: gVal,
			Blue:  bVal,
		})
	}
	for i := 0; i < 10; i++ {
		rVal := colorValues[pmod(i, 10)]
		gVal := colorValues[0]
		bVal := colorValues[pmod(i, 10)/2]
		leds = append(leds, RGBPixel{
			Red:   rVal,
			Green: gVal,
			Blue:  bVal,
		})
	}
	return leds
}

//RunLEDS lets LEDs move for time given
func (conn *ConnectionSPI) RunLEDS(leds []RGBPixel, runTime time.Duration) {
	logFields := log.Fields{"package": logPkg, "func": "RunLEDS"}
	log.WithFields(logFields).Traceln("RunLEDS")
	runTime = runTime / 4
	endTime := time.Now().Add(runTime)
	for time.Now().Before(endTime) {
		conn.RenderLEDs(leds)
		leds = placeInFront(leds, leds[len(leds)-1])
		waitTime := endTime.Sub(time.Now()) / 100
		time.Sleep(waitTime)

	}

	endTime = time.Now().Add(runTime)
	for time.Now().Before(endTime) {
		conn.RenderLEDs(leds)
		leds = placeInFront(leds, leds[len(leds)-1])
		waitTime := (runTime - endTime.Sub(time.Now())) / 100
		time.Sleep(waitTime)
	}

	leds = inverse(leds)
	endTime = time.Now().Add(runTime)
	for time.Now().Before(endTime) {
		conn.RenderLEDs(leds)
		leds = placeInBack(leds, leds[0])
		waitTime := endTime.Sub(time.Now()) / 100
		time.Sleep(waitTime)

	}
	endTime = time.Now().Add(runTime)
	for time.Now().Before(endTime) {
		conn.RenderLEDs(leds)
		leds = placeInBack(leds, leds[0])
		waitTime := (runTime - endTime.Sub(time.Now())) / 100
		time.Sleep(waitTime)
	}
}

// Clear switches off all LEDs
func (conn *ConnectionSPI) Clear(nrLeds int) {
	logFields := log.Fields{"package": logPkg, "func": "Clear"}
	log.WithFields(logFields).Traceln("Clear")

	var clearArray []uint8

	for i := 0; i < nrLeds; i++ {
		clearArray = append(clearArray, conn.getTranslatedColor([3]uint8{0, 0, 0})...)
	}
	conn.transfer(clearArray)
}

func placeInFront(leds []RGBPixel, led RGBPixel) []RGBPixel {
	newLeds := append([]RGBPixel{led}, leds...)
	leds = newLeds[:len(leds)]
	return leds
}

func placeInBack(leds []RGBPixel, led RGBPixel) []RGBPixel {
	newLeds := append(leds, []RGBPixel{led}...)
	leds = newLeds[1 : len(leds)+1]
	return leds
}

func inverse(leds []RGBPixel) []RGBPixel {
	var newLeds []RGBPixel
	for _, led := range leds {
		newLeds = append([]RGBPixel{led}, newLeds...)
	}
	return newLeds
}

func pmod(a, b int) int {
	m := a % b
	if a < 0 && b < 0 {
		m -= b
	}
	if a < 0 && b > 0 {
		m += b
	}
	return m
}
