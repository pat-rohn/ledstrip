package main

import (
	"fmt"
	"time"

	"github.com/pat-rohn/ledstrip"
	log "github.com/sirupsen/logrus"
)

func runExample0(c *ledstrip.ConnectionSPI, nrOfLeds int) {
	fmt.Println("example1")

	leds := []ledstrip.RGBPixel{{Red: 30, Blue: 6, Green: 10}}
	rDiff := -3
	gDiff := -2
	bDiff := 1
	max := uint8(70)
	min := uint8(5)
	for {
		time.Sleep(time.Millisecond * 20)
		oldColors := leds[0]
		if oldColors.Red > max {
			rDiff *= -1
		}
		if oldColors.Red <= min {
			rDiff *= -1
		}
		if oldColors.Green > max {
			gDiff *= -1
		}
		if oldColors.Green <= min {
			gDiff *= -1
		}
		if oldColors.Blue > max {
			bDiff *= -1
		}
		if oldColors.Blue <= min {
			bDiff *= -1
		}
		newColor := ledstrip.RGBPixel{
			Red:   uint8(int(oldColors.Red) + rDiff),
			Green: uint8(int(oldColors.Green) + gDiff),
			Blue:  uint8(int(oldColors.Blue) + bDiff),
		}
		fmt.Printf("newColor: %d %d %d\n", newColor.Red, newColor.Green, newColor.Blue)
		for i := range leds {
			leds[i] = newColor
		}
		if len(leds) < nrOfLeds {
			leds = append(leds, newColor)
		}
		c.Render(leds)
	}

}

func runExample1(c *ledstrip.ConnectionSPI, nrOfLeds int) {
	fmt.Println("example1")
	ledsWorms := createExample1(nrOfLeds)

	runner := ledstrip.FaderRunner{
		Conn: c,
	}
	for {
		runner.RunLEDS(ledsWorms, 0.5)
	}
}

func runExample2(c *ledstrip.ConnectionSPI, nrOfLeds int) {

	fmt.Println("example2")
	example := createExample2(nrOfLeds)
	runner := ledstrip.FaderRunner{
		Conn: c,
	}

	for {
		runner.RunLEDS(example, 1)
	}
}

func runExample3(c *ledstrip.ConnectionSPI, nrOfLeds int, maskLength int, color1 ledstrip.RGBPixel, color2 ledstrip.RGBPixel) {
	fmt.Println("example 3: Color Fading")

	leds := createExample3(nrOfLeds, maskLength, color1, color2)

	runner := ledstrip.FaderRunner{
		Conn: c,
	}

	for {
		runner.RunLEDS(leds, 1)
	}
}

func runExample4(c *ledstrip.ConnectionSPI, nrOfLeds int, maskLength int, color1 ledstrip.RGBPixel, color2 ledstrip.RGBPixel) {
	fmt.Println("example 4: Color Fading ")

	leds := createExample3(nrOfLeds, maskLength, color1, color2)
	go func() {
		runner := ledstrip.FaderRunner{
			Conn: c,
		}

		for {
			runner.RunLEDSSplit(leds, 2, false)
		}
	}()
	for {
		time.Sleep(time.Second * 1)
		fmt.Printf(" ... ")
	}
}

func createExample1(nrOfLeds int) []ledstrip.RGBPixel {
	logFields := log.Fields{"func": "CreateExample1"}

	var leds []ledstrip.RGBPixel
	colorValues := [10]uint8{uint8(0), uint8(2), uint8(4), uint8(4), uint8(8), uint8(8), uint8(16), uint8(32), uint8(32), uint8(64)}
	for len(leds) < nrOfLeds {
		for i := 0; i < 10; i++ {
			rVal := colorValues[ledstrip.Pmod(i, 10)]
			gVal := colorValues[0]
			bVal := colorValues[0]
			leds = append(leds, ledstrip.RGBPixel{
				Red:   rVal,
				Green: gVal,
				Blue:  bVal,
			})
		}
		for i := 0; i < 10; i++ {
			rVal := colorValues[0]
			gVal := colorValues[ledstrip.Pmod(i, 10)]
			bVal := colorValues[0]
			leds = append(leds, ledstrip.RGBPixel{
				Red:   rVal,
				Green: gVal,
				Blue:  bVal,
			})

		}
		for i := 0; i < 10; i++ {
			rVal := colorValues[0]
			gVal := colorValues[0]
			bVal := colorValues[ledstrip.Pmod(i, 10)]
			leds = append(leds, ledstrip.RGBPixel{
				Red:   rVal,
				Green: gVal,
				Blue:  bVal,
			})
		}
	}
	leds = leds[:nrOfLeds]
	log.WithFields(logFields).Infof("Example with %d LEDs created", len(leds))
	return leds
}

