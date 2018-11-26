package influx

import (
	"github.com/influxdata/influxdb/client/v2"
	"time"
)

type ClientConfig struct {
	Addr     string
	Username string
	Password string
	DB       string
}


func NewInfluxClient(config ClientConfig) (client.Client, error) {
	return client.NewHTTPClient(client.HTTPConfig{
		Addr:     config.Addr,
		Username: config.Username,
		Password: config.Password,
		Timeout:  15 * time.Second,
	})
}