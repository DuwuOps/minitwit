package snapshots

import (
	"log/slog"
	"minitwit/src/utils"
	"os"
	"strconv"
	"time"
)

func RecordSnapshots() {
	intervalStr := os.Getenv("SNAPSHOT_TIME_INTERVAL_SECONDS")
	if intervalStr == "" {
		slog.Info("ℹ️ Info: unable to get SNAPSHOT_TIME_INTERVAL_SECONDS, using default.")
		intervalStr = "300"
	}

	snapshotInterval, err := strconv.Atoi(intervalStr)
	if err != nil {
		utils.LogError("❌ Error: Invalid SNAPSHOT_TIME_INTERVAL_SECONDS", err)
		return
	}

	ticker := time.NewTicker(time.Duration(snapshotInterval) * time.Second)

	RunUserSnapshotsAsync(ticker)
	RunFollowerSnapshotsAsync(ticker)
	RunMessagesSnapshotsAsync(ticker)
}
