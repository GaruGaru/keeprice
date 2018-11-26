package cassandra

import (
	"github.com/GaruGaru/keeprice/models"
	"testing"
	"time"
)

var TestConfig = Config{
	Hosts:    []string{"127.0.0.1:9042"},
	Username: "cassandra",
	Password: "keeprice-cassandra-password",
	KeySpace: "keeprice_test",
	Table:    "prices_test",
}

func TestCassandraInsert(t *testing.T) {
	storage, err := NewCassandraStorage(TestConfig)

	if err != nil {
		t.Fatal(err)
	}

	err = storage.Init()

	if err != nil {
		t.Fatal(err)
	}

	err = storage.Store(models.ProductPrice{
		SiteID:       "test_site_id",
		ProductID:    "test_product_id",
		ProductPrice: 1.5,
		Time:         time.Now().Unix(),
	})

	if err != nil {
		t.Fatal(err)
	}

}


func TestCassandraInsertRead(t *testing.T) {
	t.Parallel()
	storage, err := NewCassandraStorage(TestConfig)

	if err != nil {
		t.Fatal(err)
	}

	err = storage.Init()

	if err != nil {
		t.Fatal(err)
	}

	err = storage.Store(models.ProductPrice{
		SiteID:       "test_site_id",
		ProductID:    "test_product_id",
		ProductPrice: 1.5,
		Time:         time.Now().Unix(),
	})

	if err != nil {
		t.Fatal(err)
	}

}
