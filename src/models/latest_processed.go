package models

type LatestProcessed struct {
	LatestProcessedID int64 `db:"latest_processed_id"`
	DateUpdated       int64 `db:"date_updated"`
}
