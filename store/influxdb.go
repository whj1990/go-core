package store

import (
	"github.com/whj1990/go-core/config"
	client "github.com/influxdata/influxdb1-client/v2"
)

func NewInfluxDB() (client.Client, error) {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     config.GetString("influxdb.url", ""),
		Username: config.GetString("influxdb.username", ""),
		Password: config.GetString("influxdb.password", ""),
	})
	if err != nil {
		return nil, err
	}
	return c, nil
}
