package store

import (
	client "github.com/influxdata/influxdb1-client/v2"
	"github.com/whj1990/go-core/config"
)

func NewInfluxDB() (client.Client, error) {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     config.GetNacosConfigData().Influxdb.Url,
		Username: config.GetNacosConfigData().Influxdb.Username,
		Password: config.GetNacosConfigData().Influxdb.Password,
	})
	if err != nil {
		return nil, err
	}
	return c, nil
}
