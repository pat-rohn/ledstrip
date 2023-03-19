package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/pat-rohn/ledstrip"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var loglevel string

func main() {
	fmt.Println("Led Test Suite")
	var rootCmd = &cobra.Command{
		Use:   "worms",
		Short: "LED Strip Test Suite",
	}

	rootCmd.PersistentFlags().StringVarP(&loglevel, "verbose", "v", "w", "verbosity")
	var testCmd = &cobra.Command{
		Use:   "test ledNr Red(int) Green(int) Blue(int)",
		Args:  cobra.MinimumNArgs(4),
		Short: "Some LED test",
		Long: `
Examples: 
	worms test 20 50 60 70 -v i`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("test %+v\n", args)

			exitFunc := colorTest(args)

			waitForExitSignal(exitFunc)
			return nil
		},
	}
	var testRunCmd = &cobra.Command{
		Use:   "test-run ledNr test-version",
		Args:  cobra.MinimumNArgs(2),
		Short: "Some LED test-run",
		Long: `
Examples: 
	worms test-run 30 1 -v i`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("test %+v\n", args)
			exitFunc := colorTestRun(args)

			waitForExitSignal(exitFunc)
			return nil
		},
	}

	rootCmd.AddCommand(testCmd)
	rootCmd.AddCommand(testRunCmd)
	rootCmd.Execute()
}

func waitForExitSignal(exitFunction func()) {
	sigchnl := make(chan os.Signal, 1)
	signal.Notify(sigchnl)
	exitchnl := make(chan int)

	go func() {
		for {
			s := <-sigchnl
			if s == syscall.SIGTERM {
				fmt.Println("Got kill signal. ")
				exitFunction()
				fmt.Println("Program will terminate now.")
				os.Exit(0)
			} else if s == syscall.SIGINT {
				fmt.Println("Got CTRL+C signal")
				exitFunction()
				fmt.Println("Closing.")

				os.Exit(0)
			} else {
				fmt.Println("Ignoring signal: ", s)
			}
		}
	}()

	exitcode := <-exitchnl
	os.Exit(exitcode)
}

func colorTest(args []string) func() {
	fmt.Printf("args:%+v", args)
	ledNr, err := strconv.Atoi(args[0])
	if err != nil {
		log.Errorf("Invalid led number %s", args[0])
		ledNr = 2
	}
	var color ledstrip.RGBPixel
	c, err := strconv.Atoi(args[1])
	if err != nil {
		log.Error("Invalid red number %s", args[1])
	}
	color.Red = uint8(c)
	c, err = strconv.Atoi(args[2])
	if err != nil {
		log.Errorf("Invalid green number %s", args[2])
	}
	color.Green = uint8(c)
	c, err = strconv.Atoi(args[3])
	if err != nil {
		log.Errorf("Invalid blue number %s", args[3])
	}

	color.Blue = uint8(c)
	log.Infof("Color%+v", color)
	var leds []ledstrip.RGBPixel
	conn := ledstrip.NewSPI("/dev/spidev0.0", ledNr)

	for i := 0; i < ledNr; i++ {
		leds = append(leds, color)
	}

	conn.RenderLEDs(leds)

	return conn.Exit

}

func colorTestRun(args []string) func() {
	ledNr, err := strconv.Atoi(args[0])
	if err != nil {
		log.Error("Invalid led number %s", args[0])
		ledNr = 30
	}
	testVersion, err := strconv.Atoi(args[1])
	if err != nil {
		log.Error("Invalid led number %s", args[1])
		ledNr = 30
	}

	c := ledstrip.NewSPI("/dev/spidev0.0", ledNr)

	ledsTest := CreateTest()
	ledsWorms := CreateWorms()
	switch testVersion {
	case 0:
		go func() {
			RunLEDS(&c, ledsWorms, time.Second*2)
		}()
	case 1:
		go func() {
			RunLEDS(&c, ledsTest, time.Second*2)
		}()
	case 2:
		go func() {
			leds := []ledstrip.RGBPixel{{Red: 70, Blue: 35, Green: 0}}
			rDiff := 1
			gDiff := -2
			bDiff := -1
			max := uint8(70)
			for {
				time.Sleep(time.Millisecond * 100)
				oldColors := leds[0]
				if oldColors.Red > max {
					rDiff = -1
				}
				if oldColors.Red <= 4 {
					rDiff = 1
				}
				if oldColors.Green > max {
					gDiff = -2
				}
				if oldColors.Green <= 4 {
					gDiff = 2
				}
				if oldColors.Blue > max {
					bDiff = -1
				}
				if oldColors.Blue <= 4 {
					bDiff = 1
				}
				newColor := ledstrip.RGBPixel{
					Red:   uint8(int(oldColors.Red) + rDiff),
					Green: uint8(int(oldColors.Green) + gDiff),
					Blue:  uint8(int(oldColors.Blue) + bDiff)}
				fmt.Printf("newColor: %d %d %d\n", newColor.Red, newColor.Green, newColor.Blue)
				for i := range leds {
					leds[i] = newColor
					//leds[i] = ledstrip.RGBPixel{Red: oldValues.Red, Green: oldValues.Green, Blue: oldValues.Blue}
				}
				if len(leds) < ledNr {
					leds = append(leds, newColor)
				}
				c.RenderLEDs(leds)
			}
		}()

	}

	return c.Exit
}

