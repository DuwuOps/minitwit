package helpers

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func ParseFollowerBuckets(envVarName string) ([][2]uint32, error) {
	val := os.Getenv(envVarName)
	if val == "" {
		return nil, fmt.Errorf("environment variable %s is not set", envVarName)
	}

	pairs := strings.Split(val, ",")
	buckets := make([][2]uint32, 0, len(pairs))

	for _, pair := range pairs {
		if strings.Count(pair, "-") <= 1 {
			return nil, fmt.Errorf("invalid format in %s, expecting something like '100-200'", pair)
		}
		
		rng := strings.Split(pair, "-")
		low, err := strconv.ParseUint(rng[0], 10, 32)
		if err != nil {
			return nil, fmt.Errorf("failed to parse lower bound %s: %w", rng[0], err)
		}
		high, err := strconv.ParseUint(rng[1], 10, 32)
		if err != nil {
			return nil, fmt.Errorf("failed to parse upper bound %s: %w", rng[1], err)
		}

		buckets = append(buckets, [2]uint32{uint32(low), uint32(high)})
	}

	return buckets, nil
}
