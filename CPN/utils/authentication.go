package utils

import (
	"github.com/gorilla/securecookie"
)

type Cookie struct {
	UserID string
}

var (
	hashKey       = []byte("LadderCompetitionPlatform")
	blockKey      = securecookie.GenerateRandomKey(16)
	cookieManager = securecookie.New(hashKey, blockKey)
)

func CookieEncoder(userId string) string {
	cookie := Cookie{
		UserID: userId,
	}
	encodedCookie, err := cookieManager.Encode("cookie", cookie)
	if err != nil {
		Logger.Errorf("UserID: %s\nError: %s", userId, err)
	}
	return encodedCookie
}

func CookieDecoder(cookie string) string {
	decodedCookie := &Cookie{}
	err := cookieManager.Decode("cookie", cookie, decodedCookie)
	if err != nil {
		Logger.Errorf("Cookie: %s\nError: %s", cookie, err)
		return ""
	}
	return decodedCookie.UserID
}
