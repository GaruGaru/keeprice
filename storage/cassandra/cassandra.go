package cassandra

import (
	"fmt"
	"github.com/GaruGaru/keeprice/models"
	"github.com/gocql/gocql"
)

type CassandraStorage struct {
	session           gocql.Session
	keyspace          string
	table             string
	replicationFactor int
}

func NewCassandraStorage(config Config) (CassandraStorage, error) {

	if config.Consistency == 0 {
		config.Consistency = gocql.Quorum
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
	session, err := cluster.CreateSession()

	if err != nil {
		return CassandraStorage{}, err
	}

	return CassandraStorage{
		session:           *session,
		keyspace:          config.KeySpace,
		table:             config.Table,
		replicationFactor: config.ReplicationFactor,
	}, nil
}

func (s CassandraStorage) Init() error {
	if err := s.createKeyspaceIfNotExists(); err != nil {
		return err
	}
	if err := s.createTableIfNotExists(); err != nil {
		return err
	}
	return nil
}

func (s CassandraStorage) Store(itemPrice models.ProductPrice) error {
	queryTemplate := `INSERT INTO %s (site_id, product_id, time, price) VALUES (%s, %s, %d, %f)`
	query := fmt.Sprintf(queryTemplate, s.tableName(), itemPrice.SiteID, itemPrice.ProductID, itemPrice.Time, itemPrice.ProductPrice)
	return s.session.Query(query).Exec()
}

func (s CassandraStorage) Get(siteID string, productID string) (models.ProductPriceHistory, error) {
	return models.ProductPriceHistory{}, fmt.Errorf("not implmemented yet")
}

func (s CassandraStorage) createKeyspaceIfNotExists() error {
	queryTemplate := `CREATE KEYSPACE IF NOT EXISTS %s
	WITH replication = {
		'class' : 'SimpleStrategy',
		'replication_factor' : %d
	}`
	query := fmt.Sprintf(queryTemplate, s.keyspace, s.replicationFactor)
	return s.session.Query(query).Exec()
}

func (s CassandraStorage) createTableIfNotExists() error {
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

func (s CassandraStorage) tableName() string {
	return fmt.Sprintf("%s.%s", s.keyspace, s.table)
}
