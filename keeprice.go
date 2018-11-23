package main

import (
	"github.com/GaruGaru/keeprice/api"
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

	keepriceApi := api.Api{
		Storage: priceStorage,
	}

	err = keepriceApi.Run("0.0.0.0:8976")

	if err != nil {
		panic(err)
	}

}

func WriteTestData(priceStorage storage.PriceStorage) {
	for i := 0; i < 1000; i++ {

		price := models.ProductPrice{
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