// CreateWorms creats a slice of RGBPixels
func CreateWorms() []ledstrip.RGBPixel {
	logFields := log.Fields{"func": "CreateWorms"}
	log.WithFields(logFields).Traceln("CreateWorms")

	var leds []ledstrip.RGBPixel
	colorValues := [10]uint8{uint8(0), uint8(2), uint8(4), uint8(4), uint8(8), uint8(8), uint8(16), uint8(32), uint8(32), uint8(64)}

	for i := 0; i < 10; i++ {
		rVal := colorValues[ledstrip.Pmod(i, 10)]
		gVal := colorValues[0]
		bVal := colorValues[0]
		leds = append(leds, ledstrip.RGBPixel{
			Red:   rVal,
			Green: gVal,
			Blue:  bVal,
		})
	}
	for i := 0; i < 10; i++ {
		rVal := colorValues[0]
		gVal := colorValues[ledstrip.Pmod(i, 10)]
		bVal := colorValues[0]
		leds = append(leds, ledstrip.RGBPixel{
			Red:   rVal,
			Green: gVal,
			Blue:  bVal,
		})

	}
	for i := 0; i < 10; i++ {
		rVal := colorValues[0]
		gVal := colorValues[0]
		bVal := colorValues[ledstrip.Pmod(i, 10)]
		leds = append(leds, ledstrip.RGBPixel{
			Red:   rVal,
			Green: gVal,
			Blue:  bVal,
		})

	}
	return leds
}

// CreateTest creats a slice of RGBPixels
func CreateTest() []ledstrip.RGBPixel {
	logFields := log.Fields{"func": "Test1"}
	log.WithFields(logFields).Traceln("Test")

	var leds []ledstrip.RGBPixel
	colorValues := [10]uint8{uint8(0), uint8(4), uint8(16), uint8(32), uint8(64), uint8(32), uint8(16), uint8(8), uint8(4), uint8(0)}

	for i := 0; i < 10; i++ {
		rVal := colorValues[ledstrip.Pmod(i, 10)/2]
		gVal := colorValues[ledstrip.Pmod(i, 10)]
		bVal := colorValues[0]
		leds = append(leds, ledstrip.RGBPixel{
			Red:   rVal,
			Green: gVal,
			Blue:  bVal,
		})
	}
	for i := 0; i < 10; i++ {
		rVal := colorValues[0]
		gVal := colorValues[ledstrip.Pmod(i, 10)/2]
		bVal := colorValues[ledstrip.Pmod(i, 10)]
		leds = append(leds, ledstrip.RGBPixel{
			Red:   rVal,
			Green: gVal,
			Blue:  bVal,
		})
	}
	for i := 0; i < 10; i++ {
		rVal := colorValues[ledstrip.Pmod(i, 10)]
		gVal := colorValues[0]
		bVal := colorValues[ledstrip.Pmod(i, 10)/2]
		leds = append(leds, ledstrip.RGBPixel{
			Red:   rVal,
			Green: gVal,
			Blue:  bVal,
		})
	}
	return leds
}

func RunLEDS(c *ledstrip.ConnectionSPI, leds []ledstrip.RGBPixel, runTime time.Duration) {
	logFields := log.Fields{"func": "RunLEDS"}
	log.WithFields(logFields).Traceln("RunLEDS")
	runTime = runTime / 4
	endTime := time.Now().Add(runTime)
	for time.Now().Before(endTime) {
		c.RenderLEDs(leds)
		leds = ledstrip.PlaceInFront(leds, leds[len(leds)-1])
		waitTime := time.Until(endTime)
		time.Sleep(waitTime)

	}
	endTime = time.Now().Add(runTime)
	for time.Now().Before(endTime) {
		c.RenderLEDs(leds)
		leds = ledstrip.PlaceInFront(leds, leds[len(leds)-1])
		waitTime := time.Until(endTime)
		time.Sleep(waitTime)
	}

	leds = ledstrip.Inverse(leds)
	endTime = time.Now().Add(runTime)
	for time.Now().Before(endTime) {
		c.RenderLEDs(leds)
		leds = ledstrip.PlaceInBack(leds, leds[0])
		waitTime := time.Until(endTime)
		time.Sleep(waitTime)
	}
	endTime = time.Now().Add(runTime)
	for time.Now().Before(endTime) {
		c.RenderLEDs(leds)
		leds = ledstrip.PlaceInBack(leds, leds[0])
		waitTime := time.Until(endTime)
		time.Sleep(waitTime)
	}
}
