package middlewares

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"github.com/Vrg26/shortener-tpl/internal/app/types"
	"net/http"
)

const userKey types.ContextKey = 0

func Auth(secretKey string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var id uint32
			user, err := r.Cookie("User")
			if err != nil {
				idBytes, err := generateBytesUserID()
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				user = createAuthCookie(idBytes, secretKey)
				http.SetCookie(w, user)
				id = binary.BigEndian.Uint32(idBytes)
			} else {
				id, err = decodeAuthCookie(user.Value, secretKey)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			}
			ctx := context.WithValue(r.Context(), userKey, id)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func decodeAuthCookie(value string, secretKey string) (uint32, error) {
	data, err := hex.DecodeString(value)
	if err != nil {
		return 0, err
	}

	id := binary.BigEndian.Uint32(data[:8])

	hm := hmac.New(sha256.New, []byte(secretKey))
	hm.Write(data[:8])
	sign := hm.Sum(nil)
	if hmac.Equal(data[8:], sign) {
		return id, nil
	}
	return 0, http.ErrNoCookie
}

func generateBytesUserID() ([]byte, error) {
	id := make([]byte, 8)

	if _, err := rand.Read(id); err != nil {
		return nil, err
	}
	return id, nil
}

func createAuthCookie(id []byte, secretKey string) *http.Cookie {
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write(id)
	sign := hex.EncodeToString(append(id, h.Sum(nil)...))
	return &http.Cookie{
		Name:  "User",
		Value: sign,
	}
}
