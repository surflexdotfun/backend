package binance

import (
	"errors"
	"surflex-backend/common/config"
	"surflex-backend/common/db"
	"surflex-backend/common/model"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

type ChartFetcher struct{}

func NewChartFetcher() *ChartFetcher {
	return &ChartFetcher{}
}

func (e *ChartFetcher) Fetch() {
	for _, coinID := range config.COIN_LIST {
		// Get latest timestamp from ChartData for this coinID
		var from time.Time
		err := db.GetConnection().Model(&model.ChartData{}).
			Select("close_time").
			Where("symbol = ?", coinID).
			Order("close_time desc").
			First(&from).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			// If no existing data found, set from to 30 days ago
			from = time.Now().AddDate(0, 0, -30)
		} else if err != nil {
			log.Error(err)
			continue
		}

		chartData, err := GetChartData(coinID, from)
		if err != nil {
			log.Info(err)
			continue
		}

		if len(chartData) == 0 {
			continue
		}
		// Save to database
		if err := db.GetConnection().Create(&chartData).Error; err != nil {
			log.Error(err)
		}

		time.Sleep(time.Millisecond * 500) // Prevent request limit
	}
}
