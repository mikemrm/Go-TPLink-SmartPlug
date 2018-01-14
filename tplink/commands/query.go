package tpcmds

import (
	"fmt"
	"../devices"
	"../outputs"
	"errors"
)

func Query(output tpoutput.Output) error {
	fmt.Println("Discovering devices...")
	err, discovered, devices := tpdevices.DiscoverDevices(1)
	if err != nil {
		return err
	}
	if len(discovered) > 0 {
		fmt.Println("Found", len(discovered), "devices!")
	} else {
		return errors.New("No devices found.")
	}

	if err, _ = devices.GetAllData(); err != nil {
		return err
	}
	err = output.Write(devices)
	return err
}