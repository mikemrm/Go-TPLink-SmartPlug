package tpdevices

import ".."
import "encoding/json"

type SystemInfo struct {
	Mode			string	`json:"active_mode,omitempty"`
	Alias			string	`json:"alias,omitempty"`
	Product			string	`json:"dev_name,omitempty"`
	DeviceId		string	`json:"device_id,omitempty"`
	ErrorCode		int		`json:"err_code,omitempty"`
	Features		string	`json:"feature,omitempty"`
	FirmwareId		string	`json:"fwId,omitempty"`
	HardwareId		string	`json:"hwId,omitempty"`
	HardwareVersion	string	`json:"hw_ver,omitempty"`
	IconHash		string	`json:"icon_hash,omitempty"`
	GpsLatitude		float32	`json:"latitude,omitempty"`
	GpsLongitude	float32	`json:"longitude,omitempty"`
	LedOff			uint8	`json:"led_off,omitempty"`
	Mac				string	`json:"mac,omitempty"`
	Model			string	`json:"model,omitempty"`
	OemId			string	`json:"odemId,omitempty"`
	OnTime			uint32	`json:"on_time,omitempty"`
	RelayOn			uint8	`json:"relay_state,omitempty"`
	Rssi			int		`json:"rssi,omitempty"`
	SoftwareVersion	string	`json:"sw_ver,omitempty"`
	ProductType		string	`json:"type,omitempty"`
	Updating		uint8	`json:"updating,omitempty"`
}

type System struct {
	Info SystemInfo `json:"get_sysinfo"`
}

type SystemStructure struct {
	System	System	`json:"system,omitempty"`
}

func (d *TPDevice) getSystemInfo() (error, []string, SystemInfo) {
	tags := []string{
		"Mac", "Alias", "Features", "FirmwareId",
		"GpsLatitude", "OemId", "SoftwareVersion", "Mode", "HardwareVersion",
		"OnTime", "RelayOn", "ProductType", "Updating", "Model", "HardwareId",
		"GpsLongitude", "LedOff", "Rssi",
	}
	data := &SystemStructure{}
	err, resp := tplink.Query(d.Addr, data)
	if err != nil {
		return err, tags, SystemInfo{}
	}
	if err := json.Unmarshal(resp, &data); err != nil {
		return err, tags, SystemInfo{}
	}
	return nil, tags, data.System.Info
}