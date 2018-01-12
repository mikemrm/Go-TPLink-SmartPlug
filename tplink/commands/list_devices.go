package tpcmds

import (
	"time"
	"fmt"
	"../devices"
	"../outputs"
)

func PrintDevices() int {
	s := time.Now()
	hosts := []string{"10.7.74.240:9999","10.7.74.241:9999","10.7.74.242:9999","10.7.74.243:9999"}

	devices := tpdevices.TPDevices{}
	for _, host := range hosts {
		devices.AddHost(host)
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
	e := time.Now()
	d := e.Sub(s)
	fmt.Println(d)
	return 0
}