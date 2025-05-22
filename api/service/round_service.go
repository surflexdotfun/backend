package service

import (
	"log"
	"surflex-backend/api/config"
	"surflex-backend/api/types"
	"surflex-backend/common/db"
	"surflex-backend/common/model"
	"time"
)

var round uint64 = 0

func GetRound() types.GetRoundResponse {
	return types.GetRoundResponse{
		Round: round,
	}
}

func StartRoundUpdater() {
	currentRound, err := fetchRoundFromDB()
	if err != nil {
		log.Println(err)
		return
	}
	round = currentRound
	ticker := time.NewTicker(config.ROUND_REFRESH_INTERVAL)
	go func() {
		for {
			<-ticker.C
			currentRound, err := fetchRoundFromDB()
			if err != nil {
				log.Println(err)
				continue
			}
			round = currentRound
		}
	}()
}

func fetchRoundFromDB() (uint64, error) {
	conn := db.GetConnection()

	var currentRound uint64
	if err := conn.Model(&model.KeyValueStore{}).
		Select("value").
		Where("key = ?", model.CURRENT_ROUND_KEY).
		Find(&currentRound).Error; err != nil {
		log.Println(err)
		return 0, err
	}
	return currentRound, nil
}
