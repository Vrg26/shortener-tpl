package middlewares

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"net/http"
)

func Auth(secretKey string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var id uint32
			user, err := r.Cookie("User")
			if err != nil {
				idBytes, err := generateBytesUserId()
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				user = createAuthCookie(idBytes, secretKey)
				http.SetCookie(w, user)
				id = binary.BigEndian.Uint32(idBytes)
			}
			id, err = decodeAuthCookie(user.Value, secretKey)
			ctx := context.WithValue(r.Context(), "user", id)
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

func generateBytesUserId() ([]byte, error) {
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
