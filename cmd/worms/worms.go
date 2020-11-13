package main

import (
	"time"

	"github.com/pat-rohn/ledstrip"
)

func main() {
	conn := ledstrip.NewSPI("/dev/spidev0.0")
	leds := ledstrip.CreateWorms()
	runTime := time.Second * 30
	conn.RunLEDS(leds, runTime)
	conn.Clear(30)
	conn.Close()
}
