package models

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/structs"
	"github.com/influxdata/influxdb1-client/v2"
	"log"
	"time"
)

const (
	// MyDB specifies name of database
	MyDB = "readings"
)

type TempReading struct {
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
	Battery     int64   `json:"batteryLevel"`
}

func NewTempReading(reading []byte) *TempReading {
	r := TempReading{}
	err := json.Unmarshal(reading, &r)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return &r
}

func Insert(reading *TempReading) {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: "http://localhost:8086",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  MyDB,
		Precision: "s",
	})
	if err != nil {
		log.Fatal(err)
	}

	// Create a point and add to batch
	// tags := map[string]string{"productView": productMeasurement["ProductName"].(string)}
	fields := structs.Map(reading)

	pt, err := client.NewPoint("sensor", nil, fields, time.Now())
	if err != nil {
		log.Fatal(err)
	}
	bp.AddPoint(pt)

	// Write the batch
	if err := c.Write(bp); err != nil {
		log.Fatal(err)
	}

	// Close client resources
	if err := c.Close(); err != nil {
		log.Fatal(err)
	}
}
