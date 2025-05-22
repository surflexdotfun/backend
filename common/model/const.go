package model

const CURRENT_ROUND_KEY = "current_round"
const CURRENT_ROUND_STATUS_KEY = "current_round_status"
const NEXT_CURSOR = "next_cursor"

type RoundStatus string

const (
	RoundStatusStart RoundStatus = "start"
	RoundStatusEnd   RoundStatus = "end"
	RoundStatusReady RoundStatus = "ready"
)
