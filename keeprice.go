package main

import (
	"fmt"
	"github.com/GaruGaru/keeprice/models"
	"github.com/GaruGaru/keeprice/storage"
	"time"
)

const (
	influxAddress  = "http://localhost:8086/"
	influxDB       = "keeprice"
	influxUsername = "keeprice"
	influxPassword = "keeprice-password"
)

func main() {

	influxConfig := storage.InfluxClientConfig{
		Addr:     influxAddress,
		Username: influxUsername,
		Password: influxPassword,
		DB:       influxDB,
	}

	priceStorage, err := storage.NewInfluxDBStorage(influxConfig)

	if err != nil {
		panic(err)
	}

	result, err := priceStorage.Get("tannico", "vino-buono")

	if err != nil {
		panic(err)
	}

	fmt.Printf("got %d results\n", result.Count)

	for _, r := range result.History {
		fmt.Printf("%f - %s\n", r.Price, time.Unix(r.Time, 0).String())
	}

}

func WriteTestData(priceStorage storage.PriceStorage) {
	for i := 0; i < 1000; i++ {

		price := models.ItemPrice{
			SiteID:       "tannico",
			ProductID:    "vino-buono",
			ProductPrice: float32(i),
			Time:         time.Now().Unix(),
		}

		err := priceStorage.Store(price)
		if err != nil {
			panic(err)
		}

		time.Sleep(1 * time.Second)

	}
}
