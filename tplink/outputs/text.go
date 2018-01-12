package tpoutput
import (
	"fmt"
	"strings"
)

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
