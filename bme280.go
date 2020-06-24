package main

import (
	"github.com/maciej/bme280"
	"golang.org/x/exp/io/i2c"
)

func getBME280() *bme280.Driver {
	device, err := i2c.Open(&i2c.Devfs{Dev: "/dev/i2c-1"}, 0x76)
	if err != nil {
		log.Fatal(err)
	}

	driver := bme280.New(device)
	err = driver.InitWith(bme280.ModeForced, bme280.Settings{
		Filter:                  bme280.FilterOff,
		Standby:                 bme280.StandByTime1000ms,
		PressureOversampling:    bme280.Oversampling16x,
		TemperatureOversampling: bme280.Oversampling16x,
		HumidityOversampling:    bme280.Oversampling16x,
	})

	if err != nil {
		log.Fatal(err)
	}

	return driver
}
