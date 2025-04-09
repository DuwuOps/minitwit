package snapshots

import (
	"log"
	"os"
	"strconv"
	"time"
)

func RecordSnapshots() {
	intervalStr := os.Getenv("SNAPSHOT_TIME_INTERVAL_SECONDS")
	if intervalStr == "" {
		log.Printf("ℹ️ Info: unable to get SNAPSHOT_TIME_INTERVAL_SECONDS, using default.")
		intervalStr = "300"
	}

	snapshotInterval, err := strconv.Atoi(intervalStr)
	if err != nil {
		log.Printf("❌ Error: Invalid SNAPSHOT_TIME_INTERVAL_SECONDS: %v", err)
		return
	}

	ticker := time.NewTicker(time.Duration(snapshotInterval) * time.Second)

	RunUserSnapshotsAsync(ticker)
}