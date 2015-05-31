package main

import (
	"github.com/kidoman/embd"
	"github.com/kidoman/embd/controller/hd44780"
	_ "github.com/kidoman/embd/host/rpi"
	"github.com/kidoman/embd/interface/display/characterdisplay"
	"os"
	"os/signal"
	"pi-clock/weather"
	"syscall"
)

func main() {
	quitChannel := make(chan os.Signal)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)

	go DisplayServices()

	<-quitChannel
}

func NewDisplay() *characterdisplay.Display {
	if err := embd.InitI2C(); err != nil {
		panic(err)
	}

	bus := embd.NewI2CBus(1)

	controller, err := hd44780.NewI2C(
		bus,
		0x27,
		hd44780.PCF8574PinMap,
		hd44780.RowAddress20Col,
		hd44780.TwoLine,
		hd44780.BlinkOn,
	)
	if err != nil {
		panic(err)
	}

	display := characterdisplay.New(controller, 20, 4)
	display.Clear()

	return display
}

func DisplayServices() {
	display := NewDisplay()
	defer embd.CloseI2C()
	defer display.Close()
	display.BacklightOn()

	weatherService := weather.NewWeatherService()
	forcast := weatherService.GetCurrentForcast()

	display.Message(forcast)
}
