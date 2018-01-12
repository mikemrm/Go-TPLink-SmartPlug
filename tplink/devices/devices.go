package tpdevices

type TPDevices struct {
	hosts	[]string
	devices	[]TPDevice
}

func (ds *TPDevices) AddHost(host string) {
	ds.hosts = append(ds.hosts, host)
	device := TPDevice{Addr: host}
	device.initialize()
	ds.devices = append(ds.devices, device)
}

func (ds *TPDevices) AddHosts(hosts []string) {
	for _, host := range hosts {
		ds.AddHost(host)
	}
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