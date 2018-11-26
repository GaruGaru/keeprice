package cmd

import (
	"fmt"
	"github.com/GaruGaru/keeprice/api"
	"github.com/GaruGaru/keeprice/storage"
	"github.com/GaruGaru/keeprice/storage/cassandra"
	"github.com/GaruGaru/keeprice/storage/influx"
	"github.com/spf13/cobra"
	"strings"
)

var addr string
var backendStorage string

var (
	influxAddress  string
	influxDB       string
	influxUsername string
	influxPassword string
)

var (
	cassandraNodes    []string
	cassandraNodesStr string
	cassandraUsername string
	cassandraPassword string
)

func init() {
	apiCmd.Flags().StringVar(&addr, "bind", "0.0.0.0:8976", "api interface address binding")

	apiCmd.Flags().StringVar(&backendStorage, "storage", "influxdb", "backend storage type")

	apiCmd.Flags().StringVar(&influxAddress, "influx-address", "http://localhost:8086/", "influx backend db address")
	apiCmd.Flags().StringVar(&influxDB, "influx-db", "keeprice", "influx backend db name")
	apiCmd.Flags().StringVar(&influxUsername, "influx-username", "keeprice", "influx backend db username")
	apiCmd.Flags().StringVar(&influxPassword, "influx-password", "keeprice-password", "influx backend db password")

	apiCmd.Flags().StringVar(&cassandraUsername, "cassandra-username", "cassandra", "cassandra db username")
	apiCmd.Flags().StringVar(&cassandraPassword, "cassandra-password", "keeprice-password", "cassandra db password")
	apiCmd.Flags().StringVar(&cassandraNodesStr, "cassandra-nodes", "localhost:9042", "cassandra nodes separated by comma. eg: localhost:9042,localhost:9043...")

	rootCmd.AddCommand(apiCmd)
}

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "start keeprice api server",
	Run: func(cmd *cobra.Command, args []string) {

		cassandraNodes = strings.Split(cassandraNodesStr, ",")

		var priceStorage storage.PriceStorage

		switch backendStorage {
		case "influxdb":
			priceStorage = InfluxDBStorage()
			break
		case "cassandra":
			priceStorage = CassandraStorage()
			break
		default:
			fmt.Printf("unsupported storage type %s", backendStorage)
		}

		err := priceStorage.Init()

		if err != nil {
			panic(err)
		}

		keepriceApi := api.Api{Storage: priceStorage}

		err = keepriceApi.Run(addr)

		if err != nil {
			panic(err)
		}
	},
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
		Hosts:    cassandraNodes,
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
