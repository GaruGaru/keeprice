package main

import (
	"fmt"
	"github.com/GaruGaru/keeprice/cmd"
	"github.com/GaruGaru/keeprice/models"
	"github.com/GaruGaru/keeprice/storage"
	"math/rand"
	"time"
)

func main() {
	store := cmd.CassandraStorage()

//	WriteTestData(store)

	_, _ = store.Get("tannico", "vino-buono")

	cmd.Execute()
}

func WriteTestData(priceStorage storage.PriceStorage) {

	//err := priceStorage.Init()
	//
	//if err != nil {
	//	panic(err)
	//}

	fmt.Println("Start write")
	start := time.Now().UnixNano() / int64(time.Millisecond)
	for i := 0; i < 10; i++ {

		price := models.ProductPrice{
			SiteID:       "tannico",
			ProductID:    "vino-buono",
			ProductPrice: float32(i)+rand.Float32(),
			Time:         int64(i),
		}

		err := priceStorage.Store(price)
		if err != nil {
			fmt.Println(err.Error())
			//panic(err)
		}

	}

	delay := time.Now().UnixNano()/int64(time.Millisecond) - start
	fmt.Printf("Done write in %d ms\n", delay)

}
