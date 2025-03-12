package template_rendering

import (
	"crypto/md5"
	"fmt"
	"strings"
)

func gravatarUrl(email string, size int) string {

	email = strings.TrimSpace(strings.ToLower(email))
	hash := md5.Sum([]byte(email))

	// Convert hash to hexadecimal string
	hashString := fmt.Sprintf("%x", hash)

	return fmt.Sprintf("http://www.gravatar.com/avatar/%s?d=identicon&s=%d", hashString, size)
}
