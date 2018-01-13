package main

import (
	"os"
	"fmt"
	"./tplink/commands"
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
		os.Exit(tpcmds.PrintDevices())
	}
	switch os.Args[1] {
		case "-l":
			os.Exit(tpcmds.PrintDevices())
		case "-p":
			os.Exit(tpcmds.PollDevices())
		case "-P":
			os.Exit(tpcmds.LoopPollDevices())
		default:
			PrintHelp()
			os.Exit(1)
	}
}