package models

type LatestProcessed struct {
	LatestProcessedID int    `db:"latest_processed_id"`
	Date   int64  `db:"date_updated"`
}