func createExample2(nrOfLeds int) []ledstrip.RGBPixel {
	logFields := log.Fields{"func": "CreateExample2"}
	log.WithFields(logFields).Traceln("Example 2")

	var leds []ledstrip.RGBPixel
	colorValues := [10]uint8{uint8(0), uint8(4), uint8(16), uint8(32), uint8(64), uint8(32), uint8(16), uint8(8), uint8(4), uint8(0)}
	for len(leds) < nrOfLeds {
		for i := 0; i < 10; i++ {
			rVal := colorValues[ledstrip.Pmod(i, 10)/2]
			gVal := colorValues[ledstrip.Pmod(i, 10)]
			bVal := colorValues[0]
			leds = append(leds, ledstrip.RGBPixel{
				Red:   rVal,
				Green: gVal,
				Blue:  bVal,
			})
		}
		for i := 0; i < 10; i++ {
			rVal := colorValues[0]
			gVal := colorValues[ledstrip.Pmod(i, 10)/2]
			bVal := colorValues[ledstrip.Pmod(i, 10)]
			leds = append(leds, ledstrip.RGBPixel{
				Red:   rVal,
				Green: gVal,
				Blue:  bVal,
			})
		}
		for i := 0; i < 10; i++ {
			rVal := colorValues[ledstrip.Pmod(i, 10)]
			gVal := colorValues[0]
			bVal := colorValues[ledstrip.Pmod(i, 10)/2]
			leds = append(leds, ledstrip.RGBPixel{
				Red:   rVal,
				Green: gVal,
				Blue:  bVal,
			})
		}
	}
	leds = leds[:nrOfLeds]
	log.WithFields(logFields).Infof("Example with %d LEDs created", len(leds))
	return leds
}

func createExample3(nrOfLeds int, maskLength int, color1 ledstrip.RGBPixel, color2 ledstrip.RGBPixel) []ledstrip.RGBPixel {
	leds := make([]ledstrip.RGBPixel, nrOfLeds)

	if nrOfLeds%maskLength != 0 {
		log.Warnf("Length should be divisible with mask length %d %d", nrOfLeds, maskLength)
	}
	mask := make([]float32, maskLength)

	factor := 2.0 / float32(maskLength)
	for i := range mask {
		if i < maskLength/2 {
			mask[i] = float32(i) * factor
		} else {
			mask[i] = (float32(maskLength) - float32(i+1)) * factor
			if mask[i] < 0 {
				mask[i] = 0
			}
		}
	}
	for i := range leds {
		leds[i%nrOfLeds].Red = uint8(mask[i%maskLength] * float32(color1.Red))
		leds[i%nrOfLeds].Green = uint8(mask[i%maskLength] * float32(color1.Green))
		leds[i%nrOfLeds].Blue = uint8(mask[i%maskLength] * float32(color1.Blue))
	}

	for i := range leds {
		maskIndex := (i + maskLength/2) % maskLength
		red := float32(leds[i].Red) + float32(color2.Red)*mask[maskIndex]
		if red > 255 {
			red = 255
		}
		green := float32(leds[i].Green) + float32(color2.Green)*mask[maskIndex]
		if green > 255 {
			green = 255
		}
		blue := float32(leds[i].Blue) + float32(color2.Blue)*mask[maskIndex]
		if blue > 255 {
			blue = 255
		}
		leds[i].Red = uint8(red)
		leds[i].Green = uint8(green)
		leds[i].Blue = uint8(blue)
	}
	return leds
}
