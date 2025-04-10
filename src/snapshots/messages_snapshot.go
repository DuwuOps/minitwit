package snapshots

import (
	"context"
	"log"
	"minitwit/src/handlers/repo_wrappers"
	"minitwit/src/metrics"
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
	log.Printf("📸 Info: Updating MessagesTotal Snapshots")
	count, err := repo_wrappers.CountAllMessages(ctx)
	if err != nil {
		log.Printf("❌ Snapshot Error: counting all messages: %v", err)
	}
	metrics.MessagesTotal.Set(float64(count))
}

func updateFlaggedMessagesTotal(ctx context.Context) {
	log.Printf("📸 Info: Updating FlaggedMessagesTotal Snapshots")
	condition := map[string]any{
		"flagged": 1,
	}

	count, err := repo_wrappers.CountFilteredMessages(ctx, condition)
	if err != nil {
		log.Printf("❌ Snapshot Error: counting all flagged messages: %v", err)
	}

	metrics.FlaggedMessagesTotal.Set(float64(count))
}
