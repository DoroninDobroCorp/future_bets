package entity

import "time"

type PriceRecord struct {
	Bookmaker string    `json:"bookmaker"`
	Score     float64   `json:"score"`
	CreatedAt time.Time `json:"createdAt"`
}

type FullPriceRecord struct {
	First     PriceRecord `json:"first"`
	Second    PriceRecord `json:"second"`
	Outcome   string      `json:"outcome"`
	ROI       float64     `json:"roi"`
	Margin    float64     `json:"margin"`
}

type ResponsePriceRecord struct {
	Key       string      `json:"key"`
	First     PriceRecord `json:"first"`
	Second    PriceRecord `json:"second"`
	Outcome   string      `json:"outcome"`
	ROI       float64     `json:"roi"`
	Margin    float64     `json:"margin"`
	CreatedAt time.Time   `json:"createdAt"`
}
