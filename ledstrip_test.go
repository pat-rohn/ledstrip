package ledstrip_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/pat-rohn/ledstrip"
	log "github.com/sirupsen/logrus"
)

func TestColorValue(t *testing.T) {
	log.SetLevel(log.InfoLevel)

	log.Traceln("TestSPI")

	emptyData := ledstrip.GetColorData(ledstrip.RGBPixel{0, 0, 0})

	fmt.Printf("%08b\n", emptyData)
	//fmt.Printf("%08b\n", oldEmptyData)

	colorData := []uint8{}
	oldColorData := []uint8{}
	leds := 1
	color := ledstrip.RGBPixel{7, 32, 128}

	for l := 0; l < leds; l++ {
		newData := ledstrip.GetColorData(color)
		for _, d := range newData {
			colorData = append(colorData, d)
		}
		//oldColorData = append(oldColorData,
		//	conn.DeprecatedgetTranslatedColor([3]uint8{color.Green, color.Red, color.Blue})...)
	}
	fmt.Printf("\n%08b\n", colorData)
	fmt.Printf("%08b\n\n", oldColorData)
	fmt.Printf("%0X\n", colorData)
	fmt.Printf("%0X\n", oldColorData)

	/* G / R / B
	[10010011 01001001 00100100   10010010 01001001 10100100  11010010 01001001 00100100]
	[10010011 01001001 00100100   10010010 01001001 10100100  11010010 01001001 00100100]
	*/
}

func TestPerformance(t *testing.T) {
	log.SetLevel(log.ErrorLevel)

	log.Traceln("TestSPI")

	emptyData := ledstrip.GetColorData(ledstrip.RGBPixel{0, 0, 0})

	fmt.Printf("%08b\n", emptyData)

	leds := 200
	startTime := time.Now()
	for i := 0; i < 500000; i++ {
		for l := 0; l < leds; l++ {
			ledstrip.GetColorData(ledstrip.RGBPixel{10, 20, 30})
		}
	}

	fmt.Printf("Time used %v for %d leds\n", time.Since(startTime)/500000, leds)

}
