package tpcmds

import (
	"fmt"
	"../devices"
	"../outputs"
)

func PrintDevices() int {

	fmt.Println("Discovering devices...")
	err, discovered, devices := tpdevices.DiscoverDevices(1)
	if err != nil {
		fmt.Println(err)
		return 1
	}
	if len(discovered) > 0 {
		fmt.Println("Found", len(discovered), "devices!")
	} else {
		fmt.Println("No devices found.")
		return 1
	}

	err, data := devices.GetAllData()
	if err != nil {
		panic(err)
	}

	headers := []string{
		"Name", "IP", "Model", "Version",
		"Mode", "State", "Watts", "Current",
		"Voltage", "PF", "Total Kwh",
	}
	var rows [][]string

	for _, device := range data {
		v := device.Data
		relay_state := "Off"
		if v["RelayOn"] == 1 {
			relay_state = "On"
		}
		pf := v["Power"].(float32) / (v["Voltage"].(float32) * v["Current"].(float32))
		row := []string{
			v["Alias"].(string),
			device.Addr,
			v["Model"].(string),
			v["HardwareVersion"].(string),
			v["Mode"].(string),
			relay_state,
			fmt.Sprintf("%0.2f", v["Power"].(float32)),
			fmt.Sprintf("%0.2f", v["Current"].(float32)),
			fmt.Sprintf("%0.2f", v["Voltage"].(float32)),
			fmt.Sprintf("%0.2f", pf),
			fmt.Sprintf("%0.2f", v["TotalKwh"].(float32)),
		}
		rows = append(rows, row)
	}

	tpoutput.AsciiTable(headers, rows)
	return 0
}