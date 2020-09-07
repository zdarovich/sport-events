package influxdb

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/zdarovich/sport-events/config"
	"log"
	"sync"
	"sync/atomic"
	"time"

	_ "github.com/influxdata/influxdb1-client"
	influx "github.com/influxdata/influxdb1-client/v2"
)

var (
	influxMu     sync.Mutex
	influxInit   uint32
	influxClient influx.Client
)
var dbName = config.Config.Influxdb.Name

type Response = influx.Response

func newClient() (influx.Client, error) {
	influxAddress := config.Config.Influxdb.Address
	influxTimeout := time.Duration(config.Config.Influxdb.Timeout) * time.Second
	config := influx.HTTPConfig{Addr: influxAddress, Timeout: influxTimeout}
	client, err := influx.NewHTTPClient(config)
	if err != nil {
		return nil, err
	}

	_, _, err = client.Ping(time.Duration(0))
	if err != nil {
		return nil, err
	}

	return client, nil
}

// Singleton instance of InfluxDB connection
func GetClient() (influx.Client, error) {

	if atomic.LoadUint32(&influxInit) == 1 {
		return influxClient, nil
	}

	influxMu.Lock()
	defer influxMu.Unlock()

	if influxInit == 0 {
		client, err := newClient()
		if err != nil {
			return nil, err
		}

		influxClient = client
		atomic.StoreUint32(&influxInit, 1)
	}

	return influxClient, nil
}

func createDatabase(client influx.Client) (*influx.Response, error) {
	query := influx.NewQuery(fmt.Sprintf("CREATE DATABASE %s", dbName), "", "")
	return client.Query(query)
}

func Initialize(newDbName string) error {
	if newDbName != "" {
		dbName = newDbName
	}
	client, err := GetClient()
	if err != nil {
		return err
	}

	_, err = createDatabase(client)
	if err != nil {
		return err
	}
	return nil
}

func AddNewPoint(name string, tags map[string]string, fields map[string]interface{}, t time.Time) error {
	client, err := GetClient()
	if err != nil {
		return err
	}

	batchPointsConfig := influx.BatchPointsConfig{Database: dbName}
	batchPoints, err := influx.NewBatchPoints(batchPointsConfig)
	if err != nil {
		return err
	}
	point, err := influx.NewPoint(name, tags, fields, t)
	if err != nil {
		return err
	}
	batchPoints.AddPoint(point)
	err = client.Write(batchPoints)
	return err
}

func Query(queryCommand string) (*Response, error) {
	client, err := GetClient()
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, errors.Wrap(err, "Database connection failed")
	}
	query := influx.NewQuery(queryCommand, dbName, "")
	response, err := client.Query(query)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, errors.Wrap(err, "Cannot perform database query")
	}
	return response, nil
}
