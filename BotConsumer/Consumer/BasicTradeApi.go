package BotConsumer

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/beldur/kraken-go-api-client"
)

type TradeApi interface {
	Bid(amount float64, price float64) float64
	Take(amount float64, price float64) float64
	Buy(amount float64, price float64) bool
	GetTickerInfo() (askPrice float64, highPrice float64)
	GetSMA(timespan string) float64
	Make(currency string, retryTimeDuration time.Duration)
}

type BasicTradeApi struct {
	testingIndex      int
	lastPrice         float64
	KrakenApi         *krakenapi.KrakenApi
	pair              string
	retryTimeDuration time.Duration
}

func (api *BasicTradeApi) Make(currency string, retryTimeDuration time.Duration) {
	if currency == BitcoinCurrency {
		api.pair = KrakenBitcoinPair
	} else if currency == EtherCurrency {
		api.pair = KrakenEtherPair
	} else {
		panic("Unknown currency")
	}
	api.retryTimeDuration = retryTimeDuration
	api.KrakenApi = krakenapi.New("KEY", "SECRET")
}

const BitcoinCurrency = "Bitcoin"
const EtherCurrency = "Ether"
const KrakenBitcoinPair = "XXBTZUSD"
const KrakenEtherPair = "XETHZUSD"
const BlockChainTimespan30Days = "30days"
const BlockChainTimespan60Days = "60days"
const BlockChainTimespan180Days = "180days"
const BlockChainTimespan1Year = "1year"
const BlockChainTimespan2Years = "2year"

type BlockChainApiResponse struct {
	Status      string `json:"status"`
	Name        string `json:"name"`
	Unit        string `json:"unit"`
	Period      string `json:"period"`
	Description string `json:"description"`
	Values      []struct {
		X int     `json:"x"`
		Y float64 `json:"y"`
	} `json:"values"`
}

func (api *BasicTradeApi) GetTickerInfo() (askPrice float64, highPrice float64) {
	for {
		res, err := api.KrakenApi.Query("Ticker", map[string]string{
			"pair": api.pair,
		})

		if err != nil {
			fmt.Println(err)
		}
		if res != nil {
			queryResult := (res.(map[string]interface{})[api.pair]).(map[string]interface{})
			askPriceStr := queryResult["a"].([]interface{})[0].(string)
			askPrice, err = strconv.ParseFloat(askPriceStr, 64)

			highPriceStr := queryResult["h"].([]interface{})[0].(string)
			highPrice, err = strconv.ParseFloat(highPriceStr, 64)

			if err != nil {
				fmt.Println(err)
			}
			return askPrice, highPrice
		}
		log.Println(time.Now(),"api response is empty")
		time.Sleep(api.retryTimeDuration)
	}
}
func (api *BasicTradeApi) GetSMA(timespan string) float64 {
	for i := 0; i <= 3; i++ {
		var url2 = "https://api.blockchain.info/charts/market-price?format=json&timespan=" + timespan
		res, err := http.Get(url2)
		if err != nil {
			log.Print(err.Error())
		} else {
			resp := BlockChainApiResponse{}
			_ = json.NewDecoder(res.Body).Decode(&resp)
			sum := 0.0
			for j := 0; j < len(resp.Values); j++ {
				sum += resp.Values[j].Y
			}
			return sum / float64(len(resp.Values))
		}
		log.Println(err)
	}
	return 0.0
}

func (api *BasicTradeApi) Bid(amount float64, price float64) float64 {
	var commission float64 = 0.9998
	return amount * price * commission //stimulate commission
}
func (api *BasicTradeApi) Take(amount float64, price float64) float64 {
	var commission float64 = 0.9998
	return amount * price * commission //stimulate commission
}
func (api *BasicTradeApi) Buy(amount float64, price float64) bool {
	return true
}
