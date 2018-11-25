package main

import (
	"fmt"
	"github.com/GaruGaru/keeprice/api"
	"github.com/GaruGaru/keeprice/models"
	"github.com/GaruGaru/keeprice/storage"
	"github.com/GaruGaru/keeprice/storage/cassandra"
	"github.com/GaruGaru/keeprice/storage/influx"
	"github.com/gocql/gocql"
	"log"
	"time"
)

const (
	influxAddress  = "http://localhost:8086/"
	influxDB       = "keeprice"
	influxUsername = "keeprice"
	influxPassword = "keeprice-password"
)

var cassandraHosts = []string{"localhost:9042", "localhost:9043", "localhost:9044"}

const (
	cassandraUsername = "cassandra"
	cassandraPassword = "keeprice-cassandra-password"
)

func InfluxDBStorage() influx.InfluxDBStorage {

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

func CassandraStorage() cassandra.CassandraStorage {
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

func main() {

	priceStorage := CassandraStorage()

	keepriceApi := api.Api{
		Storage: priceStorage,
	}

	WriteTestData(priceStorage)

	err := keepriceApi.Run("0.0.0.0:8976")

	if err != nil {
		panic(err)
	}

}

func TestCassandra() {
	cluster := gocql.NewCluster("localhost:9042", "localhost:9043", "localhost:9044")
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: "cassandra",
		Password: "keeprice-cassandra-password",
	}

	cluster.Consistency = gocql.Quorum
	session, err := cluster.CreateSession()

	if err != nil {
		panic(err)
	}
	defer session.Close()

	err = session.Query(fmt.Sprintf(`CREATE KEYSPACE IF NOT EXISTS %s
	WITH replication = {
		'class' : 'SimpleStrategy',
		'replication_factor' : 2
	}`, "prices")).Exec()

	if err != nil {
		panic(err)
	}

	//if err := session.Query(`create table if not exists prices.tweet(timeline text, id UUID, text text, PRIMARY KEY(id));`).Exec(); err != nil {
	///	log.Fatal(err)
	//è }

	if err := session.Query(`INSERT INTO "prices"."tweet" (timeline, id, text) VALUES (?, ?, ?)`,
		"me", gocql.TimeUUID(), "hello world").Exec(); err != nil {
		log.Fatal(err)
	}

	var id gocql.UUID
	var text string

	/* Search for a specific set of records whose 'timeline' column matches
	 * the value 'me'. The secondary index that we created earlier will be
	 * used for optimizing the search */
	if err := session.Query(`SELECT id, text FROM prices.tweet WHERE timeline = ? LIMIT 1`,
		"me").Consistency(gocql.One).Scan(&id, &text); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Tweet:", id, text)

	// list all tweets
	iter := session.Query(`SELECT id, text FROM prices.tweet WHERE timeline = ?`, "me").Iter()
	for iter.Scan(&id, &text) {
		fmt.Println("Tweet:", id, text)
	}
	if err := iter.Close(); err != nil {
		log.Fatal(err)
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
