package store

import (
	client "github.com/influxdata/influxdb1-client/v2"
	"github.com/whj1990/go-core/config"
)

func NewInfluxDB() (client.Client, error) {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     config.GetNaCosString("influxdb.url", ""),
		Username: config.GetNaCosString("influxdb.username", ""),
		Password: config.GetNaCosString("influxdb.password", ""),
	})
	if err != nil {
		return nil, err
	}
	return c, nil
}
