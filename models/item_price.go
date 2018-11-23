package models

type ItemPrice struct {
	SiteID       string `json:"site_id"`
	ProductID    string `json:"product_id"`
	ProductPrice float32 `json:"product_price"`
	Time         int64  `json:"time"`
}
