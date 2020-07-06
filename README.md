# ledstrips
WS281* library to control leds fully written in Go. 
Tested on Raspberry Pi using interface "/dev/spidev0.0".


## Example
```go
    conn := ledstrip.NewSPI()
	leds := ledstrip.CreateWorms()
	runTime := time.Second * 30
	conn.RunLEDS(leds, runTime)
	conn.Clear(30)
	conn.Close()

```

## Crosscompile for ARM
Run following commands in terminal to build for ARM
```terminal
  export GOARCH=arm
  export GOARM=7
  export GOBIN=
```

