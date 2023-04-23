package ledstrip

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

type FaderRunner struct {
	NrOfLeds int
	Colors   []RGBPixel
	Conn     *ConnectionSPI
}

func (f *FaderRunner) RunLEDS(leds []RGBPixel, speedFactor float32) {
	logFields := log.Fields{"func": "RunLEDS", "SpeedFactor": speedFactor}
	log.WithFields(logFields).Traceln("RunLEDS")
	ledsTarget := make([]RGBPixel, len(leds))
	ledsOld := make([]RGBPixel, len(leds))

	stepSize := speedFactor / 15
	sleepTime := time.Duration(10/speedFactor) * time.Millisecond

	logFields["stepSize"] = stepSize
	logFields["sleepTime"] = sleepTime
	log.WithFields(logFields).Infoln("Start Fader")
	for {

		copy(ledsTarget, leds)
		copy(ledsOld, leds)
		ledsTarget = PlaceInFront(ledsTarget, ledsTarget[len(ledsTarget)-1])
		for fadeProgress := float32(0.0); fadeProgress <= 1; fadeProgress += stepSize {
			for i := range leds {
				redShift := (float32(ledsTarget[i].Red) - float32(ledsOld[i].Red)) * fadeProgress
				greenShift := (float32(ledsTarget[i].Green) - float32(ledsOld[i].Green)) * fadeProgress
				blueShift := (float32(ledsTarget[i].Blue) - float32(ledsOld[i].Blue)) * fadeProgress

				if i == len(leds)-1 && false {
					fmt.Printf(" fadeFactor %v \n", fadeProgress)
					fmt.Printf("Diffs %v / %v / %v \n",
						float32(ledsTarget[i].Red)-float32(ledsOld[i].Red),
						float32(ledsTarget[i].Green)-float32(ledsOld[i].Green),
						float32(ledsTarget[i].Blue)-float32(ledsOld[i].Blue))
					fmt.Printf("Shifts %v / %v / %v \n", redShift, greenShift, blueShift)
					fmt.Printf(" --- LED old %v  ---- LED Target %v\n", leds[i], ledsTarget[i])
				}

				r := float32(ledsOld[i].Red) + redShift
				g := float32(ledsOld[i].Green) + greenShift
				b := float32(ledsOld[i].Blue) + blueShift
				leds[i] = RGBPixel{
					Red:   uint8(r),
					Green: uint8(g),
					Blue:  uint8(b),
				}

			}

			f.Conn.Render(leds)
			time.Sleep(sleepTime)

		}
		copy(leds, ledsTarget)

	}
}
