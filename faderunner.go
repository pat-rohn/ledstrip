package ledstrip

import (
	"time"

	log "github.com/sirupsen/logrus"
)

type FaderRunner struct {
	Conn *ConnectionSPI
}

func (f *FaderRunner) RunLEDS(leds []RGBPixel, speedFactor float32) {
	logFields := log.Fields{"func": "RunLEDS", "SpeedFactor": speedFactor}
	log.WithFields(logFields).Traceln("RunLEDS")
	ledsTarget := make([]RGBPixel, len(leds))
	ledsOld := make([]RGBPixel, len(leds))

	stepSize := speedFactor / 15
	sleepTime := time.Duration(10/speedFactor) * time.Millisecond

	log.WithFields(logFields).Infoln("Start Fader")
	copy(ledsTarget, leds)

	for {
		ledsTarget = PlaceInFront(ledsTarget, ledsTarget[len(ledsTarget)-1])
		for fadeProgress := float32(0.0); fadeProgress <= 1; fadeProgress += stepSize {
			for i := range leds {
				redShift := (float32(ledsTarget[i].Red) - float32(ledsOld[i].Red)) * fadeProgress
				greenShift := (float32(ledsTarget[i].Green) - float32(ledsOld[i].Green)) * fadeProgress
				blueShift := (float32(ledsTarget[i].Blue) - float32(ledsOld[i].Blue)) * fadeProgress

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
		copy(ledsOld, leds)
	}
}

func (f *FaderRunner) RunLEDSSplit(leds []RGBPixel, speedFactor float32, invert bool) {
	logFields := log.Fields{"func": "RunLEDS", "SpeedFactor": speedFactor}
	log.WithFields(logFields).Traceln("RunLEDS")
	middle := len(leds) / 2

	stepSize := speedFactor / 15
	sleepTime := time.Duration(10/speedFactor) * time.Millisecond

	ledsTargetLeft := make([]RGBPixel, middle)
	ledsOldLeft := make([]RGBPixel, middle)
	ledsTargetRight := make([]RGBPixel, middle)
	ledsOldRight := make([]RGBPixel, middle)

	copy(ledsTargetLeft, leds[:middle])
	copy(ledsOldLeft, leds[:middle])
	copy(ledsTargetRight, leds[middle:])
	copy(ledsOldRight, leds[middle:])

	log.WithFields(logFields).Infoln("Start faderunner")

	for {
		if invert {
			ledsTargetLeft = PlaceInFront(ledsTargetLeft, ledsTargetLeft[len(ledsTargetLeft)-1])
		} else {
			ledsTargetLeft = PlaceAtBack(ledsTargetLeft, ledsTargetLeft[0])
		}
		if invert {
			ledsTargetRight = PlaceAtBack(ledsTargetRight, ledsTargetRight[0])
		} else {
			ledsTargetRight = PlaceInFront(ledsTargetRight, ledsTargetRight[len(ledsTargetRight)-1])
		}
		for fadeProgress := float32(0.0); fadeProgress <= 1; fadeProgress += stepSize {
			for i := 0; i < middle; i++ {
				redShift := (float32(ledsTargetLeft[i].Red) - float32(ledsOldLeft[i].Red)) * fadeProgress
				greenShift := (float32(ledsTargetLeft[i].Green) - float32(ledsOldLeft[i].Green)) * fadeProgress
				blueShift := (float32(ledsTargetLeft[i].Blue) - float32(ledsOldLeft[i].Blue)) * fadeProgress

				r := float32(ledsOldLeft[i].Red) + redShift
				g := float32(ledsOldLeft[i].Green) + greenShift
				b := float32(ledsOldLeft[i].Blue) + blueShift
				leds[i] = RGBPixel{
					Red:   uint8(r),
					Green: uint8(g),
					Blue:  uint8(b),
				}
			}
			for i := 0; i < middle; i++ {
				redShift := (float32(ledsTargetRight[i].Red) - float32(ledsOldRight[i].Red)) * fadeProgress
				greenShift := (float32(ledsTargetRight[i].Green) - float32(ledsOldRight[i].Green)) * fadeProgress
				blueShift := (float32(ledsTargetRight[i].Blue) - float32(ledsOldRight[i].Blue)) * fadeProgress

				r := float32(ledsOldRight[i].Red) + redShift
				g := float32(ledsOldRight[i].Green) + greenShift
				b := float32(ledsOldRight[i].Blue) + blueShift
				leds[i+middle] = RGBPixel{
					Red:   uint8(r),
					Green: uint8(g),
					Blue:  uint8(b),
				}
			}
			f.Conn.Render(leds)

			time.Sleep(sleepTime)
		}
		copy(ledsOldLeft, ledsTargetLeft)
		copy(ledsOldRight, ledsTargetRight)
	}
}
