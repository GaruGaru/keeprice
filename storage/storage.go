package storage

import "github.com/GaruGaru/keeprice/models"

type PriceStorage interface {
	Store(itemPrice models.ProductPrice) error
	Get(siteID string, productID string) (models.ProductPriceHistory, error)
}
