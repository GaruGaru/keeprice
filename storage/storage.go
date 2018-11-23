package storage

import "github.com/GaruGaru/keeprice/models"

type PriceStorage interface {
	Store(itemPrice models.ItemPrice) error
	Get(siteID string, productID string) (models.ItemPriceHistory, error)
}
