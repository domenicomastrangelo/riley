package auth

import "net/url"

func IsAuthorized(url *url.URL, userID uint64) bool {
	return true
}
