package db

import (
	"surflex-backend/common/config"
	"surflex-backend/common/model"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func Init() {
	conn, err := gorm.Open(sqlite.Open(config.SQLITE_DB_PATH), &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		panic(err.Error())
	}

	conn.AutoMigrate(&model.ChartData{})
	conn.AutoMigrate(&model.Position{})
	conn.AutoMigrate(&model.Account{})
	conn.AutoMigrate(&model.AddressName{})
	conn.AutoMigrate(&model.KeyValueStore{})
	conn.AutoMigrate(&model.LeaderBoard{})
	conn.AutoMigrate(&model.DepositEvent{})

	db = conn
}

func GetConnection() *gorm.DB {
	return db
}
