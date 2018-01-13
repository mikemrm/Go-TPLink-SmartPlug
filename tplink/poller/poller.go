package tppoller

import (
	"fmt"
	"syscall"
	"time"
	"os"
	"os/signal"
	"github.com/influxdata/influxdb/client/v2"
	"../devices"
	"../outputs"
)

func RunPoll(devices tpdevices.TPDevices, host string, database string, measurement string, precision string, retention string) error {
	var points []*client.Point

	err, data_devices := devices.GetAllData()

	if err != nil {
		return err
	}

	for _, device := range data_devices {
		tags := make(map[string]string)
		fields := make(map[string]interface{})
		for f, v := range device.Data {
			if device.TagExists(f) {
				tags[f] = fmt.Sprintf("%v", v)
			} else {
				fields[f] = v
			}
		}
		pt, err := client.NewPoint(measurement, tags, fields, time.Now())
		if err != nil {
			fmt.Println(err)
			return err
		}
		points = append(points, pt)
	}
	if err := tpoutput.Influx(host, database, precision, retention, points); err != nil {
		return err
	}
	return nil
}

func StartPolling(devices tpdevices.TPDevices, interval uint8, host string, database string, measurement string, precision string, retention string) error {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	is_polling := 0
	polling := make(chan int)
	ret_err := make(chan error)
	go func(devices tpdevices.TPDevices){
		for range time.Tick(time.Second * time.Duration(interval)) {
			t := time.Now()
			polling <- 1
			fmt.Printf("[%s]: Polling...\n", t.Format("2006-01-02 15:04:05"))
			if err := RunPoll(devices, host, database, measurement, precision, retention); err != nil {
				ret_err <- err
				break
			}
			polling <- 0
		}
		polling <- 0
	}(devices)
	Q:
	for {
		select {
			case x := <-c:
				fmt.Println("Caught", x)
				break Q
			case x := <-polling:
				is_polling = x
			case x := <-ret_err:
				if x != nil {
					fmt.Println("Found Error")
					return x
				}
		}
	}
	fmt.Println("Shutting down...")
	if is_polling == 0 {
		return nil
	}
	fmt.Println("Polling in progress, waiting for update.")
	running := <-polling
	if running == 0 {
		return nil
	}
	return nil
}