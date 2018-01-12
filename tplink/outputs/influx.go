package tpoutput
import (
	"github.com/influxdata/influxdb/client/v2"
)

func Influx(host string, database string, precision string, retention string, points []*client.Point) error {
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