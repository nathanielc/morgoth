package influxdb

import (
	"fmt"
	"github.com/influxdb/influxdb/client"
)

func connect(config *InfluxDBConf) (*client.Client, error) {
	c := client.ClientConfig{
		Host:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Username: config.User,
		Password: config.Password,
		Database: config.Database,
	}

	return client.NewClient(&c)
}
