package service

import (
	"log"
	"strconv"
	"surflex-backend/api/config"
	"surflex-backend/api/types"
	commonconfig "surflex-backend/common/config"
	"surflex-backend/common/db"
	"surflex-backend/common/model"
	"time"
)

var priceMap map[string]float64

func GetPrice() types.GetPriceResponse {
	return priceMap
}

func StartPriceUpdater() {
	prices, err := FetchPriceFromDB(commonconfig.COIN_LIST)
	if err != nil {
		log.Println(err)
		return
	}
	priceMap = prices
	ticker := time.NewTicker(config.PRICE_REFRESH_INTERVAL)
	go func() {
		for {
			<-ticker.C
			prices, err := FetchPriceFromDB(commonconfig.COIN_LIST)
			if err != nil {
				log.Println(err)
				continue
			}
			priceMap = prices
		}
	}()
}

func FetchPriceFromDB(coinIDs []string) (map[string]float64, error) {
	conn := db.GetConnection()

	var keyValueStores []model.KeyValueStore
	err := conn.Where("key IN ?", coinIDs).Find(&keyValueStores).Error
	if err != nil {
		return nil, err
	}

	newPriceMap := make(map[string]float64)
	for _, keyValueStore := range keyValueStores {
		price, err := strconv.ParseFloat(keyValueStore.Value, 64)
		if err != nil {
			return nil, err
		}
		newPriceMap[keyValueStore.Key] = price
	}
	return newPriceMap, nil
}
