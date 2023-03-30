# ledstrips

WS281* library to control leds fully written in Go. 
Tested on Raspberry Pi 3b and 4.

## Crosscompile on Linux for ARM

Export following variables:

```terminal
export GOARCH=arm
export GOARM=7
```

## Crosscompile on Windows for ARM

Set following path variables:

```cmd
set GOOS=linux
set GOARCH=arm
set GOARM=7
```

```powershell
$Env:GOOS="linux"
$Env:GOARCH="arm"
$Env:GOARM=7
```

# Troubleshooting
See [readme](https://github.com/jgarff) of jgarff's repository. Thanks at this point.


## SPI driver issue

Some issues occure within the SPI driver using DMA
 - https://github.com/jgarff/rpi_ws281x/issues/499

Workaround ignoring lowest bits:

```diff
diff --git a/colordata.go b/colordata.go
index 6ed29dd..b379a93 100644
--- a/colordata.go
+++ b/colordata.go
@@ -87,7 +87,7 @@ func (c *ColorData) addBits(high bool) {
 
 func (c *ColorData) addColorValue(color uint8) {
        for i := 7; i >= 0; i-- { // Most Significant bit first
-               if HasBit(color, uint8(i)) {
+               if HasBit(color, uint8(i)) && i > 2 { // Workaround for SPI DMA issue
                        c.addBits(true)
                } else {
                        c.addBits(false)
```