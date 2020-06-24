package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go"
	lcd "github.com/wjessop/lcm1602_lcd"
	"golang.org/x/exp/io/i2c"
)

const (
	databaseName = "climate"
)

var (
	// Flag, wether to update an LCD display
	updateLCD bool

	// Flag, name of the location to tag data with
	locationName string

	// Address of the influx DB host to send data to
	influxhost string

	lcdDisplay *lcd.LCM1602LCD
)

func init() {
	flag.BoolVar(&updateLCD, "lcd", false, "wether to output temp/humidity data to an attached LCD")
	flag.StringVar(&locationName, "name", "", "the location name to tag metrics with")
	flag.StringVar(&influxhost, "influxhost", "127.0.0.1", "the host address of the influx DB")
	flag.Parse()
}

func main() {
	debug := false
	if os.Getenv("DEBUG") != "" {
		debug = true
	}

	createLogger(debug)

	if locationName == "" {
		log.Fatal("Please provide a non-blank location name")
	}

	log.Debugf("Update LCD display set to %v", updateLCD)

	// create new client with default option for server url authenticate by token
	client := influxdb2.NewClient(
		fmt.Sprintf("http://%s:8086", influxhost),
		fmt.Sprintf("climate-writer:%s", os.Getenv("INFLUX_DB_SECRET")),
	)

	// user blocking write client for writes to desired bucket
	writeAPI := client.WriteApiBlocking("", databaseName)

	bme280 := getBME280()
	defer bme280.Close()

	var ctx = context.Background()

	if updateLCD {
		lcdDevice, err := i2c.Open(&i2c.Devfs{Dev: "/dev/i2c-1"}, 0x27)
		if err != nil {
			log.Fatal(err)
		}
		defer lcdDevice.Close()

		lcdDisplay, err = lcd.NewLCM1602LCD(lcdDevice)
		if err != nil {
			log.Fatal(err)
		}
	}

	for {
		response, err := bme280.Read()
		if err != nil {
			log.Fatal(err)
		}

		log.Debugf("Temp: %.1fC, Press: %.1fhPa, Hum: %.1f%%\n", response.Temperature, response.Pressure, response.Humidity)

		now := time.Now()

		p := influxdb2.NewPoint(
			"stat",
			map[string]string{"unit": "temperature", "location": locationName},
			map[string]interface{}{"value": response.Temperature},
			now,
		)

		if err = writeAPI.WritePoint(ctx, p); err != nil {
			log.Fatal(err)
		}

		p = influxdb2.NewPoint(
			"stat",
			map[string]string{"unit": "pressure", "location": locationName},
			map[string]interface{}{"value": response.Pressure},
			now,
		)

		if err = writeAPI.WritePoint(ctx, p); err != nil {
			log.Fatal(err)
		}

		p = influxdb2.NewPoint(
			"stat",
			map[string]string{"unit": "humidity", "location": locationName},
			map[string]interface{}{"value": response.Humidity},
			now,
		)

		if err = writeAPI.WritePoint(ctx, p); err != nil {
			log.Fatal(err)
		}

		if lcdDisplay != nil {
			log.Debug("Updating LCD")
			if err := lcdDisplay.WritePaddedString(fmt.Sprintf("Temp: %.1f C", response.Temperature), 1, 0); err != nil {
				log.Fatal(err)
			}
			if err := lcdDisplay.WritePaddedString(fmt.Sprintf("Press: %.1f hPa", response.Pressure), 2, 0); err != nil {
				log.Fatal(err)
			}
			if err := lcdDisplay.WritePaddedString(fmt.Sprintf("Hum: %.1f", response.Humidity), 3, 0); err != nil {
				log.Fatal(err)
			}
		}

		time.Sleep(10 * time.Second)
	}
}
