package snapshots

import (
	"context"
	"log/slog"
	"minitwit/src/handlers/repo_wrappers"
	"minitwit/src/metrics"
	"minitwit/src/utils"
	"time"
)

func RunMessagesSnapshotsAsync(ticker *time.Ticker) {
	ctx := context.Background()
	go func() {
		defer ticker.Stop()

		for {
			<-ticker.C
			updateMessagesTotal(ctx)
			updateFlaggedMessagesTotal(ctx)
		}
	}()
}

func updateMessagesTotal(ctx context.Context) {
	slog.InfoContext(ctx, "📸 Info: Updating MessagesTotal Snapshots")
	count, err := repo_wrappers.CountAllMessages(ctx)
	if err != nil {
		utils.LogErrorContext(ctx, "❌ Snapshot Error: counting all messages", err)
	}
	metrics.MessagesTotal.Set(float64(count))
}

func updateFlaggedMessagesTotal(ctx context.Context) {
	slog.InfoContext(ctx, "📸 Info: Updating FlaggedMessagesTotal Snapshots")
	condition := map[string]any{
		"flagged": 1,
	}

	count, err := repo_wrappers.CountFilteredMessages(ctx, condition)
	if err != nil {
		utils.LogErrorContext(ctx, "❌ Snapshot Error: counting all flagged messages", err)
	}

	metrics.FlaggedMessagesTotal.Set(float64(count))
}
