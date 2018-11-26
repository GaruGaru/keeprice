package cassandra

import (
	"fmt"
	"github.com/GaruGaru/keeprice/models"
	"github.com/gocql/gocql"
	"time"
)

var EmptyResult = models.ProductPriceHistory{}

type Storage struct {
	session           gocql.Session
	keyspace          string
	table             string
	replicationFactor int
}

func NewCassandraStorage(config Config) (Storage, error) {

	if config.Consistency == 0 {
		config.Consistency = gocql.One
	}

	if config.ReplicationFactor == 0 {
		config.ReplicationFactor = 1
	}

	cluster := gocql.NewCluster(config.Hosts...)

	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: config.Username,
		Password: config.Password,
	}

	cluster.Consistency = config.Consistency
	cluster.Timeout = 15 * time.Second
	session, err := cluster.CreateSession()

	if err != nil {
		return Storage{}, err
	}

	return Storage{
		session:           *session,
		keyspace:          config.KeySpace,
		table:             config.Table,
		replicationFactor: config.ReplicationFactor,
	}, nil
}

func (s Storage) Init() error {
	if err := s.createKeyspaceIfNotExists(); err != nil {
		return err
	}
	if err := s.createTableIfNotExists(); err != nil {
		return err
	}
	return nil
}

func (s Storage) Store(itemPrice models.ProductPrice) error {
	queryTemplate := `INSERT INTO %s (site_id, product_id, time, price) VALUES ('%s', '%s', %d, %f)`
	query := fmt.Sprintf(queryTemplate, s.tableName(), itemPrice.SiteID, itemPrice.ProductID, itemPrice.Time, itemPrice.ProductPrice)
	return s.session.Query(query).Exec()
}

func (s Storage) Get(siteID string, productID string) (models.ProductPriceHistory, error) {
	queryTemplate := `SELECT time, price FROM %s WHERE site_id='%s' and product_id='%s' ORDER BY time DESC;`
	queryStr := fmt.Sprintf(queryTemplate, s.tableName(), siteID, productID)
	query := s.session.Query(queryStr)

	var results []models.ProductPriceHistoryEntry

	var t int64
	var price float32
	iter := query.Iter()
	for iter.Scan(&t, &price) {
		results = append(results, models.ProductPriceHistoryEntry{
			Time:  t,
			Price: price,
		})
	}

	err := iter.Close()
	if err != nil {
		return EmptyResult, err
	}

	return models.ProductPriceHistory{
		Count:       len(results),
		PeriodStart: results[len(results)-1].Time,
		PeriodEnd:   results[0].Time,
		History:     results,
	}, nil
}

func (s Storage) createKeyspaceIfNotExists() error {
	queryTemplate := `CREATE KEYSPACE IF NOT EXISTS %s
    WITH replication = {
        'class' : 'SimpleStrategy',
        'replication_factor' : %d
    }`
	query := fmt.Sprintf(queryTemplate, s.keyspace, s.replicationFactor)
	return s.session.Query(query).Exec()
}

func (s Storage) createTableIfNotExists() error {
	queryTemplate := `CREATE TABLE IF NOT EXISTS %s(
     site_id text,
     product_id text,
     time timestamp,
     price float,
     PRIMARY KEY((site_id, product_id), time)
    );`
	query := fmt.Sprintf(queryTemplate, s.tableName())
	return s.session.Query(query).Exec()
}

func (s Storage) tableName() string {
	return fmt.Sprintf("%s.%s", s.keyspace, s.table)
}
