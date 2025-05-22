package okx

import (
	"fmt"
	"log"
	"surflex-backend/common/config"
	"surflex-backend/common/db"
	"surflex-backend/common/model"
	"surflex-backend/indexer/cache"
	"surflex-backend/indexer/liquidation"

	"gorm.io/gorm/clause"
)

type PriceFetcher struct{}

func NewPriceFetcher() *PriceFetcher {
	return &PriceFetcher{}
}

func (e *PriceFetcher) Fetch() {
	prices, err := GetTokenPrices(config.COIN_LIST)
	if err != nil {
		log.Println(err)
		return
	}
	cache.SetPrice(prices)
	liquidation.FindAndLiquidate()

	db := db.GetConnection()
	var keyValueStores []model.KeyValueStore
	for coinID, price := range prices {
		keyValueStores = append(keyValueStores, model.KeyValueStore{
			Key:   coinID,
			Value: fmt.Sprintf("%f", price),
		})
	}
	db.Clauses(clause.OnConflict{
		DoUpdates: clause.AssignmentColumns([]string{"value"}),
	}).Create(keyValueStores)
}
