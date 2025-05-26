package auth

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/apetsko/gophkeeper/utils"
	"github.com/gorilla/securecookie"
)

func securedCookie(secret string) *securecookie.SecureCookie {
	secretLen := 32
	id := utils.GenerateID(secret, secretLen)
	sharedSecret := []byte(id)
	return securecookie.New(sharedSecret, sharedSecret)
}

func CookieSetUserID(w http.ResponseWriter, userID int, secret string) (err error) {
	sc := securedCookie(secret)

	encoded, err := sc.Encode("gophkeeper ", userID)
	if err != nil {
		err = fmt.Errorf("error encoding userid cookie: %v", err)
		return err
	}

	oneDay := time.Hour * 24
	http.SetCookie(w, &http.Cookie{
		Name:     "gophkeeper ",
		Value:    encoded,
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Now().Add(oneDay),
		SameSite: http.SameSiteLaxMode,
	})
	return nil
}

func CookieGetUserID(r *http.Request, secret string) (userID *int, err error) {
	cookie, err := r.Cookie("gophkeeper ")
	if err != nil {
		return nil, http.ErrNoCookie
	}

	if err := cookie.Valid(); err != nil {
		return nil, http.ErrNoCookie
	}

	sc := securedCookie(secret)
	if err := sc.Decode("gophkeeper ", cookie.Value, &userID); err != nil {
		err = fmt.Errorf("error decoding user cookie: %w", err)
		return nil, err
	}

	if userID == nil {
		return nil, errors.New("userid not found in cookie")
	}
	return userID, nil
}
