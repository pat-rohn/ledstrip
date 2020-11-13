package ledstrip

import (
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
)

func TestSPI(t *testing.T) {
	log.SetLevel(log.TraceLevel)

	log.Traceln("TestSPI")
	ledsTest := CreateTest()
	ledsWorms := CreateWorms()
	endTime := time.Now().Add(10 * time.Minute)
	conn := NewSPI("/dev/spidev0.0")
	waitTime := 15 * time.Second
	for time.Now().Before(endTime) {
		go conn.RunLEDS(ledsTest, waitTime)
		time.Sleep(waitTime)
		go conn.RunLEDS(ledsWorms, waitTime)
		time.Sleep(waitTime)
	}

	//go
	conn.Clear(30)
	defer conn.Close()

}
