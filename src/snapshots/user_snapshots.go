package snapshots

import (
	"context"
	"log/slog"
	"minitwit/src/handlers/repo_wrappers"
	"minitwit/src/metrics"
	"minitwit/src/utils"
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
	slog.InfoContext(ctx, "📸 Info: Updating TotalUsers Snapshot")
	count, err := repo_wrappers.CountAllUsers(ctx)
	if err != nil {
		utils.LogErrorContext(ctx, "❌ Snapshot Error: counting all users", err)
	}
	metrics.TotalUsers.Set(float64(count))
}