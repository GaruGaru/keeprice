package models

type ItemPriceHistoryEntry struct {
	Time  int64   `json:"time"`
	Price float32 `json:"price"`
}

type ItemPriceHistory struct {
	Count       int                     `json:"count"`
	PeriodStart int64                     `json:"period_start"`
	PeriodEnd   int64                     `json:"period_end"`
	History     []ItemPriceHistoryEntry `json:"history"`
}
