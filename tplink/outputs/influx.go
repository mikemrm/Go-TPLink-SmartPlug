package tpoutput

import (
	"os"
	"fmt"
	"flag"
	"time"
	"errors"
	"regexp"
	"syscall"
	"os/signal"
	"../devices"
//	"strings"
	"github.com/influxdata/influxdb/client/v2"
)

type Influx struct {
	Output
	Host		string
	Database	string
	Measurement	string
	Precision	string
	Retention	string
}

func (i *Influx) BuildPoints(devices tpdevices.TPDevices) (error, []*client.Point) {
	var points []*client.Point

	for _, device := range devices.GetDevices() {
		tags := make(map[string]string)
		fields := make(map[string]interface{})
		for f, v := range device.Data {
			if device.TagExists(f) {
				tags[f] = fmt.Sprintf("%v", v)
			} else {
				fields[f] = v
			}
		}
		pt, err := client.NewPoint(i.Measurement, tags, fields, time.Now())
		if err != nil {
			return err, points
		}
		points = append(points, pt)
	}
	return nil, points
}

func (i *Influx) WritePoints(points []*client.Point) error {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:	i.Host,
	})
	if err != nil {
		return err
	}
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:			i.Database,
		Precision:			i.Precision,
		RetentionPolicy:	i.Retention,
	})
	if err != nil {
		return err
	}
	for _, point := range points {
		bp.AddPoint(point)
	}
	return c.Write(bp)
}


func (i *Influx) Write(devices tpdevices.TPDevices) error {
	err, points := i.BuildPoints(devices)
	if err != nil {
		return err
	}
	if len(points) == 0 {
		return nil
	}
	return i.WritePoints(points)
}

func New(host, database, measurement, precision, retention string) (error, Influx) {
	influx := Influx{}
	if _, err := regexp.MatchString("^https?://[^:]+:[0-9]{1,4}$",  host); err != nil {
		return errors.New("Influx: Invalid host."), influx
	}
	if database == "" {
		return errors.New("Influx: Database not specified"), influx
	}
	switch precision {
		case "ns", "u", "ms", "s", "m", "h", "d", "w":
		default:
			return errors.New("Influx: Invalid precision."), influx
	}
	influx.Host = host
	influx.Database = database
	influx.Measurement = measurement
	influx.Precision = precision
	influx.Retention = retention
	return nil, influx
}

func NewFromArgs() (error, Influx) {
	influxCmds := flag.NewFlagSet("Influx Settings", flag.ExitOnError)

	host		:= influxCmds.String("influx.host", "http://localhost:8086", "Influx host to use. http://IP:PORT")
	database	:= influxCmds.String("influx.database", "", "Influx database to use")
	measurement	:= influxCmds.String("influx.measurement", "tpmon", "Influx measurement to use")
	precision	:= influxCmds.String("influx.precision", "s", "Influx precision to use")
	rtpolicy	:= influxCmds.String("influx.retention", "autogen", "Influx retention policy")
	
	influxCmds.Parse(os.Args[3:])
	if influxCmds.Parsed() {
		err, influx := New(*host, *database, *measurement, *precision, *rtpolicy)
		return err, influx
	}
	return errors.New("Invalid commands"), Influx{}
}

type InfluxLoop struct {
	*Influx
	Interval int
}

func (i *InfluxLoop) Write(devices tpdevices.TPDevices) error {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	is_polling := 0
	polling := make(chan int)
	ret_err := make(chan error)
	go func(devices tpdevices.TPDevices){
		for range time.Tick(time.Second * time.Duration(i.Interval)) {
			t := time.Now()
			polling <- 1
			s := time.Now()
			err, _ := devices.GetAllData()
			if err != nil {
				ret_err <- err
				break
			}
			err, points := i.BuildPoints(devices)
			if err != nil {
				ret_err <- err
				break
			}
			if len(points) == 0 {
				continue
			}
			if err = i.WritePoints(points); err != nil {
				ret_err <- err
				break
			}
			fmt.Printf("[%s]: Poll update completed in %s\n", t.Format("2006-01-02 15:04:05"), time.Now().Sub(s))
			polling <- 0
		}
		polling <- 0
	}(devices)
	Q:
	for {
		select {
			case x := <-c:
				fmt.Println("Caught", x)
				break Q
			case x := <-polling:
				is_polling = x
			case x := <-ret_err:
				if x != nil {
					fmt.Println("Found Error", x)
					return x
				}
		}
	}
	fmt.Println("Shutting down...")
	if is_polling == 0 {
		return nil
	}
	fmt.Println("Polling in progress, waiting for update.")
	running := <-polling
	if running == 0 {
		return nil
	}
	return nil
}

func NewLoopFromArgs() (error, InfluxLoop) {
	influxCmds := flag.NewFlagSet("Influx Settings", flag.ExitOnError)

	interval	:= influxCmds.Int("loop.interval", 5, "Influx loop interval")
	host		:= influxCmds.String("influx.host", "http://localhost:8086", "Influx host to use. http://IP:PORT")
	database	:= influxCmds.String("influx.database", "", "Influx database to use")
	measurement	:= influxCmds.String("influx.measurement", "tpmon", "Influx measurement to use")
	precision	:= influxCmds.String("influx.precision", "s", "Influx precision to use")
	rtpolicy	:= influxCmds.String("influx.retention", "autogen", "Influx retention policy")
	
	influxCmds.Parse(os.Args[3:])
	if influxCmds.Parsed() {
		err, influx := New(*host, *database, *measurement, *precision, *rtpolicy)
		return err, InfluxLoop{&influx, *interval}
	}
	return errors.New("Invalid commands"), InfluxLoop{}
}


func init() {
	AddOutput("influx", func() (error, Output){
		err, influx := NewFromArgs()
		return err, &influx
	})
	AddOutput("influx-loop", func() (error, Output){
		err, influx := NewLoopFromArgs()
		return err, &influx
	})
}