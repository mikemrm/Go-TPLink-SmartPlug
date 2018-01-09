package main

import (
	"os"
	"fmt"
	"./tplink/tpprocs"
)

func PrintHelp() {
	help := `Usage:
    -l	Query devices
    -p	Poll data
    -P	Poll data constantly`
    fmt.Println(help)
}

func main() {
	if len(os.Args) < 2 {
		os.Exit(tpprocs.CMD_PrintDevices())
	}
	switch os.Args[1] {
		case "-l":
			os.Exit(tpprocs.CMD_PrintDevices())
		case "-p":
			os.Exit(tpprocs.CMD_PollDevices())
		case "-P":
			os.Exit(tpprocs.CMD_LoopPollDevices())
		default:
			PrintHelp()
			os.Exit(1)
	}
}