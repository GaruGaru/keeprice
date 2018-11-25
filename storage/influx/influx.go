package influx

import (
	"encoding/json"
	"fmt"
	"github.com/GaruGaru/keeprice/models"
	"github.com/influxdata/influxdb/client/v2"
	"time"
)

func NewInfluxClient(config ClientConfig) (client.Client, error) {
	return client.NewHTTPClient(client.HTTPConfig{
		Addr:     config.Addr,
		Username: config.Username,
		Password: config.Password,
		Timeout:  15 * time.Second,
	})
}

type InfluxDBStorage struct {
	client   client.Client
	database string
}

func NewInfluxDBStorage(config ClientConfig) (InfluxDBStorage, error) {
	influxClient, err := NewInfluxClient(config)
	if err != nil {
		return InfluxDBStorage{}, err
	}
	return InfluxDBStorage{
		client:   influxClient,
		database: config.DB,
	}, nil
}

func (s InfluxDBStorage) Init() error {
	return nil
}

func (s InfluxDBStorage) Store(itemPrice models.ProductPrice) error {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  s.database,
		Precision: "s",
	})

	if err != nil {
		return err
	}

	tags := map[string]string{
		"site_id":    itemPrice.SiteID,
		"product_id": itemPrice.ProductID,
	}

	fields := make(map[string]interface{})

	fields["price"] = itemPrice.ProductPrice

	pt, err := client.NewPoint("prices", tags, fields, time.Unix(itemPrice.Time, 0))

	if err != nil {
		return err
	}

	bp.AddPoint(pt)

	if err := s.client.Write(bp); err != nil {
		return err
	}

	return nil
}

func (s InfluxDBStorage) Get(siteID string, productID string) (models.ProductPriceHistory, error) {

	EmptyHistory := models.ProductPriceHistory{
		History: []models.ProductPriceHistoryEntry{},
	}

	//queryTemplate := `SELECT "price" FROM "%s"."autogen"."prices" WHERE "product_id"='%s' AND "site_id"='%s' GROUP BY time(1d) ORDER BY time DESC`
	queryTemplate := `SELECT mean("price") FROM "%s"."autogen"."prices" WHERE "product_id"='%s' AND "site_id"='%s' GROUP BY time(10s) ORDER BY time DESC`
	query := fmt.Sprintf(queryTemplate, s.database, productID, siteID)

	response, err := s.client.Query(client.NewQuery(query, s.database, "s"))

	if err != nil {
		return EmptyHistory, err
	}

	if len(response.Results[0].Series) == 0 {
		return EmptyHistory, nil
	}

	series := response.Results[0].Series[0]

	if len(series.Columns) == 0 {
		return EmptyHistory, nil
	}

	var results []models.ProductPriceHistoryEntry

	for _, val := range series.Values {
		if val[1] == nil {
			continue
		}

		price, err := val[1].(json.Number).Float64()

		if err != nil {
			return EmptyHistory, err
		}

		timestamp, err := val[0].(json.Number).Int64()

		if err != nil {
			return EmptyHistory, err
		}

		results = append(results, models.ProductPriceHistoryEntry{
			Time:  timestamp,
			Price: float32(price),
		})

	}
	return models.ProductPriceHistory{
		Count:       len(results),
		PeriodStart: results[len(results)-1].Time,
		PeriodEnd:   results[0].Time,
		History:     results,
	}, nil
}
