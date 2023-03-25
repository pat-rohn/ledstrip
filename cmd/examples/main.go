package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/pat-rohn/ledstrip"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var device string
var loglevel string

func main() {
	fmt.Println("Led Test Suite")
	var rootCmd = &cobra.Command{
		Use:   "examples",
		Short: "LED Strip Test Suite",
	}

	rootCmd.PersistentFlags().StringVarP(&loglevel, "verbose", "v", "w", "loggign verbosity")
	rootCmd.PersistentFlags().StringVarP(&device, "spi-evice", "d", "/dev/spidev0.0", "SPI Device")

	var colorCmd = &cobra.Command{
		Use:   "color ledNr Red(int) Green(int) Blue(int)",
		Args:  cobra.MinimumNArgs(4),
		Short: "Some LED test",
		Long: `
Examples: 
	examples color 20 50 60 70 -v i`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("test %+v\n", args)

			showColor(args)
			return nil
		},
	}
	var testRunCmd = &cobra.Command{
		Use:   "test ledNr test-version",
		Args:  cobra.MinimumNArgs(2),
		Short: "Some LED test-run",
		Long: `
Examples: 
	examples test 30 1 -v i`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("test %+v\n", args)
			RunTest(args)

			return nil
		},
	}
	var clearCmd = &cobra.Command{
		Use:   "clear ledNr ",
		Args:  cobra.MinimumNArgs(1),
		Short: "Turn off LEDs",
		Long: `
Examples: 
	examples clear 30 -v i`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("test %+v\n", args)
			nrOfLeds, err := strconv.Atoi(args[0])
			if err != nil {
				log.Error("Invalid led number %s", args[0])
				nrOfLeds = 30
			}
			c := ledstrip.NewSPI(device, nrOfLeds)
			c.Clear()

			return nil
		},
	}
	rootCmd.AddCommand(colorCmd)
	rootCmd.AddCommand(testRunCmd)
	rootCmd.AddCommand(clearCmd)
	cobra.OnInitialize(initGlobalFlags)
	rootCmd.Execute()
}

func initGlobalFlags() {
	setLogLevel(loglevel)
}

func setLogLevel(level string) {
	switch level {
	case "-trace", "t":
		log.SetLevel(log.TraceLevel)
	case "-info", "i":
		log.SetLevel(log.InfoLevel)
	case "-warn", "w":
		log.SetLevel(log.WarnLevel)
	case "-error", "e":
		log.SetLevel(log.ErrorLevel)
	default:
		fmt.Printf("Invalid log-level %s\n", level)
		return
	}
	fmt.Printf("LogLevel is set to %s\n", level)
}

func RunTest(args []string) {

	nrOfLeds, err := strconv.Atoi(args[0])
	if err != nil {
		log.Error("Invalid led number %s", args[0])
		nrOfLeds = 30
	}
	testVersion, err := strconv.Atoi(args[1])
	if err != nil {
		log.Error("Invalid test %s", args[1])
		nrOfLeds = 30
	}

	switch testVersion {
	case 0:
		runExample0(nrOfLeds)
	case 1:
		runExample1(nrOfLeds)
	case 2:
		runExample2(nrOfLeds)
	default:
		log.Fatal("Unknown test")

	}
}

type Runner struct {
	c *ledstrip.ConnectionSPI
}

func (r *Runner) RunLEDS(leds []ledstrip.RGBPixel, runTime time.Duration) {
	logFields := log.Fields{"func": "RunLEDS"}
	log.WithFields(logFields).Traceln("RunLEDS")
	runTime = runTime / 4
	endTime := time.Now().Add(runTime)
	for time.Now().Before(endTime) {
		r.c.Render(leds)
		leds = ledstrip.PlaceInFront(leds, leds[len(leds)-1])
		time.Sleep(time.Millisecond * 50)
	}

	endTime = time.Now().Add(runTime)
	for time.Now().Before(endTime) {
		r.c.Render(leds)
		leds = ledstrip.PlaceInFront(leds, leds[len(leds)-1])
		time.Sleep(time.Millisecond * 50)
	}

	leds = ledstrip.Inverse(leds)
	endTime = time.Now().Add(runTime)
	for time.Now().Before(endTime) {
		r.c.Render(leds)
		leds = ledstrip.PlaceAtBack(leds, leds[0])
		time.Sleep(time.Millisecond * 50)
	}
	endTime = time.Now().Add(runTime)
	for time.Now().Before(endTime) {
		r.c.Render(leds)
		leds = ledstrip.PlaceAtBack(leds, leds[0])
		time.Sleep(time.Millisecond * 50)
	}
}

