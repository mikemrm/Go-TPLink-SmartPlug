package tpdevices

type EnergyRealTime struct {
	ErrorCode	uint8	`json:"err_code,omitempty"`
	Power		float32	`json:"power,omitempty"`
	Voltage		float32	`json:"voltage,omitempty"`
	Current		float32	`json:"current,omitempty"`
	TotalKwh	float32	`json:"total,omitempty"`
}

type Energy struct {
	RealTime EnergyRealTime `json:"get_realtime"`
}

type EnergyStructure struct {
	Energy	Energy	`json:"emeter,omitempty"`
}