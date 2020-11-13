# ledstrips

WS281* library to control leds fully written in Go. 
Tested on Raspberry Pi 3b and 4.


## Raspberry Pi Settings
[Jeremy Garff](https://github.com/jgarff/rpi_ws281x "github.com") has a good guide how to configure the SPI.
For me it worked already by activating SPI (sudo raspi-config) using Raspberry Pi OS (stretch / buster).

## Crosscompile on Linux for ARM

Export following variables:

```terminal
  export GOARCH=arm
  export GOARM=7
  export GOBIN=
```

## Crosscompile on Windows for ARM

Set following path variables:

```cmd
set GOOS=linux
set GOARCH=arm
set GOARM=7
``` 

## Example

```go
    conn := ledstrip.NewSPI("/dev/spidev0.0")
	leds := ledstrip.CreateWorms()
	runTime := time.Second * 30
	conn.RunLEDS(leds, runTime)
	conn.Clear(30)
	conn.Close()
```

## To Do
- Support for PWM and PCM

# Special Thanks
Thank you to [Jeremy Garff](https://github.com/jgarff) for writing the C library rpi_ws281x.

