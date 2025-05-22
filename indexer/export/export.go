package export

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	commonconfig "surflex-backend/common/config"
	"surflex-backend/common/db"
	"surflex-backend/common/model"
	"surflex-backend/indexer/cache"
	"surflex-backend/indexer/config"
	"time"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var EVENT_LIMIT = uint64(5)

type ProgramExporter struct {
	client           sui.ISuiAPI
	stakeEventFilter interface{}
}

func NewProgramExporter() *ProgramExporter {
	client := sui.NewSuiClient(commonconfig.RPC_URL)

	moveType := commonconfig.PACKAGE_ID + "::" + commonconfig.MODULE_NAME + "::" + commonconfig.STAKE_EVENT_NAME
	stakeEventFilter := map[string]interface{}{
		"MoveModule": map[string]string{
			"package": commonconfig.PACKAGE_ID,
			"module":  commonconfig.MODULE_NAME,
			"type":    moveType,
		},
	}

	return &ProgramExporter{
		client:           client,
		stakeEventFilter: stakeEventFilter,
	}
}

func (e *ProgramExporter) Export() {
	go func() {
		for {
			log.Println("Run ProgramExporter")
			start := time.Now()
			e.FetchProgramData()
			log.Println("Run ProgramExporter Completed, elapsed:", time.Since(start))
			time.Sleep(config.DEPOSIT_UPDATE_INTERVAL)
		}
	}()
}

func (e *ProgramExporter) FetchProgramData() error {
	// 1. 최근에 처리한 트랜잭션을 가져온다.
	db := db.GetConnection()
	var nextCursorStr string
	var nextCursor *models.EventId

	db.Model(&model.KeyValueStore{}).
		Select("value").
		Where("key = ?", model.NEXT_CURSOR).
		Find(&nextCursorStr)
	if len(nextCursorStr) != 0 {
		nextCursor = &models.EventId{}
		if err := json.Unmarshal([]byte(nextCursorStr), nextCursor); err != nil {
			return err
		}
	}

	events, err := e.extractParticipatedEvents(nextCursor)
	if err != nil {
		log.Println("Failed to extract participated events:", err)
		return err
	}

	for _, event := range events {
		if err := e.saveDepositAndCreateAccount(event); err != nil {
			return err
		}
	}

	return nil
}

func (e *ProgramExporter) extractParticipatedEvents(nextCursor *models.EventId) ([]models.SuiEventResponse, error) {
	res, err := e.client.SuiXQueryEvents(
		context.Background(),
		models.SuiXQueryEventsRequest{
			SuiEventFilter:  e.stakeEventFilter,
			Cursor:          nextCursor,
			Limit:           EVENT_LIMIT,
			DescendingOrder: false,
		},
	)

	if err != nil {
		return nil, err
	}
	return res.Data, nil
}

func (e *ProgramExporter) saveDepositAndCreateAccount(event models.SuiEventResponse) error {
	db := db.GetConnection()
	currentRound := cache.GetRound()

	address := event.Sender
	nextCursorByte, err := json.Marshal(event.Id)
	if err != nil {
		return err
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		// round와 address가 같고 deleted가 false인 account가 이미 있다면 deleted를 true로 업데이트
		if err := tx.Model(&model.Account{}).
			Where("round = ? AND address = ? AND deleted = ?", currentRound, address, false).
			Update("deleted", true).Error; err != nil {
			return err
		}

		if err := tx.Model(&model.Position{}).
			Where("round = ? AND address = ? AND status = ?", currentRound, address, model.StatusActive).
			Update("status", model.StatusRemoved).Error; err != nil {
			return err
		}

		keyValueStore := model.KeyValueStore{
			Key:   model.NEXT_CURSOR,
			Value: string(nextCursorByte),
		}
		if err := tx.Clauses(clause.OnConflict{
			DoUpdates: clause.AssignmentColumns([]string{"value"}),
		}).Create(&keyValueStore).Error; err != nil {
			return err
		}

		// 새로운 account 생성
		if err := tx.Create(&model.Account{
			Round:     currentRound,
			Address:   address,
			USDAmount: config.USER_USD_AMOUNT,
		}).Error; err != nil {
			return err
		}

		// depositEvent 생성
		if err := tx.Create(&model.DepositEvent{
			CursorId: nextCursorByte,
			Address:  address,
		}).Error; err != nil {
			return err
		}

		return nil
	})
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		// 이미 존재하는 데이터는 무시한다.
		return nil
	}
	return err
}
