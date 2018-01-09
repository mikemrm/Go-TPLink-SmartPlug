package tpdevices

import (
	".."
	"encoding/json"
)

type SmartPlug struct {
	Address	string
}

func (d *SmartPlug) GetSystemInfo() (error, SystemInfo) {
	data := &SystemStructure{}
	err, results := tplink.Query(d.Address, data)
	if err != nil {
		return err, SystemInfo{}
	}
	if err := json.Unmarshal(results, &data); err != nil {
		return err, SystemInfo{}
	}
	return nil, data.System.Info
}

func (d *SmartPlug) GetRealTimeEnergy() (error, EnergyRealTime) {
	data := &EnergyStructure{}
	err, results := tplink.Query(d.Address, data)
	if err != nil {
		return err, EnergyRealTime{}
	}
	if err := json.Unmarshal(results, &data); err != nil {
		return err, EnergyRealTime{}
	}
	return nil, data.Energy.RealTime
}