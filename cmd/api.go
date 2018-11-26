package cmd

import (
	"github.com/GaruGaru/keeprice/api"
	"github.com/GaruGaru/keeprice/storage/cassandra"
	"github.com/GaruGaru/keeprice/storage/influx"
	"github.com/spf13/cobra"
)

const (
	influxAddress  = "http://localhost:8086/"
	influxDB       = "keeprice"
	influxUsername = "keeprice"
	influxPassword = "keeprice-password"
)

var cassandraHosts = []string{"localhost:9042"}

const (
	cassandraUsername = "cassandra"
	cassandraPassword = "keeprice-cassandra-password"
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "start keeprice api server",
	Run: func(cmd *cobra.Command, args []string) {

		priceStorage := CassandraStorage()

		err := priceStorage.Init()

		if err != nil {
			panic(err)
		}

		keepriceApi := api.Api{Storage: priceStorage}

		err = keepriceApi.Run("0.0.0.0:8976")

		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(apiCmd)
}

func InfluxDBStorage() influx.Storage {

	influxConfig := influx.ClientConfig{
		Addr:     influxAddress,
		Username: influxUsername,
		Password: influxPassword,
		DB:       influxDB,
	}

	influxStorage, err := influx.NewInfluxDBStorage(influxConfig)

	if err != nil {
		panic(err)
	}

	return influxStorage
}

func CassandraStorage() cassandra.Storage {
	config := cassandra.Config{
		Hosts:    cassandraHosts,
		Username: cassandraUsername,
		Password: cassandraPassword,
		KeySpace: "keeprice",
		Table:    "prices",
	}

	cassandraStorage, err := cassandra.NewCassandraStorage(config)

	if err != nil {
		panic(err)
	}

	return cassandraStorage
}
