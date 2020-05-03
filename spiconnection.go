package ledstrip

import (
	"time"

	log "github.com/sirupsen/logrus"
	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/conn/spi"
	"periph.io/x/periph/conn/spi/spireg"
	"periph.io/x/periph/host"
)

//https://cdn-reichelt.de/documents/datenblatt/A300/ADAFRUIT_1643_ENG_TDS.pdf
//https://cdn-shop.adafruit.com/datasheets/WS2812.pdf

// Connection contains connection and closer
type Connection struct {
	portCloser spi.PortCloser
	spiDev     spi.Conn
}

//New creates a SPI-Connection
func New() Connection {

	devicePath := "/dev/spidev0.0"
	hz := physic.Hertz * 2400000
	log.WithField("package", logPkg).Tracef("Open spi-dev '%s' with max speed '%v'", devicePath, hz)
	spiMode := spi.Mode(spi.Mode0)
	log.WithField("package", logPkg).Infof("SPI info:  %v \n", spiMode.String())

	host.Init()

	s, err := spireg.Open(devicePath)
	if err != nil {
		log.WithField("package", logPkg).Fatalf(
			"Failed to open Pt100Connection connection:  %v \n", err)
	}
	c, err := s.Connect(hz, spiMode, 8)
	if err != nil {
		log.WithField("package", logPkg).Fatalf(
			"Failed to connect %v\n", err)
	}

	conn := Connection{portCloser: s, spiDev: c}
	log.WithField("package", logPkg).Tracef("SPI info:  %v \n", c.String())
	log.WithField("package", logPkg).Tracef("SPI info:  %v \n", c.Duplex().String())
	return conn
}

//Close closes SPIConnection
func (conn *Connection) Close() error {
	logFields := log.Fields{"package": logPkg, "func": "Close"}
	log.WithField("package", logPkg).Tracef("Close")
	err := conn.portCloser.Close()
	if err != nil {
		log.WithFields(logFields).Errorf(
			"Failed to close SPIConnection connection  %v on \n", err)
		return err
	}
	return nil
}

//RenderLEDs translates RGBPixels into SPI message and transfers the message
func (conn *Connection) RenderLEDs(pixels []RGBPixel) {

	var translatedRGBs []uint8
	for _, pixel := range pixels {
		// Composition of 24bit data of a pixel, is ordered GRB
		translatedRGBs = append(translatedRGBs, getTranslatedColor([3]uint8{pixel.Green, pixel.Red, pixel.Blue})...)
	}
	conn.transfer(translatedRGBs)
}

/*
 * Copyright (c) 2014 Jeremy Garff <jer @ jers.net>
 * https://github.com/jgarff/rpi_ws281x/blob/master/ws2811.c
 * https://github.com/jgarff/rpi_ws281x/blob/master/LICENSE
 */
func getTranslatedColor(pixel [3]uint8) []uint8 {
	logFields := log.Fields{"package": logPkg, "func": "getTranslatedColor"}
	log.WithFields(logFields).Traceln("getTranslatedColor")

	var rgbTranslated []uint8
	low := uint8(0x4)  // 1 0 0
	high := uint8(0x6) // 1 1 0
	//lowInv := 0x1  // 0 0 1
	//highInv := 0x3 // 0 1 1
	bitpos := 7

	translatedByte := uint8(0)
	for colorNr := 0; colorNr < 3; colorNr++ {

		log.WithFields(logFields).Tracef("COLOR %v\n", pixel[colorNr])

		for bitNr := 7; bitNr >= 0; bitNr-- {
			symbol := uint8(low)
			/*if (channel->invert)
			{
				symbol = SYMBOL_LOW_INV;
			}    */
			log.WithFields(logFields).Tracef("translatedByte 1 %v\n", translatedByte)
			log.WithFields(logFields).Tracef("color[j] & (1 << k) 1 %v\n", pixel[colorNr]&(1<<bitNr))

			if pixel[colorNr]&(1<<bitNr) > 0 {
				symbol = high
				log.WithFields(logFields).Tracef("Set symbol to high %v\n", symbol)

				/*if (channel->invert)
				{
					symbol = SYMBOL_HIGH_INV;
				} */
			}
			log.WithFields(logFields).Tracef("bit %v\n", bitNr)
			for l := 2; l >= 0; l-- {
				if (symbol & (1 << l)) > 0 {
					translatedByte |= (1 << bitpos)
				}

				log.WithFields(logFields).Tracef("%v translatedByte %v\n", colorNr, translatedByte)

				bitpos--
				if bitpos < 0 {
					rgbTranslated = append(rgbTranslated, translatedByte)
					translatedByte = 0
					log.WithFields(logFields).Tracef("%v translatedByte  %v\n", colorNr, translatedByte)

					bitpos = 7
				}
			}
		}
		log.WithFields(logFields).Tracef("%v --- My Byte %v", colorNr, translatedByte)
	}
	return rgbTranslated
}

func (conn *Connection) transfer(msg []byte) []byte {

	res := make([]byte, len(msg))
	for i := range res {
		res[i] = 0xff
	}
	err := conn.spiDev.Tx(msg, res)
	if err != nil {
		log.WithField("package", logPkg).Errorf("Tx failed  %v  \n", err)
		return res
	}
	time.Sleep(time.Microsecond * 50)
	return res
}
