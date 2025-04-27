package snapshots

import (
	"context"
	"fmt"
	"log"
	"minitwit/src/handlers/helpers"
	"minitwit/src/handlers/repo_wrappers"
	"minitwit/src/metrics"
	"time"
)

func RunFollowerSnapshotsAsync(ticker *time.Ticker) {
	ctx := context.Background()

	bounds, err := helpers.ParseFollowerBuckets("FOLLOWER_BUCKETS")
	if err != nil {
		log.Printf("❌ Error parsing follower buckets: %v", err)
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
	log.Printf("📸 Info: Updating Follower Snapshots")
	for _, bound := range bounds {
		low, high := int(bound[0]), int(bound[1])

		followerCount, err := repo_wrappers.CountFieldInRange(ctx, "following_id", low, high)
		if err != nil {
			log.Printf("❌ Snapshot Error: counting followers in range %d-%d: %v", low, high, err)
		}

		followeesCount, err := repo_wrappers.CountFieldInRange(ctx, "follower_id", low, high)
		if err != nil {
			log.Printf("❌ Snapshot Error: counting followees in range %d-%d: %v", low, high, err)
		}

		metrics.FollowerTotal.WithLabelValues(fmt.Sprintf("%d-%d", low, high)).Set(float64(followerCount))
		metrics.FolloweesTotal.WithLabelValues(fmt.Sprintf("%d-%d", low, high)).Set(float64(followeesCount))
	}
}
