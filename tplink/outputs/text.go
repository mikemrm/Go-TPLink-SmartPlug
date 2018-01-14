package tpoutput
import (
	"fmt"
	"strings"
	"../devices"
)

type TableBrief struct {
	something string
}

func (t *TableBrief) Write(devices tpdevices.TPDevices) error {
	headers := []string{
		"Name", "IP", "Model", "Version",
		"Mode", "State", "Watts", "Current",
		"Voltage", "PF", "Total Kwh",
	}
	var rows [][]string

	for _, device := range devices.GetDevices() {
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

	AsciiTable(headers, rows)
	return nil
}
type Table struct {}

func (t *Table) Write(devices tpdevices.TPDevices) error {
	fmt.Println(devices)
	return nil
}

func init(){
	AddOutput("brief", func() (error, Output) {
		return nil, &TableBrief{}
	})
	AddOutput("table", func() (error, Output) {
		return nil, &Table{}
	})
}

func AsciiTable(header []string, rows [][]string) {
	col_widths := make([]int, len(header))
	for i, s := range header {
		if len(s) > col_widths[i] {
			col_widths[i] = len(s)
		}
	}
	for _, row := range rows {
		for i, s := range row {
			if len(s) > col_widths[i] {
				col_widths[i] = len(s)
			}
		}
	}

	var separations []string
	for _, count := range col_widths {
		separations = append(separations, strings.Repeat("-", count))

	}
	separator := "+-" + strings.Join(separations, "-+-") + "-+"

	var padded_header []string
	for i, col := range header {
		diff := col_widths[i] - len(col)
		if diff > 0 {
			padded_header = append(padded_header, col + strings.Repeat(" ", diff))
		} else {
			padded_header = append(padded_header, col)
		}
	}

	fmt.Println(separator)
	fmt.Println("| " + strings.Join(padded_header, " | ") + " |")
	fmt.Println(separator)
	for _, row := range rows {
		var padded []string
		for i, col := range row {
			diff := col_widths[i] - len(col)
			if diff > 0 {
				padded = append(padded, col + strings.Repeat(" ", diff))
			} else {
				padded = append(padded, col)
			}
		}
		fmt.Println("| " + strings.Join(padded, " | ") + " |")
	}
	fmt.Println(separator)
}
