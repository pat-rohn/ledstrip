package ledstrip

import (
	"time"

	log "github.com/sirupsen/logrus"
	"periph.io/x/conn/v3/driver/driverreg"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/conn/v3/spi"
	"periph.io/x/conn/v3/spi/spireg"
	host "periph.io/x/host/v3"
)

//https://cdn-reichelt.de/documents/datenblatt/A300/ADAFRUIT_1643_ENG_TDS.pdf
//https://cdn-shop.adafruit.com/datasheets/WS2812.pdf

// ConnectionSPIs
type ConnectionSPI struct {
	portCloser spi.PortCloser
	spiDev     spi.Conn
	NrOfLeds   int
	FixSPI     bool
}

// NewSPI creates a SPI-Connection
func NewSPI(devicePath string, nrOfLEDs int, fixSPI bool) ConnectionSPI {
	logFields := log.Fields{"package": logPkg, "conn": "SPI", "func": "NewSPI"}
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}
	hz := physic.Hertz * 2400000
	log.WithFields(logFields).Infof("Open spi-dev '%s' with max speed '%v'", devicePath, hz)
	spiMode := spi.Mode(spi.Mode0)
	log.WithFields(logFields).Infof("SPI mode:  %v %b\n", spiMode.String(), spiMode)
	if _, err := driverreg.Init(); err != nil {
		log.Fatal(err)
	}

	s, err := spireg.Open(devicePath)
	if err != nil {
		log.WithFields(logFields).Fatalf("Failed to open SPI connection:  %v \n", err)
	}
	c, err := s.Connect(hz, spiMode, 8)
	if err != nil {
		log.WithFields(logFields).Fatalf("Failed to connect %v\n", err)
	}

	conn := ConnectionSPI{
		portCloser: s,
		spiDev:     c,
		NrOfLeds:   nrOfLEDs,
		FixSPI:     fixSPI,
	}

	return conn
}

// Render translates RGBPixels into SPI message and transfers the message
func (c *ConnectionSPI) Render(pixels []RGBPixel) {
	logFields := log.Fields{"package": logPkg, "conn": "SPI", "func": "RenderLEDs"}
	log.WithFields(logFields).Infof("Render %d LEDs", len(pixels))
	// Fix for Raspberry Pi 3 Model B+ (5.15.84-v7+)
	// Data signal seems to be splitted sending less than 11 LEDS
	if len(pixels) < 11 {
		pixels = append(pixels, []RGBPixel{{}, {}, {}, {}, {}, {}, {}, {}, {}, {}}...)
	}
	var translatedRGBs []uint8
	for _, pixel := range pixels {

		colorData := GetColorData(pixel, c.FixSPI)
		log.WithFields(logFields).Tracef("%08b", pixel)

		for _, c := range colorData {
			translatedRGBs = append(translatedRGBs, c)
		}
	}

	c.transfer(translatedRGBs)
}

// Close closes SPIConnection
func (c *ConnectionSPI) Close() error {
	logFields := log.Fields{"package": logPkg, "conn": "SPI", "func": "Close"}
	log.WithFields(logFields).Tracef("Close")
	err := c.portCloser.Close()
	if err != nil {
		log.WithFields(logFields).Errorf(
			"Failed to close SPIConnection connection: %v\n", err)
		return err
	}
	return nil
}

// Exit turns off LEDs and closes SPIConnection
func (c *ConnectionSPI) Exit() {
	logFields := log.Fields{"package": logPkg, "conn": "SPI", "func": "Exit"}
	log.WithFields(logFields).Warn("Exit")

	c.Clear()
	c.Close()
}

func (c *ConnectionSPI) Clear() {
	logFields := log.Fields{"package": logPkg, "func": "Clear"}
	log.WithFields(logFields).Infof("Clear")

	var clearArray []uint8
	for i := 0; i < c.NrOfLeds; i++ {
		for _, color := range GetColorData(RGBPixel{Red: 0, Green: 0, Blue: 0}, c.FixSPI) {
			clearArray = append(clearArray, color)
		}
	}

	c.transfer(clearArray)
}

func (c *ConnectionSPI) transfer(msg []byte) []byte {

	log.Infof("transfer %d bytes", len(msg))

	res := make([]byte, len(msg))
	for i := range res {
		res[i] = 0xff
	}
	err := c.spiDev.Tx(msg, res)
	if err != nil {
		log.WithField("package", logPkg).Errorf("Tx failed  %v  \n", err)
		return res
	}
	time.Sleep(time.Microsecond * 300) // Above 50μs to reset (some need more than 300us)
	return res
}
