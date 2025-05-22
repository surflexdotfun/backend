package config

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

const RPC_LATEST_CACHE_INTERVAL = 1 * time.Second // 1s

const INITIAL_USD_BALANCE = float64(10000)

var COIN_LIST = []string{
	"BTCUSDT",
	"ETHUSDT",
	"SUIUSDT",
}

var (
	RPC_URL          string
	PACKAGE_ID       string
	MODULE_NAME      = "surflex"
	STAKE_EVENT_NAME = "StakeEvent"
	SQLITE_DB_PATH   string
	ADMIN_PUBKEY     string
)

func init() {
	// 테스트 환경이면 init 실행 안 함
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, "-test.") {
			log.Println("Skipping init() in test mode")
			return
		}
	}

	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	adminPubkey := os.Getenv("ADMIN_PUBKEY")
	if adminPubkey == "" {
		log.Fatal("ADMIN_PUBKEY is not set in the environment variables")
	}
	ADMIN_PUBKEY = adminPubkey

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		log.Fatal("DB_PATH is not set in the environment variables")
	}
	SQLITE_DB_PATH = dbPath

	rpcUrl := os.Getenv("RPC_URL")
	if rpcUrl == "" {
		log.Fatal("RPC_URL is not set in the environment variables")
	}
	RPC_URL = rpcUrl

	packageId := os.Getenv("PACKAGE_ID")
	if packageId == "" {
		log.Fatal("PACKAGE_ID is not set in the environment variables")
	}
	PACKAGE_ID = packageId
}
