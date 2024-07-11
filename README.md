# ledstrips

WS281* library to control LEDs fully written in Go.
Tested on Raspberry Pi 3b and 4.

## Build for ARM

Set following path variables:

```powershell
$Env:GOOS="linux"
$Env:GOARCH="arm"
$Env:GOARM=7
```

## Troubleshooting

See [readme](https://github.com/jgarff) of jgarff's repository. Thanks at this point.

## SPI driver issue

Some issues occure within the SPI driver using DMA
 - https://github.com/jgarff/rpi_ws281x/issues/499
