package tpcmds

import (
	"os"
	"fmt"
	flag "github.com/ogier/pflag"
	"../devices"
	"../poller"
)


func PollDevices() int {
	hosts := []string{"10.7.74.240:9999","10.7.74.241:9999","10.7.74.242:9999","10.7.74.243:9999"}

	influxCmds := flag.NewFlagSet("Loop Polling", flag.ExitOnError)

	host		:= influxCmds.StringP("host", "h", "localhost:8086", "Influx host to use. IP:PORT")
	database	:= influxCmds.StringP("database", "d", "", "Influx database to use")
	measurement	:= influxCmds.StringP("measurement", "m", "tplink", "Influx measurement to use")
	precision	:= influxCmds.StringP("precision", "p", "s", "Influx precision to use")
	rtpolicy	:= influxCmds.StringP("retention", "r", "autogen", "Influx retention policy")

	if len(os.Args) < 3 {
		influxCmds.PrintDefaults()
		return 1
	}
	influxCmds.Parse(os.Args[2:])
	if influxCmds.Parsed() {
		if *database == "" {
			influxCmds.PrintDefaults()
			return 1
		}
	}
	devices := tpdevices.TPDevices{}
	devices.AddHosts(hosts)

	if err := tppoller.RunPoll(devices, *host, *database, *measurement, *precision, *rtpolicy); err != nil {
		panic(err)
		return 1
	}
	fmt.Println("Points Written!")
	return 0
}

func LoopPollDevices() int {
	hosts := []string{"10.7.74.240:9999","10.7.74.241:9999","10.7.74.242:9999","10.7.74.243:9999"}
	
	pollCmds := flag.NewFlagSet("Loop Polling", flag.ExitOnError)

	interval	:= pollCmds.IntP("interval", "i", 5, "How often to poll for updates. Must be > 0")
	host		:= pollCmds.StringP("host", "h", "localhost:8086", "Influx host to use. IP:PORT")
	database	:= pollCmds.StringP("database", "d", "", "Influx database to use")
	measurement	:= pollCmds.StringP("measurement", "m", "tplink", "Influx measurement to use")
	precision	:= pollCmds.StringP("precision", "p", "s", "Influx precision to use")
	rtpolicy	:= pollCmds.StringP("retention", "r", "autogen", "Influx retention policy")

	if len(os.Args) < 3 {
		pollCmds.PrintDefaults()
		return 1
	}
	pollCmds.Parse(os.Args[2:])
	if pollCmds.Parsed() {
		if *interval < 1 {
			pollCmds.PrintDefaults()
			return 1
		}
		if *database == "" {
			pollCmds.PrintDefaults()
			return 1
		}
		devices := tpdevices.TPDevices{}
		devices.AddHosts(hosts)
		if err := tppoller.StartPolling(devices, uint8(*interval), *host, *database, *measurement, *precision, *rtpolicy); err != nil {
			panic(err)
			return 1
		}
	}
	return 1
}