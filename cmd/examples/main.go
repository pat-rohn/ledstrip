package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pat-rohn/ledstrip"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"periph.io/x/conn/v3/spi"
)

var device string
var loglevel string
var fixSPI bool // https://github.com/jgarff/rpi_ws281x/issues/499

func main() {
	fmt.Println("LED Examples")
	var rootCmd = &cobra.Command{
		Use:   "examples",
		Short: "LED Strip Test Suite",
	}

	rootCmd.PersistentFlags().StringVarP(&loglevel, "verbose", "v", "w", "loggign verbosity")
	rootCmd.PersistentFlags().StringVarP(&device, "spi-device", "d", "/dev/spidev0.0", "SPI Device")
	rootCmd.PersistentFlags().BoolVarP(&fixSPI, "fix-spi", "f", false, "Fix SPI/DMA issue")

	var colorCmd = &cobra.Command{
		Use:   "color nrOfLeds Red Green Blue",
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
		Use:   "test [nrOfLeds] [test-version]",
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
		Use:   "clear nrOfLeds ",
		Args:  cobra.MinimumNArgs(1),
		Short: "Turn off LEDs",
		Long: `
Examples: 
	examples clear 30 -v i`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("test %+v\n", args)
			nrOfLeds, err := strconv.Atoi(args[0])
			if err != nil {
				log.Errorf("Invalid led number %s", args[0])
				nrOfLeds = 30
			}
			c, err := ledstrip.NewSPI(device, nrOfLeds, fixSPI)
			if err != nil {
				return err
			}
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
	if fixSPI {
		log.Warn("Fixing SPI issue is active")
	}
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
		log.Errorf("Invalid number of leds: %s", args[0])
		nrOfLeds = 30
	}
	testVersion, err := strconv.Atoi(args[1])
	if err != nil {
		log.Errorf("Invalid test %s", args[1])
		nrOfLeds = 30
	}
	c, err := ledstrip.New(device, nrOfLeds, spi.Mode2, fixSPI)
	if err != nil {
		log.Fatalf("Failed to create SPI connection: %v", err)
	}
	switch testVersion {
	case 0:
		runExample0(&c, nrOfLeds)
	case 1:
		runExample1(&c, nrOfLeds)
	case 2:
		runExample2(&c, nrOfLeds)
	case 3:
		maskLength := 20
		if len(args) > 2 {
			maskLength, err = strconv.Atoi(args[2])
			if err != nil {
				log.Errorf("Invalid mask length: %s", args[2])
				maskLength = 20
			}
		}

		color1 := ledstrip.RGBPixel{Red: 0, Green: 100, Blue: 0}
		color2 := ledstrip.RGBPixel{Red: 0, Green: 0, Blue: 100}
		if len(args) > 4 {
			colorVals1 := strings.Split(args[3], ",")
			colorVals2 := strings.Split(args[4], ",")
			if len(colorVals1) == 3 && len(colorVals2) == 3 {
				v, err := strconv.Atoi(colorVals1[0])
				if err != nil {
					log.Errorf("Invalid color value: %s", colorVals1[0])
				}
				color1.Red = uint8(v)
				v, err = strconv.Atoi(colorVals1[1])
				if err != nil {
					log.Errorf("Invalid color value: %s", colorVals1[1])

				}
				color1.Green = uint8(v)
				v, err = strconv.Atoi(colorVals1[2])
				if err != nil {
					log.Errorf("Invalid color value: %s", colorVals1[2])

				}
				color1.Blue = uint8(v)

				v, err = strconv.Atoi(colorVals2[0])
				if err != nil {
					log.Errorf("Invalid color value: %s", colorVals2[0])

				}
				color2.Red = uint8(v)
				v, err = strconv.Atoi(colorVals2[1])
				if err != nil {
					log.Errorf("Invalid color value: %s", colorVals2[1])

				}
				color2.Green = uint8(v)
				v, err = strconv.Atoi(colorVals2[2])
				if err != nil {
					log.Errorf("Invalid color value: %s", colorVals2[2])

				}
				color2.Blue = uint8(v)
			}
		}

		runExample3(&c, nrOfLeds, maskLength, color1, color2)
	case 4:
		maskLength := 20
		if len(args) > 2 {
			maskLength, err = strconv.Atoi(args[2])
			if err != nil {
				log.Errorf("Invalid mask length: %s", args[2])
				maskLength = 20
			}
		}

		color1 := ledstrip.RGBPixel{Red: 0, Green: 100, Blue: 0}
		color2 := ledstrip.RGBPixel{Red: 0, Green: 0, Blue: 100}
		if len(args) > 4 {
			colorVals1 := strings.Split(args[3], ",")
			colorVals2 := strings.Split(args[4], ",")
			if len(colorVals1) == 3 && len(colorVals2) == 3 {
				v, err := strconv.Atoi(colorVals1[0])
				if err != nil {
					log.Errorf("Invalid color value: %s", colorVals1[0])
				}
				color1.Red = uint8(v)
				v, err = strconv.Atoi(colorVals1[1])
				if err != nil {
					log.Errorf("Invalid color value: %s", colorVals1[1])

				}
				color1.Green = uint8(v)
				v, err = strconv.Atoi(colorVals1[2])
				if err != nil {
					log.Errorf("Invalid color value: %s", colorVals1[2])

				}
				color1.Blue = uint8(v)

				v, err = strconv.Atoi(colorVals2[0])
				if err != nil {
					log.Errorf("Invalid color value: %s", colorVals2[0])

				}
				color2.Red = uint8(v)
				v, err = strconv.Atoi(colorVals2[1])
				if err != nil {
					log.Errorf("Invalid color value: %s", colorVals2[1])

				}
				color2.Green = uint8(v)
				v, err = strconv.Atoi(colorVals2[2])
				if err != nil {
					log.Errorf("Invalid color value: %s", colorVals2[2])

				}
				color2.Blue = uint8(v)
			}
		}

		runExample4(&c, nrOfLeds, maskLength, color1, color2)
	default:
		log.Fatal("Unknown test")

	}
}

func showColor(args []string) {
	fmt.Printf("args: %+v\n", args)
	ledNr, err := strconv.Atoi(args[0])
	if err != nil {
		log.Errorf("Invalid number leds: %s", args[0])
		ledNr = 2
	}
	var color ledstrip.RGBPixel
	c, err := strconv.Atoi(args[1])
	if err != nil {
		log.Errorf("Invalid red number %s", args[1])
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
	conn, err := ledstrip.NewSPI(device, ledNr, fixSPI)
	if err != nil {
		log.Errorf("Failed to create SPI connection: %v", err)
		return
	}
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
