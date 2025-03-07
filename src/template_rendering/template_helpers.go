package template_rendering

import (
	"crypto/md5"
	"fmt"
	"strings"
	"time"
)

func gravatarUrl(email string, size int) string {

	email = strings.TrimSpace(strings.ToLower(email))
	hash := md5.Sum([]byte(email))

	// Convert hash to hexadecimal string
	hashString := fmt.Sprintf("%x", hash)

	return fmt.Sprintf("http://www.gravatar.com/avatar/%s?d=identicon&s=%d", hashString, size)
}

func formatDatetime(timestamp int64) string {
	time := time.Unix(timestamp, 0).UTC().Local()
	return time.Format("2006-01-02 @ 15:04")
}
