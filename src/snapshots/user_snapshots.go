package snapshots

import (
	"log"
	"context"
	"minitwit/src/handlers/repo_wrappers"
	"minitwit/src/metrics"
	"time"
)

func RunUserSnapshotsAsync(ticker *time.Ticker) {
	ctx := context.Background()
	go func() {
		defer ticker.Stop()

		for {
			<-ticker.C
			updateTotalUsers(ctx)
		}
	}()
}

func updateTotalUsers(ctx context.Context) {
	log.Printf("📸 Info: Updating TotalUsers Snapshot")
	count, err := repo_wrappers.CountAllUsers(ctx)
	if err != nil {
		log.Printf("❌ Snapshot Error: counting all users: %v", err)
	}
	metrics.TotalUsers.Set(float64(count))
}