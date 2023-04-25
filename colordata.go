package ledstrip

import (
	log "github.com/sirupsen/logrus"
)

type ColorData struct {
	bitCounter  uint8
	byteCounter int
	ColorData   [9]uint8
	FixSPI      bool
	ignoreNext  bool
}

/*
one color-value has 3 bits

high = 1 1 0
low = 1 0 0

0 24bit Green
1 24bit Red
2 24bit Blue
=> 9 byte

Set all to 0

	  G7      G6      G5      G4      G3      G2      G1      G0      /      R7      R6      R5   ...
	[1 0 0] [1 0 0] [1 0 0] [1 0 0] [1 0 0] [1 0 0] [1 0 0] [1 0 0]       [1 0 0] [1 0 0] [1 0 0] ...
	   0       0       0       0       0       0       0       0             0       0       0    ...

Set all to 1

		G7      G6      G5      G4      G3      G2      G1      G0      /      R7      R6      R5   ...
	[1 1 0] [1 1 0] [1 1 0] [1 1 0] [1 1 0] [1 1 0] [1 1 0] [1 1 0]       [1 1 0] [1 1 0] [1 1 0] ...
		1       1       1       1       1       1       1       1             1       1       1    ...
	   128     64       32      16      8       4       2       1             128     64      32   ...

SPI issue

		G7      G6     G5          G4      G3          G2      G1      G0        /      R7      R6      R5   ...
	[1 1 0] [1 1 0] [1 1 XX 0 ]  [1 1 0] [1 1 0] [1  XX 1 0] [1 1 0] [1 1 0] XX      [1 1 0] [1 1 0] [1 1 0] ...
	   1       1       1 XX         1       1        XX 1       1       1    XX         1       1       1    ...
	  128      64     32 XX         16      8        XX 4       2       1    XX         128     64      32   ...

Set Red

	  G7      G6      G5      G4      G3      G2      G1      G0      /      R7      R6      R5   ...
	[1 0 0] [1 0 0] [1 0 0] [1 0 0] [1 0 0] [1 0 0] [1 0 0] [1 0 0]       [1 0 0] [1 1 0] [1 0 0] ...
	   0       0       0       0       0       0       0       0             0       1       0    ...


	100 100 100 100 100 100 100 100       100 110 100 100 100 100 100 100 100  100 100 100 100 100 100 100
	100 100 100 100 100 100 100 100       100 100 100 100 100 100 100 100 100  100 100 100 100 100 100 100
	100 100 100 100 100 100 100 110       100 100 100 100 100 100 100 100 100  100 100 100 100 100 100 100

	10010010 01001001 00100100  11011011 01101101 10110110  10010010 01001001 00100100
	10010010 01001001 00100100  11011011 01101101 10110110  10010010 01001001 00100100
*/
func GetColorData(pixel RGBPixel, fixSPI bool) [9]uint8 {
	colorData := ColorData{
		bitCounter:  0,
		byteCounter: 0,
		ColorData:   [9]uint8{},
		FixSPI:      fixSPI,
		ignoreNext:  false,
	}
	colorData.addColorValue(pixel.Green)
	colorData.addColorValue(pixel.Red)
	colorData.addColorValue(pixel.Blue)

	log.Tracef("%08b", colorData.ColorData)
	return colorData.ColorData
}

func (c *ColorData) setNextBit(high bool) {
	if high {
		c.ColorData[c.byteCounter] = SetBit(c.ColorData[c.byteCounter], c.bitCounter)
	}
	if c.ignoreNext && c.FixSPI {
		c.ignoreNext = false
		return
	}

	c.bitCounter++
	if c.bitCounter > 7 {
		c.ignoreNext = true
		log.Tracef("%d - %08b", c.byteCounter, c.ColorData[c.byteCounter])
		c.bitCounter = 0
		c.byteCounter++
		if c.byteCounter > 8 {
			c.byteCounter = 0
		}
	}
}

func (c *ColorData) addBits(high bool) {
	// color-value as 3 bits for high or low (1 1 0 / 1 0 0)
	c.setNextBit(true)
	if high {
		c.setNextBit(true)
	} else {
		c.setNextBit(false)
	}
	c.setNextBit(false)
}

func (c *ColorData) addColorValue(color uint8) {
	for i := 7; i >= 0; i-- { // Most Significant bit first
		if HasBit(color, uint8(i)) {
			c.addBits(true)
		} else {
			c.addBits(false)
		}
	}
}

func HasBit(n uint8, pos uint8) bool {
	val := n & uint8(1<<pos)
	return (val > 0)
}

func SetBit(n uint8, pos uint8) uint8 {
	pos = 7 - pos
	n |= (1 << pos)
	return n
}
