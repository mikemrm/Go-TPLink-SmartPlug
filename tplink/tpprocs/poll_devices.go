package tpprocs

import (
	"fmt"
	"time"
	flag "github.com/ogier/pflag"
	"os"
	"os/signal"
	"github.com/influxdata/influxdb/client/v2"
)

func write_points(host string, database string, precision string, retention string, points []*client.Point) error {
	if len(points) == 0 {
		return nil
	}
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:	"http://" + host,
	})
	if err != nil {
		return err
	}
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:			database,
		Precision:			precision,
		RetentionPolicy:	retention,
	})
	if err != nil {
		return err
	}
	for _, point := range points {
		bp.AddPoint(point)
	}
	return c.Write(bp)
}

func runPoll(host string, database string, measurement string, precision string, retention string) int {
	_, config := loadDevices(false)

	devices := Devices{}
	devices.UpdateDeviceList(config.Devices)
	_, sysinfos := devices.GetSystemInfo()
	_, power := devices.GetRealTimeEnergy()

	var points []*client.Point

	for i, device := range devices.Devices {
		d_info := sysinfos[i]
		d_pow := power[i]
		tags := map[string]string {
			"name":	d_info.Alias,
			"host":	device.Address,
			"model": d_info.Model,
			"state": fmt.Sprintf("%d", d_info.RelayOn),
		}
		fields := map[string]interface{} {
			"watts": d_pow.Power,
			"amps": d_pow.Current,
			"volts": d_pow.Voltage,
			"kwh": d_pow.TotalKwh,
		}
		pt, err := client.NewPoint(measurement, tags, fields, time.Now())
		if err != nil {
			fmt.Println(err)
			return 1
		}
		points = append(points, pt)
	}
	if err := write_points(host, database, precision, retention, points); err != nil {
		fmt.Println(err)
		return 1
	}
	return 0
}

func CMD_PollDevices() int {
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
	ecode := runPoll(*host, *database, *measurement, *precision, *rtpolicy)
	if ecode == 0 {
		fmt.Println("Points Written!")
	}
	return ecode
}

func startPolling(interval uint8, host string, database string, measurement string, precision string, retention string) int {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	is_polling := 0
	polling := make(chan int)
	go func(){
		for range time.Tick(time.Second * time.Duration(interval)) {
			t := time.Now()
			polling <- 1
			fmt.Printf("[%s]: Polling...\n", t.Format("2006-01-02 15:04:05"))
			if ecode := runPoll(host, database, measurement, precision, retention); ecode != 0 {
				fmt.Println("Error:", ecode)
			}
			polling <- 0
		}
	}()
	Q:
	for {
		select {
			case <-c:
				break Q
			case x := <-polling:
				is_polling = x
		}
	}
	fmt.Println("Shutting down...")
	if is_polling == 0 {
		return 0
	}
	fmt.Println("Polling in progress, waiting for update.")
	running := <-polling
	if running == 0 {
		return 0
	}
	return 0
}

func CMD_LoopPollDevices() int {
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
		return startPolling(uint8(*interval), *host, *database, *measurement, *precision, *rtpolicy)
	}
	return 1
}