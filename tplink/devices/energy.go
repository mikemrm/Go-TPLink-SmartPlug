package tpdevices

import ".."
import "encoding/json"

type RealTimeEnergy struct {
	ErrorCode	uint8	`json:"err_code,omitempty"`
	Power		float32	`json:"power,omitempty"`
	Voltage		float32	`json:"voltage,omitempty"`
	Current		float32	`json:"current,omitempty"`
	TotalKwh	float32	`json:"total,omitempty"`
}

type Energy struct {
	RealTime RealTimeEnergy `json:"get_realtime"`
}

type EnergyStructure struct {
	Energy	Energy	`json:"emeter,omitempty"`
}

func (d *TPDevice) GetRealTimeEnergy() (error, []string, RealTimeEnergy) {
	tags := []string{}
	data := &EnergyStructure{}
	err, resp := tplink.Query(d.Addr, data)
	if err != nil {
		return err, tags, RealTimeEnergy{}
	}
	if err := json.Unmarshal(resp, &data); err != nil {
		return err, tags, RealTimeEnergy{}
	}
	return nil, tags, data.Energy.RealTime
}