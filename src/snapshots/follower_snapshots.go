package snapshots

import (
	"context"
	"fmt"
	"log/slog"
	"minitwit/src/handlers/helpers"
	"minitwit/src/handlers/repo_wrappers"
	"minitwit/src/metrics"
	"minitwit/src/utils"
	"time"
)

func RunFollowerSnapshotsAsync(ticker *time.Ticker) {
	ctx := context.Background()

	bounds, err := helpers.ParseFollowerBuckets("FOLLOWER_BUCKETS")
	if err != nil {
		utils.LogError("❌ Error parsing follower buckets", err)
		return
	}

	go func() {
		defer ticker.Stop()

		for {
			<-ticker.C
			updateFollowerMetrics(ctx, bounds)
		}
	}()
}

func updateFollowerMetrics(ctx context.Context, bounds [][2]uint32) {
	slog.InfoContext(ctx, "📸 Info: Updating Follower Snapshots")
	for _, bound := range bounds {
		low, high := int(bound[0]), int(bound[1])

		followerCount, err := repo_wrappers.CountFieldInRange(ctx, "following_id", low, high)
		if err != nil {
			rangeStr := fmt.Sprintf("%d-%d", low, high)
			slog.Error("❌ Snapshot Error: counting followers in given range", slog.Any("error", err), slog.Any("range", rangeStr))
		}

		followeesCount, err := repo_wrappers.CountFieldInRange(ctx, "follower_id", low, high)
		if err != nil {
			rangeStr := fmt.Sprintf("%d-%d", low, high)
			slog.Error("❌ Snapshot Error: counting followers in given range", slog.Any("error", err), slog.Any("range", rangeStr))
		}

		metrics.FollowerTotal.WithLabelValues(fmt.Sprintf("%d-%d", low, high)).Set(float64(followerCount))
		metrics.FolloweesTotal.WithLabelValues(fmt.Sprintf("%d-%d", low, high)).Set(float64(followeesCount))
	}
}