func runExample0(nrOfLeds int) {

	c := ledstrip.NewSPI(device, nrOfLeds)
	leds := []ledstrip.RGBPixel{{Red: 30, Blue: 6, Green: 10}}
	rDiff := -3
	gDiff := -2
	bDiff := 1
	max := uint8(70)
	min := uint8(5)
	for {
		time.Sleep(time.Millisecond * 20)
		oldColors := leds[0]
		if oldColors.Red > max {
			rDiff *= -1
		}
		if oldColors.Red <= min {
			rDiff *= -1
		}
		if oldColors.Green > max {
			gDiff *= -1
		}
		if oldColors.Green <= min {
			gDiff *= -1
		}
		if oldColors.Blue > max {
			bDiff *= -1
		}
		if oldColors.Blue <= min {
			bDiff *= -1
		}
		newColor := ledstrip.RGBPixel{
			Red:   uint8(int(oldColors.Red) + rDiff),
			Green: uint8(int(oldColors.Green) + gDiff),
			Blue:  uint8(int(oldColors.Blue) + bDiff),
		}
		fmt.Printf("newColor: %d %d %d\n", newColor.Red, newColor.Green, newColor.Blue)
		for i := range leds {
			leds[i] = newColor
		}
		if len(leds) < nrOfLeds {
			leds = append(leds, newColor)
		}
		c.Render(leds)
	}

}

func runExample1(nrOfLeds int) {

	fmt.Println("example1")
	ledsWorms := CreateExample1(nrOfLeds)
	c := ledstrip.NewSPI(device, nrOfLeds)
	runner := Runner{
		c: &c,
	}
	for {
		runner.RunLEDS(ledsWorms, time.Second*10)
	}
}

func runExample2(nrOfLeds int) {

	fmt.Println("test2")
	example := CreateExample2(nrOfLeds)
	c := ledstrip.NewSPI(device, nrOfLeds)
	runner := Runner{
		c: &c,
	}

	for {
		// do not run in go routine
		runner.RunLEDS(example, time.Second*10)
	}
}

func showColor(args []string) {
	fmt.Printf("args: %+v\n", args)
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
	log.Infof("Color %d", ledNr)
	var leds []ledstrip.RGBPixel
	/*for i := 0; i < 25; i++ {
		leds = append(leds, color)
	}*/
	conn := ledstrip.NewSPI(device, ledNr)
	fmt.Print("No ")
	for i := len(leds); i < ledNr; i++ {
		if len(leds) <= ledNr {
			leds = append(leds, color)
			time.Sleep(time.Millisecond * 100)
			fmt.Printf(" %d .. ", i)
		}
		conn.Render(leds)
	}
	for {
		for i := 0; i < len(leds); i++ {
			leds = ledstrip.PlaceInFront(leds, ledstrip.RGBPixel{})
			time.Sleep(time.Millisecond * 20)
			conn.Render(leds)
		}
		for i := 0; i < len(leds); i++ {
			leds = ledstrip.PlaceInFront(leds, color)
			time.Sleep(time.Millisecond * 20)
			conn.Render(leds)
		}
	}
}

func CreateExample1(nrOfLeds int) []ledstrip.RGBPixel {
	logFields := log.Fields{"func": "CreateExample1"}

	var leds []ledstrip.RGBPixel
	colorValues := [10]uint8{uint8(0), uint8(2), uint8(4), uint8(4), uint8(8), uint8(8), uint8(16), uint8(32), uint8(32), uint8(64)}
	for len(leds) < nrOfLeds {
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
	}
	leds = leds[:nrOfLeds]
	log.WithFields(logFields).Infof("Example with %d LEDs created", len(leds))
	return leds
}

func CreateExample2(nrOfLeds int) []ledstrip.RGBPixel {
	logFields := log.Fields{"func": "CreateExample2"}
	log.WithFields(logFields).Traceln("Example 2")

	var leds []ledstrip.RGBPixel
	colorValues := [10]uint8{uint8(0), uint8(4), uint8(16), uint8(32), uint8(64), uint8(32), uint8(16), uint8(8), uint8(4), uint8(0)}
	for len(leds) < nrOfLeds {
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
	}
	leds = leds[:nrOfLeds]
	log.WithFields(logFields).Infof("Example with %d LEDs created", len(leds))
	return leds
}
