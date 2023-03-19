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

