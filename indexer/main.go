package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"surflex-backend/common/db"
	"surflex-backend/indexer/cache"
	"surflex-backend/indexer/export"
	"surflex-backend/indexer/leaderboard"
	"surflex-backend/indexer/okx"
	"surflex-backend/indexer/round"
)

func main() {
	flag.Parse()
	db.Init()

	cache.Init()

	okxFetcher := okx.NewFetcher()
	okxFetcher.Fetch()

	roundManager := round.NewRoundManager()
	roundManager.Start()

	LeaderBoardManager := leaderboard.NewLeaderBoardManager()
	LeaderBoardManager.Start()

	programExporter := export.NewProgramExporter()
	programExporter.Export()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Println("Gracefully shutting down...")
}
