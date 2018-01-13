package tpdevices

import (
	".."
	"fmt"
	"encoding/json"
)

type TPDevices struct {
	hosts	[]string
	devices	[]TPDevice
}

func (ds *TPDevices) AddHost(host string) TPDevice {
	ds.hosts = append(ds.hosts, host)
	device := TPDevice{Addr: host}
	device.initialize()
	ds.devices = append(ds.devices, device)
	return device
}

func (ds *TPDevices) AddHosts(hosts []string) {
	for _, host := range hosts {
		ds.AddHost(host)
	}
}

func (ds *TPDevices) GetDevices() []TPDevice {
	return ds.devices
}

func (ds *TPDevices) GetAllData() (error, []TPDevice) {
	//running := make(map[string]bool, len(ds.devices))
	completed := 0
	updates := make(chan string)
	for _, device := range ds.devices {
		go func(d TPDevice) {
			d.GetAllData()
			updates <- d.Addr
		}(device)
	}
	Q:
	for {
		select {
			case <-updates:
				completed++
				if completed == len(ds.devices) {
					break Q
				}
		}
	}
	return nil, ds.devices
}

func DiscoverDevices(timeout int) (error, []tplink.Discovered, TPDevices) {
	devices := TPDevices{}
	var discovered []tplink.Discovered

	req := SystemStructure{}
	err, discovered := tplink.Discover(req, timeout)
	if err != nil {
		return err, discovered, devices
	} else {
		for _, d := range discovered {
			device := devices.AddHost(fmt.Sprintf("%s", d.Addr))
			data := SystemStructure{}
			json.Unmarshal(d.Data, &data)
			var tags []string
			device.appendData(tags, data.System.Info)
		}
	}
	return nil, discovered, devices
}