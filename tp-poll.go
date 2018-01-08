package main

import (
	"fmt"
	"./tplink"
	"./tplink/tpdevices"
	"./renderers"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	Devices	[]string	`mapstructure:"devices"`
}

type Devices struct {
	Devices []*tpdevices.SmartPlug
}

func loadDevices(watch_changes bool) (error, *Config) {
	viper.SetConfigType("toml")
	viper.SetConfigName("tp-poll")
	viper.AddConfigPath("/etc/")
	viper.AddConfigPath(".")
	if watch_changes {
		viper.OnConfigChange(func(e fsnotify.Event){
			fmt.Println("Config file chagned:", e.Name)
		})
	}
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	var config Config
	err := viper.Unmarshal(&config)
	if err != nil {
		fmt.Println(err)
		return err, &config
	}
	return nil, &config
}
func (d *Devices) UpdateDeviceList(devices []string) *Devices {
	d.Devices = d.Devices[:0]
	for _, ip := range devices {
		d.Devices = append(d.Devices, &tpdevices.SmartPlug{ip})
	}
	return d
}

func (d *Devices) GetSystemInfo() (error, []tplink.SystemInfo) {
	var results []tplink.SystemInfo
	for _, device := range d.Devices {
		_, info := device.GetSystemInfo()
		results = append(results, info)
	}
	return nil, results
}

func (d *Devices) GetRealTimeEnergy() (error, []tpdevices.EnergyRealTime) {
	var results []tpdevices.EnergyRealTime
	for _, device := range d.Devices {
		_, info := device.GetRealTimeEnergy()
		results = append(results, info)
	}
	return nil, results
}

func main() {
	_, config := loadDevices(false)

	devices := Devices{}
	devices.UpdateDeviceList(config.Devices)
	_, sysinfos := devices.GetSystemInfo()
	_, power := devices.GetRealTimeEnergy()
	headers := []string{
		"Name", "IP", "Model", "Version",
		"Mode", "State", "Watts", "Current",
		"Voltage", "PF", "Total Kwh",
	}
	var rows [][]string
	for i, device := range devices.Devices {
		d_info := sysinfos[i]
		d_pow := power[i]
		relay_state := "Off"
		if d_info.RelayOn == 1 {
			relay_state = "On"
		}
		pf := d_pow.Power / (d_pow.Voltage * d_pow.Current)
		row := []string{
			d_info.Alias,
			device.Address,
			d_info.Model,
			d_info.HardwareVersion,
			d_info.Mode,
			relay_state,
			fmt.Sprintf("%0.2f", d_pow.Power),
			fmt.Sprintf("%0.2f", d_pow.Current),
			fmt.Sprintf("%0.2f", d_pow.Voltage),
			fmt.Sprintf("%0.2f", pf),
			fmt.Sprintf("%0.2f", d_pow.TotalKwh),
		}
		rows = append(rows, row)
	}

	renderers.AsciiTable(headers, rows)
}