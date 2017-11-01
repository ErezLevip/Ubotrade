package Global

import (
	"github.com/nu7hatch/gouuid"
	"time"
)

type ActivityModel struct {
	Id                   uuid.UUID `json:"id"`
	BotId                uuid.UUID `json:"bot_id"`
	ActivityType         string    `json:"activity_type"`
	ActivityPrice        float64   `json:"activity_price"`
	ActualAmountUSD      float64   `json:"actual_amount_usd"`
	ActualAmountCurrency float64   `json:"actual_amount_currency"`
	PriceDifference      float64   `json:"price_difference"`
	TimeStamp            time.Time `json:"time_stamp"`
}
type BotTickerDataModel struct {
	Currency string    `json:"currency"`
	P        float64   `json:"p"`
	H        float64   `json:"h"`
	Stairs   []float64 `json:"stairs"`
	Action   string    `json:"action"`
}
type BotInformation struct {
	Id                   uuid.UUID     `json:"id"`
	Configuration        TradingConfig `json:"configuration"`
	OriginalAmount       float64       `json:"original_amount"`
	CurrencyAmount       float64       `json:"currency_amount"`
	LiquidAmountCurrency float64       `json:"liquid_amount_currency"`
	LiquidAmountUSD      float64       `json:"liquid_amount_usd"`
	Amount               float64       `json:"amount"`
	//IsActive bool			`json:"is_active"`
	//LastHealthCheck time.Time	`json:"last_health_check"`
	LastActivityId       uuid.UUID `json:"last_activity_id"`
	Name                 string    `json:"name"`
	UserId               string `json:"user_id"`
}
